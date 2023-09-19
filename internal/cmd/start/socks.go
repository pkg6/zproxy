package start

import (
	"github.com/pkg6/zproxy/internal/config"
	"github.com/pkg6/zproxy/internal/signal"
	"github.com/pkg6/zproxy/kernel"
	"github.com/pkg6/zproxy/socks"
	"github.com/spf13/cobra"
)

var (
	socksHost string
	socksPort int
	socksAuth string
)

func init() {
	socksCmd.Flags().StringVarP(&socksHost, "host", "H", "127.0.0.1", "Set listening service address")
	socksCmd.Flags().IntVarP(&socksPort, "port", "P", 1080, "Set listening service port")
	socksCmd.Flags().StringVarP(&socksAuth, "auth", "A", "", "Set auth account password 。Example: admin@123456")
}

var socksCmd = &cobra.Command{
	Use:   "socks",
	Short: "Start the socks service",
	Long:  "Start the socks service",
	Run: func(cmd *cobra.Command, args []string) {
		config.RegisterProxyWithConfig(socks.Name, socksHost, socksPort, config.AuthStrToMap(socksAuth))
		kernel.Run()
		signal.Clean()
	},
}
