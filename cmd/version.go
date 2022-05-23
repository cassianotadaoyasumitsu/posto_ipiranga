package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var version = "development"
var goVersion = fmt.Sprintf("%s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH)
var buildStamp = ""

const (
	codeName     = "posto-ipiranga-go-kit"
	ReadableName = "Posto Ipiranga"
	figlet       = `                             
__________               __           .___       .__                                    
\______   \____  _______/  |_  ____   |   |_____ |__|___________    ____    _________   
 |     ___/  _ \/  ___/\   __\/  _ \  |   \____ \|  \_  __ \__  \  /    \  / ___\__  \  
 |    |  (  <_> )___ \  |  | (  <_> ) |   |  |_> >  ||  | \// __ \|   |  \/ /_/  > __ \_
 |____|   \____/____  > |__|  \____/  |___|   __/|__||__|  (____  /___|  /\___  (____  /
                    \/                    |__|                  \/     \//_____/     \/ 

                                                                              %s %s
`
)

type versionCommand struct {
	*cobra.Command
}

func NewVersionCommand() *versionCommand {
	vc := &versionCommand{}
	vc.Command = &cobra.Command{
		Use:   "version",
		Short: "Print the Version information",
		Long: fmt.Sprintf(`%s
Print the Version information`, Art()),
		Run: vc.run,
	}
	return vc
}

func (vc *versionCommand) run(cmd *cobra.Command, args []string) {
	cmd.Println(Art())
	cmd.Println(codeName)
	cmd.Println(fmt.Sprintf(" version: %s", version))
	cmd.Println(fmt.Sprintf(" go: %s", goVersion))
	cmd.Println(fmt.Sprintf(" built at: %s", buildStamp))
}

func Art() string {
	return fmt.Sprintf(figlet, ReadableName, version)
}
