package main

import (
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
		Short: "Starts resolver",
		Long:  `Starts the resolver's server to handle ruller requests`,
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
	rootCmd.PersistentFlags().StringP("geo-database-file", "g", "", "Geolite database file path")
	rootCmd.PersistentFlags().StringP("cities-database-file", "c", "", "Cities database file path")
	rootCmd.PersistentFlags().StringP("server-port", "p", "", "Server port listening")
	rootCmd.PersistentFlags().BoolP("log-json", "j", false, "Log format json (default is false)")
	rootCmd.PersistentFlags().StringP("log-level", "l", "error", "Log level (default is error)")
	viper.BindPFlag("DATABASE_GEOLITE2", rootCmd.PersistentFlags().Lookup("geo-database-file"))
	viper.BindPFlag("DATABASE_CITYSTATE", rootCmd.PersistentFlags().Lookup("cities-database-file"))
	viper.BindPFlag("SERVER_PORT", rootCmd.PersistentFlags().Lookup("server-port"))
	viper.BindPFlag("LOG_JSON", rootCmd.PersistentFlags().Lookup("log-json"))
	viper.BindPFlag("LOG_LEVEL", rootCmd.PersistentFlags().Lookup("log-level"))
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
		viper.SetConfigType("env")
		viper.SetConfigName(".env")
	}

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()
	InitLogger()
	if err := viper.ReadInConfig(); err == nil {
		logger.Infof("Using config file: %s", viper.ConfigFileUsed())
	}

}
