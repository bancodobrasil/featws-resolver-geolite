package main

import (
	"sort"

	adapter "github.com/bancodobrasil/featws-resolver-adapter-go"
	"github.com/bancodobrasil/featws-resolver-adapter-go/types"
	log "github.com/sirupsen/logrus"
)

var geoIPDatabase *GeoIPDatabase

// Init ...
func Init() {
	config := LoadConfig()
	log.Info("Recovering geo database...")
	db, err := NewDatabase(config)

	if err != nil {
		log.Fatalf("Couldn't connect to database. Error: %s", err)
	}
	geoIPDatabase = db
	log.Infof("Running resolver on port %s ...", config.Port)
	adapter.Run(resolver, adapter.Config{
		Port: config.Port,
	})
}

func resolver(resolveInput types.ResolveInput, resolveOutput *types.ResolveOutput) {
	sort.Strings(resolveInput.Load)
	if contains(resolveInput.Load, "geoip") {
		resolveGeoIP(resolveInput, resolveOutput)
	}
	if contains(resolveInput.Load, "uf") {
		resolveUF(resolveInput, resolveOutput)
	}
}

func resolveGeoIP(resolveInput types.ResolveInput, output *types.ResolveOutput) {
	remoteIP, ok := resolveInput.Context["remote_ip"]
	log.Debugf("Finding data for remote ip %s", remoteIP)
	if !ok {
		output.Errors["geoip"] = "The context 'remote_ip' is required to resolve 'geoip'"
	} else {
		geoRecord, err := geoIPDatabase.Find(remoteIP.(string))
		if err != nil {
			output.Errors["geoip"] = err.Error()
		} else {
			output.Context["geoip"] = geoRecord
		}
	}
}

func resolveUF(resolveInput types.ResolveInput, output *types.ResolveOutput) {
	remoteIP, ok := resolveInput.Context["remote_ip"]
	log.Debugf("Finding data for remote ip %s", remoteIP)
	if !ok {
		output.Errors["uf"] = "The context 'remote_ip' is required to resolve 'geoip'"
	} else {
		geoRecord, err := geoIPDatabase.Find(remoteIP.(string))
		if err != nil {
			output.Errors["uf"] = err.Error()
		} else {
			output.Context["uf"] = geoRecord.State
		}
	}
}

func contains(s []string, searchterm string) bool {
	if len(s) == 0 {
		return true
	}
	i := sort.SearchStrings(s, searchterm)
	return i < len(s) && s[i] == searchterm
}
