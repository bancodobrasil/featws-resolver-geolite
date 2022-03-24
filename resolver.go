package main

import (
	"fmt"
	"sort"

	adapter "github.com/bancodobrasil/featws-resolver-adapter-go"
	"github.com/bancodobrasil/featws-resolver-adapter-go/types"
)

var geoIPDatabase *GeoIPDatabase

func Init() {
	config := LoadConfig()
	logger.Info("Recovering geo database...")
	db, err := NewDatabase(config)

	if err != nil {
		logger.Fatalf("Couldn't connect to database. Error: %s", err)
	}
	geoIPDatabase = db
	logger.Info("configuring resolver...")
	adapter.Run(resolver, adapter.Config{
		Port: config.Port,
	})
}

func resolver(resolveInput types.ResolveInput, resolveOutput *types.ResolveOutput) {
	sort.Strings(resolveInput.Load)
	if contains(resolveInput.Load, "geoip") {
		resolveGeoIp(resolveInput, resolveOutput)
	}
}

func resolveGeoIp(resolveInput types.ResolveInput, output *types.ResolveOutput) {
	remoteIP, ok := resolveInput.Context["remote_ip"]
	logger.Debugf("Finding data for remote ip %s", remoteIP)
	if !ok {
		output.Errors["geoip"] = "The context 'remote_ip' is required to resolve 'geoip'"
	} else {
		remoteIPStr := fmt.Sprintf("%v", remoteIP)
		geoRecord, err := geoIPDatabase.Find(remoteIPStr)
		geoRecord.RemoteIp = remoteIPStr
		if err != nil {
			output.Errors["geoip"] = err
		} else {
			output.Context["geoip"] = geoRecord
		}
	}
}

func contains(s []string, searchterm string) bool {
	i := sort.SearchStrings(s, searchterm)
	return i < len(s) && s[i] == searchterm
}
