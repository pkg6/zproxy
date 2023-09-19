package start

import (
	"github.com/pkg6/zproxy/http"
	"github.com/pkg6/zproxy/internal/config"
	"github.com/pkg6/zproxy/internal/signal"
	"github.com/pkg6/zproxy/kernel"
	"github.com/spf13/cobra"
)

var (
	httpHost string
	httpPort int
	httpAuth string
)

func init() {
	httpCmd.Flags().StringVarP(&httpHost, "host", "H", "127.0.0.1", "Set listening service address")
	httpCmd.Flags().IntVarP(&httpPort, "port", "P", 1080, "Set listening service port")
	httpCmd.Flags().StringVarP(&httpAuth, "auth", "A", "", "Set auth account password 。Example: admin@123456")
}

var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "Start the http service",
	Long:  "Start the http service",
	Run: func(cmd *cobra.Command, args []string) {
		config.RegisterProxyWithConfig(http.Name, httpHost, httpPort, config.AuthStrToMap(httpAuth))
		kernel.Run()
		signal.Clean()
	},
}
