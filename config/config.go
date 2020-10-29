package config

import (
	"strings"

	"github.com/spf13/viper"
)

// Config define config contract
type Config interface {
	Init()
	GetString(key string) string
	GetInt(key string) int
	GetBool(key string) bool
}

type viperConfig struct{}

func (v *viperConfig) Init() {
	viper.SetEnvPrefix("awssh")
	viper.AutomaticEnv()

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("$HOME/awssh")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func (v *viperConfig) GetString(key string) string {
	return viper.GetString(key)
}

func (v *viperConfig) GetInt(key string) int {
	return viper.GetInt(key)
}

func (v *viperConfig) GetBool(key string) bool {
	return viper.GetBool(key)
}

// NewViperConfig return new viper config instance
func NewViperConfig() Config {
	v := &viperConfig{}
	v.Init()
	return v
}
