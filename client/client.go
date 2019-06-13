package ranvier

import (
	"encoding/json"
	"fmt"
	"github.com/eddieowens/ranvier/server/app/collections"
	"github.com/eddieowens/ranvier/server/app/exchange/response"
	"github.com/eddieowens/ranvier/server/app/model"
	"github.com/gorilla/websocket"
	"github.com/oliveagle/jsonpath"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

type ClientOptions struct {
	Url             string
	ConfigDirectory string
}

type ConnOptions struct {
	Names []string
}

type Connection struct {
	OnUpdate chan model.Config
	closer   chan bool
}

func NewClient(options *ClientOptions) Client {
	url := options.Url
	if string(url[len(url)-1]) == "/" {
		options.Url = url[:len(url)-1]
	}

	confDir := options.ConfigDirectory
	if confDir == "" {
		options.ConfigDirectory = path.Join(os.TempDir(), "ranvier", options.Url)
	}
	err := os.MkdirAll(options.ConfigDirectory, os.ModePerm)
	if err != nil {
		panic(err)
	}

	c := &clientImpl{
		Options:   options,
		ConfigMap: collections.NewConfigMap(),
	}

	return c
}

type Client interface {
	Connect(options *ConnOptions) (*Connection, error)
	Disconnect(conn *Connection)
	Query(name string, query string) (*model.Config, error)
}

type clientImpl struct {
	Options   *ClientOptions
	ConfigMap collections.ConfigMap
}

func (c *clientImpl) Disconnect(conn *Connection) {
	conn.closer <- true
	return
}

func (c *clientImpl) Connect(options *ConnOptions) (*Connection, error) {
	configChan := make(chan model.Config, 0)
	closer := make(chan bool, 0)
	url := c.Options.Url + "/api/config/ws/%s"
	for _, n := range options.Names {
		ws, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf(url, n), nil)
		if err != nil {
			return nil, err
		}

		go func(ws *websocket.Conn, closer chan bool, configChan chan model.Config) {
			for {
				select {
				case <-closer:
					return
				default:
					var conf model.Config
					err := ws.ReadJSON(&conf)
					if err != nil {
						continue
					}

					_ = c.writeToDisk(&conf)
					c.ConfigMap.Set(conf)
					configChan <- conf
				}
			}
		}(ws, closer, configChan)
	}

	conn := Connection{
		closer:   closer,
		OnUpdate: configChan,
	}

	return &conn, nil
}

func (c *clientImpl) Query(name string, query string) (*model.Config, error) {
	conf, err := c.fetchAndLoadConfig(name)
	if err != nil {
		return nil, err
	}
	conf, err = c.query(conf, query)
	if err != nil {
		return nil, err
	}

	return conf, nil
}

func (c *clientImpl) toUrl(name string, query string) string {
	return fmt.Sprintf("%s/api/config/%s?query%s", c.Options.Url, name, query)
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

func (c *clientImpl) fetchAndLoadConfig(name string) (conf *model.Config, err error) {
	exists := false
	conf, exists = c.ConfigMap.Get(name)
	if !exists {
		conf, err = c.fetchConfig(name, "")
		if err != nil {
			conf, err = c.loadFromDisk(name)
			if err != nil {
				return nil, err
			}
		} else {
			err := c.writeToDisk(conf)
			if err != nil {
				return nil, err
			}
		}
	}

	c.ConfigMap.Set(conf)

	return
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

	return &model.Config{
		Name: conf.Data.Name,
		Data: conf.Data.Config,
	}, nil
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
