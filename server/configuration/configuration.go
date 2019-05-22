package configuration

import (
	"fmt"
	"github.com/eddieowens/axon"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
	"os"
)

const ConfigKey = "Config"

type Config struct {
	Env string    `mapstructure:"env"`
	Git GitConfig `mapstructure:"git"`
}

type GitConfig struct {
	Remote          string `mapstructure:"remote"`
	Branch          string `mapstructure:"branch"`
	Directory       string `mapstructure:"directory"`
	PollingInterval int    `mapstructure:"polling_interval"`
}

func configFactory(_ axon.Args) axon.Instance {
	v := viper.New()

	v.AutomaticEnv()
	v.SetEnvPrefix("cubby")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	v.SetConfigName("config")
	if err := v.ReadInConfig(); err != nil {
		fmt.Println(err)
	}

	env := os.Getenv("ENV")
	if env == "" {
		v.SetConfigName("dev")
		v.Set("env", "dev")
	} else {
		v.SetConfigName(os.Getenv("ENV"))
		v.Set("env", env)
	}
	if err := v.MergeInConfig(); err != nil {
		log.Fatal(err)
	}

	config := Config{}
	if err := v.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}

	return axon.Any(config)
}
