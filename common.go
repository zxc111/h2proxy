package h2proxy

import (
	"crypto/tls"
	"flag"
	"fmt"
	"golang.org/x/net/http2"
	"log"
	"net"
	"os"
)

func NewTransport(proxyAddr string) *http2.Transport {
	return &http2.Transport{
		DialTLS: func(network, addr string, config *tls.Config) (net.Conn, error) {
			return tls.Dial("tcp", proxyAddr, &tls.Config{
				NextProtos:         []string{"h2"},
				InsecureSkipVerify: true,
			})
		},
		AllowHTTP: true,
	}
}

func ParseConfig() (local, proxy string) {
	var (
		localHost string
		localPort string
		proxyHost string
		proxyPort string
	)
	flag.StringVar(&localHost, "local_host", "localhost", "-local_host=127.0.0.1")
	flag.StringVar(&localPort, "local_port", "3002", "-local_port=4000")
	flag.StringVar(&proxyHost, "proxy_host", "", "-porxy_host=xxx.xxx.xxx.xxx")
	flag.StringVar(&proxyPort, "proxy_port", "", "-proxy_port=3000")

	flag.Parse()
	if proxyHost == "" {
		flag.Usage()
		log.Println("proxy_host is required")
		os.Exit(1)
	}
	if proxyPort == "" {
		log.Println("proxy_port is required")
		flag.Usage()

		os.Exit(1)
	}
	if localHost == "" {
		log.Println("local_host is required")
		flag.Usage()

		os.Exit(1)
	}
	if localPort == "" {
		log.Println("local_port is requred")
		flag.Usage()

		os.Exit(1)
	}
	proxy = fmt.Sprintf("%s:%s", proxyHost, proxyPort)
	local = fmt.Sprintf("%s:%s", localHost, localPort)
	log.Printf("local: %s", local)
	log.Printf("proxy: %s", proxy)

	return local, proxy
}
