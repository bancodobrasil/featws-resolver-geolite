package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidDatabase(t *testing.T) {
	config := &Config{
		GeoLite2DB: "test-data/test-data/GeoLite2-City-Test.mmdb",
	}
	db, err := NewDatabase(config)
	assert.NoError(t, err)
	ipStr := "81.2.69.142"
	geoRecord, err := db.Find(ipStr)
	assert.NoError(t, err)
	assert.Equal(t, ipStr, geoRecord.RemoteIP)
	assert.Equal(t, "United Kingdom", geoRecord.Country)
	assert.Equal(t, "London", geoRecord.City)
	assert.Equal(t, float64(51.5142), geoRecord.Location.Latitude)
	assert.Equal(t, float64(-0.0931), geoRecord.Location.Longitude)
	assert.Equal(t, uint16(10), geoRecord.Location.AccuracyRadius)
	db.Close()
}

func TestDatabaseNotDefined(t *testing.T) {
	config := &Config{
		GeoLite2DB: "",
	}
	db, err := NewDatabase(config)
	assert.Error(t, err, "Geolite database file not defined")
	assert.Empty(t, db)
}

func TestDatabaseInvalid(t *testing.T) {
	config := &Config{
		GeoLite2DB: "database",
	}
	db, err := NewDatabase(config)
	assert.Error(t, err, "open database: no such file or directory")
	assert.Empty(t, db)
}
