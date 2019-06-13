package configuration

import (
	"bytes"
	"github.com/eddieowens/axon"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"path"
	"runtime"
	"strings"
)

const ConfigKey = "Config"

type Config struct {
	Env      string   `mapstructure:"env"`
	Git      Git      `mapstructure:"git"`
	Compiler Compiler `mapstructure:"compiler"`
}

type Git struct {
	Remote          string `mapstructure:"remote"`
	Branch          string `mapstructure:"branch"`
	Directory       string `mapstructure:"directory"`
	PollingInterval int    `mapstructure:"polling_interval"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	SSHKey          string `mapstructure:"ssh_key"`
}

type Compiler struct {
	OutputDirectory string `mapstructure:"output_directory"`
}

func defaultConfig() *Config {
	return &Config{
		Env: "dev",
	}
}

func configFactory(_ axon.Injector, _ axon.Args) axon.Instance {
	v := viper.New()
	v.SetConfigType("yaml")

	b, _ := yaml.Marshal(defaultConfig())
	defaultConfig := bytes.NewReader(b)
	if err := v.MergeConfig(defaultConfig); err != nil {
		panic(err)
	}

	_, filename, _, _ := runtime.Caller(0)
	d, _ := path.Split(filename)

	v.AddConfigPath(path.Join(d, "..", "..", "config"))
	v.AddConfigPath("./config")

	v.SetConfigName("config")
	if err := v.MergeInConfig(); err != nil {
		panic(err)
	}

	v.SetConfigName(v.GetString("env"))
	if err := v.MergeInConfig(); err != nil {
		log.Warn(err)
	}

	v.AutomaticEnv()
	v.SetEnvPrefix("ranvier")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AllowEmptyEnv(false)

	config := Config{}
	if err := v.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}

	return axon.Any(config)
}
