package main

import (
	"fmt"
	"github.com/zxc111/h2proxy"
	"log"
	"net/http"

	_ "net/http/pprof"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {

	category, conf := h2proxy.ParseConfig()

	var debugPort int

	var server h2proxy.H2proxy
	switch category {
	case h2proxy.HTTP:
		config := conf.(*h2proxy.ClientConfig)
		server = h2proxy.HttpProxy{config}
		debugPort = config.DebugPort
	case h2proxy.SOCKSV5:
		config := conf.(*h2proxy.ClientConfig)
		server = h2proxy.Sock5Proxy{config}
		debugPort = config.DebugPort

	case h2proxy.SERVER:
		config := conf.(*h2proxy.ServerConfig)
		server = h2proxy.Http2Server{config}
		debugPort = config.DebugPort
	}
	go startPProf(debugPort)
	server.Start()
}

func startPProf(port int) {

	addr := fmt.Sprintf("localhost:%d", port)
	log.Printf("pprof is running at %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
