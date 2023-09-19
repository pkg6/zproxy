package config

import (
	"github.com/pkg6/zproxy/http"
	"github.com/pkg6/zproxy/kernel"
	"github.com/pkg6/zproxy/socks"
	"github.com/spf13/viper"
	"strings"
)

var (
	Vip    *viper.Viper
	proxys = map[string]kernel.IProxy{
		socks.Name: &socks.Proxy{},
		http.Name:  &http.Proxy{},
	}
)

func init() {
	Vip = viper.New()
	Vip.SetConfigName("config")
	Vip.SetConfigType("yaml")
	Vip.AddConfigPath("./")
	Vip.AddConfigPath("/etc/zproxy/")
}

func LoadConfigFile(configFile string, v any) {
	if configFile != "" {
		Vip.SetConfigFile(configFile)
	}
	_ = Vip.ReadInConfig()
	_ = Vip.Unmarshal(v)
}

func RegisterProxyWithConfig(name, host string, port int, auth map[string]string) {
	if proxy, ok := proxys[name]; ok {
		proxy.SetIP(host)
		proxy.SetPort(port)
		proxy.SetAuth(auth)
		kernel.Register(name, proxy)
	}
}

func AuthStrToMap(auth string) map[string]string {
	maps := map[string]string{}
	authStr := strings.Split(auth, "@")
	if len(authStr) >= 2 {
		maps[authStr[0]] = authStr[1]
	}
	return maps
}
