package main

import "github.com/bancodobrasil/featws-resolver-geolite/cmd"

// "encoding/json"
// "fmt"
// "log"
// "sort"

// adapter "github.com/bancodobrasil/featws-resolver-adapter-go"
// "github.com/bancodobrasil/featws-resolver-adapter-go/types"
// "github.com/bancodobrasil/featws-resolver-geolite/config"

// var cfg = config.Config{}

func main() {
	cmd.Execute()
	// err := config.LoadConfig(&cfg)
	// if err != nil {
	// 	log.Fatalf("Não foi possível carregar as configurações: %s\n", err)
	// }
	// adapter.Run(resolver, adapter.Config{
	// 	Port: cfg.Port,
	// })

}

// func resolver(resolveInput types.ResolveInput, resolveOutput *types.ResolveOutput) {
// 	sort.Strings(resolveInput.Load)
// 	if contains(resolveInput.Load, "geoip") {
// 		resolveGeoIp(resolveInput, resolveOutput)
// 	}
// }

// func resolveGeoIp(resolveInput types.ResolveInput, output *types.ResolveOutput) {
// 	remoteIP, ok := resolveInput.Context["remote_ip"]
// 	fmt.Printf("Finding data for remote ip %s", remoteIP)
// 	if !ok {
// 		output.Errors["geoip"] = "The context 'remote_ip' is required to resolve 'geoip'"
// 	} else {

// 		body := `{"city":"Brasilia", "country":"DF"}`
// 		result := make(map[string]interface{})
// 		json.Unmarshal([]byte(body), &result)
// 		output.Context["geoip"] = result
// 	}
// }

// func contains(s []string, searchterm string) bool {
// 	i := sort.SearchStrings(s, searchterm)
// 	return i < len(s) && s[i] == searchterm
// }
