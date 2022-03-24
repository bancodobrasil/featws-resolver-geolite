package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"net"
	"os"
	"strings"
	"time"

	"github.com/oschwald/geoip2-golang"
)

type GeoIPDatabase struct {
	GeoDb     *geoip2.Reader
	CityState map[string]map[string]string
}

type GeoRecordLocation struct {
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	AccuracyRadius uint16  `json:"accuracy_radius"`
}

type GeoRecord struct {
	RemoteIp string            `json:"remote_ip"`
	Country  string            `json:"country"`
	City     string            `json:"city"`
	State    string            `json:"state"`
	Location GeoRecordLocation `json:"location"`
}

func NewDatabase(config *Config) (*GeoIPDatabase, error) {

	getLite2DbFile := config.GeoLite2DB
	cityStateDbFile := config.CityStateDB

	geoIPDatabase := GeoIPDatabase{}

	if getLite2DbFile == "" {
		logger.Infof("Geolite database file not found. Localization capabilities based on IP will be disabled")
	} else {
		logger.Debugf("Loading GeoIP2 database %s", getLite2DbFile)
		gdb, err := geoip2.Open(getLite2DbFile)
		if err != nil {
			logger.Fatal(err)
		}
		geoIPDatabase.GeoDb = gdb

		logger.Infof("GeoIP2 database loaded")

		if cityStateDbFile == "" {
			logger.Infof("City State csv file not defined. City input won't be available")
		} else {
			logger.Debugf("Loading City State CSV file %s", cityStateDbFile)
			csvFile, err := os.Open(cityStateDbFile)
			if err != nil {
				logger.Fatal(err)
			}
			geoIPDatabase.CityState = make(map[string]map[string]string)
			reader := csv.NewReader(bufio.NewReader(csvFile))
			for {
				line, err := reader.Read()
				if err == io.EOF {
					break
				} else if err != nil {
					logger.Fatal(err)
				}
				country := strings.ToLower(line[0])
				city := strings.ToLower(line[1])
				state := line[2]
				cm, exists := geoIPDatabase.CityState[country]
				if !exists {
					cm = make(map[string]string)
					geoIPDatabase.CityState[country] = cm
				}
				cm[city] = state
			}
			logger.Infof("City State CSV loaded")
		}
	}
	return &geoIPDatabase, nil
}

func (g *GeoIPDatabase) Find(ipStr string) (*GeoRecord, error) {
	ip := net.ParseIP(ipStr)
	start := time.Now()
	ipRecord, err := g.GeoDb.City(ip)
	logger.Debugf("Time to find getIp data: %s", time.Since(start))
	if err != nil {
		logger.Debugf("Couldn't find geo info for ip %s. err=%s", ipStr, err)
		return nil, err
	} else {
		geoRecord := &GeoRecord{
			Country: ipRecord.Country.Names["en"],
			City:    ipRecord.City.Names["en"],
			Location: GeoRecordLocation{
				Latitude:       ipRecord.Location.Latitude,
				Longitude:      ipRecord.Location.Longitude,
				AccuracyRadius: ipRecord.Location.AccuracyRadius,
			},
		}
		cs, exists := g.CityState[strings.ToLower(ipRecord.Country.IsoCode)]
		if exists {
			state, exists := cs[strings.ToLower(ipRecord.City.Names["en"])]
			if exists {
				geoRecord.State = state
			}
		}
		return geoRecord, nil
	}

}

func (g *GeoIPDatabase) Close() error {
	err := g.GeoDb.Close()
	return err
}
