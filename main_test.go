package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func setup() {
	logger = logrus.New()
	logger.Out = ioutil.Discard

	config := &Config{
		GeoLite2DB: "test-data/test-data/GeoLite2-City-Test.mmdb",
	}
	db, err := NewDatabase(config)
	if err != nil {
		log.Fatal("Geolite2-Test database not found!")
	}
	geoIPDatabase = db
}

func shutdown() {
	geoIPDatabase.Close()
}
