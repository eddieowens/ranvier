package ranvier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/eddieowens/ranvier/server/app/collections"
	"github.com/eddieowens/ranvier/server/app/exchange/response"
	"github.com/eddieowens/ranvier/server/app/model"
	"github.com/gorilla/websocket"
	"github.com/oliveagle/jsonpath"
	"github.com/pkg/errors"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"
)

type ClientOptions struct {
	// The required hostname to your Ranvier server(s). Should not include the protocol, e.g. if Ranvier is pointed at the
	// url https://ranvier.mycompany.com, this hostname should strictly be ranvier.mycompany.com. Or if it's
	// https://ranvier:8080, this field should be set to ranvier:8080.
	Hostname string

	// The directory that your local config files will be stored on disk. Config is stored on disk to prevent a strong
	// reliance on the server. This directory will default to a temp directory.
	ConfigDirectory string
}

type ConnOptions struct {
	// The names of the configuration files that you want updates for. Whenever the configuration tied to this name is
	// updated, the Connection will receive a message of the new configuration.
	Names []string

	// If the websocket were to disconnect, how long until a retry is attempted? Default is 30 seconds.
	RetryInterval time.Duration
}

type Connection struct {
	// The channel for all incoming configuration changes. Listen for messages on this channel to get realtime updates on
	// configuration changes.
	OnUpdate chan model.ConfigEvent

	closer chan bool
}

// Create a new Ranvier client with the provided options. It is advised to only use a single client as all queries made
// by the client are cached.
func NewClient(options *ClientOptions) (Client, error) {
	_, _, err := net.SplitHostPort(options.Hostname)
	if err != nil {
		return nil, err
	}
	hostname := options.Hostname

	confDir := options.ConfigDirectory
	if confDir == "" {
		options.ConfigDirectory = path.Join(os.TempDir(), "ranvier", options.Hostname)
	}
	err = os.MkdirAll(options.ConfigDirectory, os.ModePerm)
	if err != nil {
		return nil, err
	}

	c := &clientImpl{
		Options: &ClientOptions{
			Hostname:        hostname,
			ConfigDirectory: confDir,
		},
		ConfigMap: collections.NewConfigMap(),
	}

	return c, nil
}

type QueryOptions struct {
	IgnoreCache bool
	Name        string
	Query       string
}

type Client interface {
	// Establish a websocket connection between the client and the Ranvier server. Whenever the Ranvier server detects
	// changes to the configuration you care about (specified in the options), it will notify your client and update its
	// state. If the connection could not be made, an error is returned.
	Connect(options *ConnOptions) (*Connection, error)

	// Sever the connection between the client and the Ranvier server. The client will no longer receive updates from the
	// Ranvier server after the connection is disconnected.
	Disconnect(conn *Connection)

	// Query Ranvier for some config. The query is a valid jsonpath query (https://restfulapi.net/json-jsonpath/). If the
	// query is unable to find any config, nil is returned. If the query is invalid, an error is returned.
	//
	// All queries are cached within the client and will be hit unless specified in the options. The order of operations
	// for retrieving a query are
	//   client cache -> Ranvier server -> local disk
	// If the query is unsuccessful in all of these operations, nil is returned.
	//
	// All successful queries will be written to disk.
	Query(options *QueryOptions) (*model.Config, error)
}

type clientImpl struct {
	Options   *ClientOptions
	ConfigMap collections.ConfigMap
}

// Only available when the Ranvier server is in Dev mode!
//
// Create a new Config in the Ranvier server. This will cause a websocket event to be emitted.
func (c *clientImpl) Create(config *model.Config) (*model.Config, error) {
	if config == nil {
		return nil, nil
	}

	d, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(c.Options.Hostname, "application/json", bytes.NewBuffer(d))
	if err != nil {
		return nil, err
	}

	d, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var conf response.Config
	err = json.Unmarshal(d, &conf)
	if err != nil {
		return nil, err
	}

	return conf.Data, nil
}

// Only available when the Ranvier server is in Dev mode!
//
// Deletes a Config object from the Ranvier server. This will cause a websocket event to be emitted.
func (c *clientImpl) Delete(name string) (*model.Config, error) {
	parsedUrl, _ := url.Parse(c.Options.Hostname)
	parsedUrl.Scheme = "http"
	parsedUrl.Path = "api/config/" + name
	req, err := http.NewRequest(http.MethodDelete, parsedUrl.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, err
	}

	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var respConf response.Config
	err = json.Unmarshal(d, &respConf)
	if err != nil {
		return nil, err
	}

	return respConf.Data, nil
}

