package cmd

import (
	"github.com/pkg6/zproxy/internal/cmd/run"
	"github.com/pkg6/zproxy/internal/cmd/start"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "zproxy",
	Short: "Welcome to zproxy",
	Long:  `Welcome to zproxy`,
}

func init() {
	rootCmd.AddCommand(start.Cmd)
	rootCmd.AddCommand(run.Cmd)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
