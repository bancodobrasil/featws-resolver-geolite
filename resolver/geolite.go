package resolver

import (
	"encoding/json"
	"log"
	"sort"

	"github.com/bancodobrasil/featws-resolver-adapter-go/types"
)

func resolver(resolveInput types.ResolveInput, resolveOutput *types.ResolveOutput) {
	sort.Strings(resolveInput.Load)
	if contains(resolveInput.Load, "geoip") {
		resolveGeoIp(resolveInput, resolveOutput)
	}
}

func resolveGeoIp(resolveInput types.ResolveInput, output *types.ResolveOutput) {
	remoteIP, ok := resolveInput.Context["remote_ip"]
	log.Printf("Finding data for remote ip %s", remoteIP)
	if !ok {
		output.Errors["geoip"] = "The context 'remote_ip' is required to resolve 'geoip'"
	} else {

		body := `{"city":"Brasilia", "country":"DF"}`
		result := make(map[string]interface{})
		json.Unmarshal([]byte(body), &result)
		output.Context["geoip"] = result
	}
}
