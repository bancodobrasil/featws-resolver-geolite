package databases

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strings"

	"github.com/bancodobrasil/featws-resolver-geolite/configuration"
	"github.com/oschwald/geoip2-golang"
	"github.com/sirupsen/logrus"
)

type GeoIPDatabase struct {
	GeoDb     *geoip2.Reader
	CityState map[string]map[string]string
}

func NewDatabase(config *configuration.Config) (*GeoIPDatabase, error) {

	getLite2DbFile := config.GeoLite2DB
	cityStateDbFile := config.CityStateDB

	geoIPDatabase := GeoIPDatabase{}

	if getLite2DbFile == "" {
		logrus.Infof("Geolite database file not found. Localization capabilities based on IP will be disabled")
	} else {
		logrus.Debugf("Loading GeoIP2 database %s", getLite2DbFile)
		gdb, err := geoip2.Open(getLite2DbFile)
		if err != nil {
			log.Fatal(err)
		}
		geoIPDatabase.GeoDb = gdb

		logrus.Infof("GeoIP2 database loaded")

		if cityStateDbFile == "" {
			logrus.Infof("City State csv file not defined. _ip_state input won't be available")
		} else {
			logrus.Debugf("Loading City State CSV file %s", cityStateDbFile)
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
			logrus.Infof("City State CSV loaded")
		}
	}
	return &geoIPDatabase, nil
}

func (g *GeoIPDatabase) Close() error {
	err := g.GeoDb.Close()
	return err
}
