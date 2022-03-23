package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use:   "resolver",
		Short: "Resolver Geolite is a resolver adapter component",
		Long: `Resolver Geolite is a resolver adapter component
created with golang employed in the FeatWS ecosystem.`,
	}
	serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "start resolver",
		Long:  `Starts the resolver and serves rulles related`,
		Run: func(cmd *cobra.Command, args []string) {
			Init()
		},
	}
)

func Execute() {
	err := rootCmd.Execute()
	cobra.CheckErr(err)
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	rootCmd.PersistentFlags().Bool("debug", false, "debug server")
	viper.BindPFlag("log.debug-mode", rootCmd.PersistentFlags().Lookup("debug"))
	rootCmd.AddCommand(serveCmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		ex, err := os.Executable()
		cobra.CheckErr(err)
		exePath := filepath.Dir(ex)
		viper.AddConfigPath(exePath)
		viper.SetConfigType("toml")
		viper.SetConfigName(".config")
	}

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file: ", viper.ConfigFileUsed())
	}

}
