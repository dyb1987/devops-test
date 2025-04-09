package config

import (
	"fmt"
	"github.com/spf13/viper"
)

var C = &Config{}

type Config struct {
	Jenkins `mapstructure:"jenkins"`
}

type Jenkins struct {
	URL      string `mapstructure:"url"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

func init() {
	viper.SetConfigFile("config.toml")
	viper.AddConfigPath("/etc")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("read config.toml failed: %v\n", err))
	}

	if err = viper.Unmarshal(C); err != nil {
		panic(fmt.Sprintf("unmarshal config.toml to C failed: %v\n", err))
	}
}
