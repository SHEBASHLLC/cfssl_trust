package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cloudflare/cfssl/log"
	"github.com/cloudflare/cfssl_trust/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	dbFile  string
)

func root(cmd *cobra.Command, args []string) {
}

var configLocations = []string{
	"/etc/cfssl",
	"/usr/local/cfssl",
	filepath.Join(config.GoPath(), "src", "github.com", "cloudflare", "cfssl_trust"),
}

var RootCmd = &cobra.Command{
	Use:   "certmgr",
	Short: "Manage TLS certificates for multiple services",
	Long:  ``,
	Run:   root,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "f", "", "config file (default is /etc/cfssl/cfssl.yaml)")
	RootCmd.PersistentFlags().StringVarP(&dbFile, "db", "d", "", "path to trust database")

	viper.BindPFlag("database.path", RootCmd.PersistentFlags().Lookup("db"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("cfssl-trust") // name of config file (without extension)
		for _, dir := range configLocations {
			viper.AddConfigPath(dir)
		}
		viper.AddConfigPath(".")
	}

	viper.SetEnvPrefix("CFSSL_TRUST")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err == nil {
		log.Info("cfssl-trust: loading from config file ", viper.ConfigFileUsed())
	}
}