// Only available when the Ranvier server is in Dev mode!
//
// Updates a pre-existing Config object. If the Config did not previously exist, a new one is created under the
// specified name.
func (c *clientImpl) Update(config *model.Config) (*model.Config, error) {
	if config == nil {
		return nil, nil
	}

	parsedUrl := &url.URL{
		Scheme: "http",
		Host:   c.Options.Hostname,
		Path:   "api/config/" + config.Name,
	}

	d, err := json.Marshal(config.Data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPut, parsedUrl.String(), bytes.NewBuffer(d))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, err
	}

	d, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var respConf response.Config
	err = json.Unmarshal(d, &respConf)
	if err != nil {
		return nil, err
	}

	return respConf.Data, nil
}

func (c *clientImpl) Disconnect(conn *Connection) {
	conn.closer <- true
	return
}

func (c *clientImpl) Connect(options *ConnOptions) (*Connection, error) {
	configChan := make(chan model.ConfigEvent, 0)
	closer := make(chan bool, 1)
	parsedUrl := &url.URL{
		Host:   c.Options.Hostname,
		Scheme: "ws",
	}
	retryInterval := options.RetryInterval
	if retryInterval == 0 {
		retryInterval = time.Second * 30
	}

	for _, n := range options.Names {
		parsedUrl.Path = fmt.Sprintf("/api/config/ws/%s", n)
		ws, _, err := websocket.DefaultDialer.Dial(parsedUrl.String(), nil)
		if err != nil {
			return nil, err
		}

		go func(ws *websocket.Conn, closer chan bool, configChan chan model.ConfigEvent) {
			for {
				if len(closer) > 0 {
					ws.Close()
					return
				}
				conf := new(model.ConfigEvent)
				err = ws.ReadJSON(conf)
				if err != nil {
					if _, ok := err.(*websocket.CloseError); ok {
						time.Sleep(retryInterval)
						ws, _, err = websocket.DefaultDialer.Dial(parsedUrl.String(), nil)
					}
					continue
				}

				_ = c.writeToDisk(&conf.Config)
				c.ConfigMap.Set(conf.Config.Name, conf.Config)
				configChan <- *conf
			}
		}(ws, closer, configChan)
	}

	conn := Connection{
		closer:   closer,
		OnUpdate: configChan,
	}

	return &conn, nil
}

func (c *clientImpl) Query(options *QueryOptions) (*model.Config, error) {
	if options == nil {
		return nil, errors.New("options are required to run a query")
	}
	conf, err := c.fetchAndLoadConfig(options.Name, options.IgnoreCache)
	if err != nil {
		return nil, err
	}
	conf, err = c.query(conf, options.Query)
	if err != nil {
		return nil, err
	}

	return conf, nil
}

func (c *clientImpl) toUrl(name string, query string) string {
	return fmt.Sprintf("%s/api/config/%s?query%s", c.Options.Hostname, name, query)
}

func (c *clientImpl) query(config *model.Config, query string) (*model.Config, error) {
	raw, err := jsonpath.JsonPathLookup(config.Data, query)
	if err != nil {
		return nil, err
	}

	return &model.Config{
		Name: config.Name,
		Data: raw,
	}, nil
}

func (c *clientImpl) fetchAndLoadConfig(name string, ignoreCache bool) (*model.Config, error) {
	var exists bool
	var conf model.Config
	if !ignoreCache {
		conf, exists = c.ConfigMap.Get(name)
	}
	if !exists {
		fetchedCfg, err := c.fetchConfig(name, "")
		if err != nil {
			fetchedCfg, err = c.loadFromDisk(name)
			if err != nil {
				return nil, err
			}
		} else {
			err = c.writeToDisk(fetchedCfg)
			if err != nil {
				return nil, err
			}
		}
		conf = *fetchedCfg
	}

	c.ConfigMap.Set(conf.Name, conf)

	return &conf, nil
}

func (c *clientImpl) fetchConfig(name string, query string) (*model.Config, error) {
	resp, err := http.Get(c.toUrl(name, query))
	if err != nil {
		return nil, err
	}

	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var conf response.Config
	err = json.Unmarshal(d, &conf)
	if err != nil {
		return nil, err
	}

	return conf.Data, nil
}

func (c *clientImpl) loadFromDisk(name string) (*model.Config, error) {
	d, err := ioutil.ReadFile(path.Join(c.Options.ConfigDirectory, name))
	if err != nil {
		return nil, err
	}

	var conf model.Config
	err = json.Unmarshal(d, &conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}

func (c *clientImpl) writeToDisk(config *model.Config) error {
	d, err := json.Marshal(config)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path.Join(c.Options.ConfigDirectory, config.Name), d, os.ModePerm)
}
