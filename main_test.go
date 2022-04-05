package main

import (
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func setup() {
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
