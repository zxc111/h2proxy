package main

import (
	"fmt"
	"github.com/zxc111/h2proxy"
	"net/http"

	_ "net/http/pprof"
)

func main() {

	category, conf := h2proxy.ParseConfig()

	h2proxy.InitLogger()
	var debugPort int

	var server h2proxy.H2proxy
	switch category {
	case h2proxy.HTTP:
		config := conf.(*h2proxy.ClientConfig)
		server = h2proxy.HttpProxy{Config: config}
		debugPort = config.Pprof
	case h2proxy.SOCKSV5:
		config := conf.(*h2proxy.ClientConfig)
		server = h2proxy.Sock5Proxy{Config: config}
		debugPort = config.Pprof

	case h2proxy.SERVER:
		config := conf.(*h2proxy.ServerConfig)
		server = h2proxy.Http2Server{Config: config}
		debugPort = config.Pprof
	}
	go startPProf(debugPort)

	server.Start()
}

func startPProf(port int) {

	addr := fmt.Sprintf("localhost:%d", port)
	h2proxy.Log.Info("pprof is running at " + addr)
	h2proxy.Log.Fatal(http.ListenAndServe(addr, nil))
}
