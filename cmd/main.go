package main

import (
	"github.com/zxc111/h2proxy"
	"log"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	category, conf := h2proxy.ParseConfig()
	switch category {
	case h2proxy.HTTP:
		config := conf.(*h2proxy.ClientConfig)
		startHttp(config)
	case h2proxy.SOCKSV5:
		config := conf.(*h2proxy.ClientConfig)
		startSocks5(config)
	case h2proxy.SERVER:
		config := conf.(*h2proxy.ServerConfig)
		StartServer(config)
	}
}
