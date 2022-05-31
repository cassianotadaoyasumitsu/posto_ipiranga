package main

import (
	"fmt"
	"os"

	"git.wealth-park.com/cassiano/posto_ipiranga/cmd"
	"github.com/spf13/cobra"
)

type rootCommand struct {
	*cobra.Command
}

func newRootCommand() *rootCommand {
	rc := &rootCommand{}
	rc.Command = &cobra.Command{
		Use:   "posto-ipirangad",
		Short: fmt.Sprintf("Command Line Interface manager for %s", cmd.ReadableName),
		Long: fmt.Sprintf(`%s
Manage %s from the command line`, cmd.Art(), cmd.ReadableName),
		Run:           rc.run,
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	return rc
}

func (rc *rootCommand) run(c *cobra.Command, args []string) {
	c.Help()
}

func main() {
	root := newRootCommand()
	root.Command.AddCommand(
		cmd.NewVersionCommand().Command,
		cmd.NewHTTPCommand().Command,
	)
	c, err := root.Command.ExecuteC()
	if err != nil {
		c.Println(cmd.Art())
		c.Println(c.UsageString())
		c.PrintErrf("ERROR: %v\n", err)
		os.Exit(1)
	}
}
