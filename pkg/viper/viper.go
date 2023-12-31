package viper

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Viper *viper.Viper
}

func InitConf(name string) Config {
	config := Config{Viper: viper.New()}
	v := config.Viper
	v.SetConfigType("yaml")
	v.SetConfigName(name)
	v.AddConfigPath("../config/")
	v.AddConfigPath("../../config/")
	v.AddConfigPath("../../../config/")
	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("InitConfig error: %v", err)
	}
	return config
}
