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
	case "http":
		config := conf.(*h2proxy.ClientConfig)
		startHttp(config)
	case "socks5":
		config := conf.(*h2proxy.ClientConfig)
		startSocks5(config)
	case "server":
		config := conf.(*h2proxy.ServerConfig)
		StartServer(config)
	}
}
