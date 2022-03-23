package resolver

import (
	"log"

	adapter "github.com/bancodobrasil/featws-resolver-adapter-go"
	"github.com/bancodobrasil/featws-resolver-geolite/configuration"
)

func Init() {
	config := configuration.LoadConfig()
	log.Println("configuring resolver...")
	adapter.Run(resolver, adapter.Config{
		Port: config.Port,
	})
}
