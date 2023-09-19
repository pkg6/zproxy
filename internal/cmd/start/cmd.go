package start

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "start",
	Short: "Start service related commands",
	Long:  "Start service related commands",
}

func init() {
	Cmd.AddCommand(socksCmd)
	Cmd.AddCommand(httpCmd)
}
