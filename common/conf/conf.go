package config

import (
	"fmt"
	json "github.com/json-iterator/go"
	"github.com/spf13/viper"
	"walrus_llm_project/log"
)

type Configuration struct {
	Env    string         `mapstructure:"env"`
	Server Server         `mapstructure:"server"`
	Log    *log.LogConfig `mapstructure:"log"`
}

func (c *Configuration) String() string {
	s, _ := json.MarshalToString(c)
	return s
}

type Server struct {
	RunMode  string `mapstructure:"run_mode"`
	HttpPort int    `mapstructure:"http_port"`
}

var conf = &Configuration{}

func GetConfig() *Configuration {
	return conf
}

var (
	configViper = viper.New()
)

func InitConfig(path string) *Configuration {
	fmt.Printf("init config path: %s\n", path)
	configViper.SetConfigFile(path)

	err := configViper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	setDefault(configViper)

	err = configViper.Unmarshal(&conf)
	if err != nil {
		panic(err)
	}
	return conf
}

func setDefault(v *viper.Viper) {
	v.SetDefault("log.path", "logs/")
}
