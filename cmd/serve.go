package cmd

import (
	"github.com/bancodobrasil/featws-resolver-geolite/resolver"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start resolver",
	Long:  `Starts the resolver and serves rulles related`,
	Run: func(cmd *cobra.Command, args []string) {
		resolver.Init()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
