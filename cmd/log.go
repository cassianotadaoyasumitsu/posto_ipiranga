package cmd

import (
	"git.wealth-park.com/cassiano/posto_ipiranga/internal/log"
	"git.wealth-park.com/cassiano/posto_ipiranga/internal/log/zerolog"
	"github.com/spf13/cobra"
)

func initLogger(cmd *cobra.Command, conf *baseConfig) log.Logger {
	lvl := log.InfoLevel
	if conf.Debug {
		lvl = log.DebugLevel
	}
	cmd.Println(Art())
	if conf.JSONLogging {
		return zerolog.NewJSONLogger(lvl)
	}
	return zerolog.NewConsoleLogger(lvl)
}
