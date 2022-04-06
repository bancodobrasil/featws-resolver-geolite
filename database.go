package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"

	"github.com/oschwald/geoip2-golang"
	log "github.com/sirupsen/logrus"
)

// GeoIPDatabase ...
type GeoIPDatabase struct {
	GeoDb     *geoip2.Reader
	CityState map[string]map[string]string
}

// GeoRecordLocation ...
type GeoRecordLocation struct {
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	AccuracyRadius uint16  `json:"accuracy_radius"`
}

// GeoRecord ...
type GeoRecord struct {
	RemoteIP string            `json:"remote_ip"`
	Country  string            `json:"country"`
	City     string            `json:"city"`
	State    string            `json:"state,omitempty"`
	Location GeoRecordLocation `json:"location"`
}

// NewDatabase ...
func NewDatabase(config *Config) (*GeoIPDatabase, error) {

	getLite2DbFile := config.GeoLite2DB
	cityStateDbFile := config.CityStateDB

	geoIPDatabase := GeoIPDatabase{}

	if getLite2DbFile == "" {
		return nil, errors.New("Geolite database file not defined")
	}
	log.Debugf("Loading GeoIP2 database %s", getLite2DbFile)
	gdb, err := geoip2.Open(getLite2DbFile)
	if err != nil {
		return nil, err
	}
	geoIPDatabase.GeoDb = gdb

	log.Infof("GeoIP2 database loaded")

	if cityStateDbFile == "" {
		log.Infof("City State csv file not defined. City input won't be available")
	} else {
		log.Debugf("Loading City State CSV file %s", cityStateDbFile)
		csvFile, err := os.Open(cityStateDbFile)
		if err != nil {
			log.Fatal(err)
		}
		geoIPDatabase.CityState = make(map[string]map[string]string)
		reader := csv.NewReader(bufio.NewReader(csvFile))
		for {
			line, err := reader.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
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
		log.Infof("City State CSV loaded")
	}
	return &geoIPDatabase, nil
}

// Find ...
func (g *GeoIPDatabase) Find(ipStr string) (*GeoRecord, error) {
	if g.GeoDb != nil {
		ip := net.ParseIP(ipStr)
		start := time.Now()
		ipRecord, err := g.GeoDb.City(ip)
		log.Debugf("Time to find getIp data: %s", time.Since(start))
		if err != nil {
			errStr := fmt.Sprintf("Couldn't find geo info for ip %s. err=%s", ipStr, err)
			log.Debugf(errStr)
			return nil, errors.New(errStr)
		}
		geoRecord := &GeoRecord{
			RemoteIP: ipStr,
			Country:  ipRecord.Country.Names["en"],
			City:     ipRecord.City.Names["en"],
			Location: GeoRecordLocation{
				Latitude:       ipRecord.Location.Latitude,
				Longitude:      ipRecord.Location.Longitude,
				AccuracyRadius: ipRecord.Location.AccuracyRadius,
			},
		}
		if g.CityState != nil {
			cs, exists := g.CityState[strings.ToLower(ipRecord.Country.IsoCode)]
			if exists {
				state, exists := cs[strings.ToLower(ipRecord.City.Names["en"])]
				if exists {
					geoRecord.State = state
				}
			}
		}
		return geoRecord, nil
	}
	return nil, errors.New("No geolite database")
}

// Close ...
func (g *GeoIPDatabase) Close() error {
	err := g.GeoDb.Close()
	return err
}
