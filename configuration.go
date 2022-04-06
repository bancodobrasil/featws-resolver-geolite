package main

import (
	"github.com/spf13/viper"
)

// Config ...
type Config struct {
	Port string

	GeoLite2DB  string
	CityStateDB string
}

// LoadConfig ...
func LoadConfig() *Config {
	return &Config{
		Port: viper.GetString("SERVER_PORT"),

		GeoLite2DB:  viper.GetString("DATABASE_GEOLITE2"),
		CityStateDB: viper.GetString("DATABASE_CITYSTATE"),
	}

}
