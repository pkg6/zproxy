package kernel

import (
	"log"
	"sync"
)

var proxyNameMap = map[string]IProxy{}

// Register 注册代理
func Register(name string, proxy IProxy) {
	proxyNameMap[name] = proxy
}

// Run 批量执行服务
func Run(names ...string) {
	var wg sync.WaitGroup
	for name, proxy := range proxyNameMap {
		for _, s := range names {
			if name != s {
				log.Printf("Skipping service %s", name)
				continue
			}
		}
		wg.Add(1)
		go func(name string, proxy IProxy) {
			defer wg.Done()
			if err := proxy.Start(); err != nil {
				log.Printf("%s servcie fail, ERR: %s", name, err)
			}
		}(name, proxy)
	}
	wg.Wait()
}
