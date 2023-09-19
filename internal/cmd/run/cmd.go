package run

import (
	"github.com/pkg6/zproxy/http"
	"github.com/pkg6/zproxy/internal/config"
	"github.com/pkg6/zproxy/internal/signal"
	"github.com/pkg6/zproxy/kernel"
	"github.com/pkg6/zproxy/socks"
	"github.com/spf13/cobra"
	"log"
)

var (
	runConfigFile string
	cfg           Config
)

type Config struct {
	Host      string
	HTTPPort  int
	SocksPort int
	Auth      map[string]string
}

func init() {
	Cmd.Flags().StringVarP(&runConfigFile, "config", "C", "", "Set Config File")
}

var Cmd = &cobra.Command{
	Use:   "run",
	Short: "Start with Config File",
	Long:  "Start with Config File",
	Run: func(cmd *cobra.Command, args []string) {
		config.LoadConfigFile(runConfigFile, &cfg)
		if cfg.SocksPort == cfg.HTTPPort {
			log.Fatalln("Ports cannot be the same")
			return
		}
		isRun := false
		if cfg.SocksPort != 0 {
			config.RegisterProxyWithConfig(socks.Name, cfg.Host, cfg.SocksPort, cfg.Auth)
			isRun = true
		}
		if cfg.HTTPPort != 0 {
			config.RegisterProxyWithConfig(http.Name, cfg.Host, cfg.HTTPPort, cfg.Auth)
			isRun = true
		}
		if isRun {
			kernel.Run()
			signal.Clean()
			return
		}
		log.Fatalln("Run failed")
		return
	},
}
