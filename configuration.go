package main

import (
	"github.com/spf13/viper"
)

type Config struct {
	Port string

	GeoLite2DB  string
	CityStateDB string
}

func LoadConfig() *Config {
	return &Config{
		Port: viper.GetString("server.port"),

		GeoLite2DB:  viper.GetString("database.geolite2"),
		CityStateDB: viper.GetString("database.citystate"),
	}

}
