package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type baseConfig struct {
	JSONLogging bool `mapstructure:"json"`
	Debug       bool `mapstructure:"debug"`
}

func unmarshalConfig(c *cobra.Command, envPrefix string, rawVal interface{}) {
	var configFilePath string

	c.Flags().SortFlags = true
	c.Flags().StringVarP(&configFilePath, "config", "c", "", "Path to the configuration file")
	c.Flags().Bool("json", false, "Output the logs in JSON format")
	c.Flags().Bool("debug", false, "Enable debug level logging")

	cobra.OnInitialize(func() {
		viper.SetEnvPrefix(envPrefix)
		viper.AutomaticEnv()
		if configFilePath != "" {
			viper.SetConfigFile(configFilePath)
			if err := viper.ReadInConfig(); err != nil {
				c.PrintErrf("ERROR: failed to read config file: %v\n", err)
				os.Exit(1)
			}
		}
		c.Flags().VisitAll(func(f *pflag.Flag) {
			viper.BindPFlag(strings.Replace(f.Name, "-", "_", -1), f)
		})
		if err := viper.Unmarshal(rawVal); err != nil {
			c.PrintErrf("ERROR: failed to unmarshal the configs: %v\n", err)
			os.Exit(1)
		}
	})
}
