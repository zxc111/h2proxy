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

func ParseConfig() (local, proxy string, needAuth bool, user *UserInfo) {
	var (
		localHost string
		localPort string
		proxyHost string
		proxyPort string
	)
	user = &UserInfo{}
	flag.StringVar(&localHost, "local_host", "localhost", "-local_host=127.0.0.1")
	flag.StringVar(&localPort, "local_port", "3002", "-local_port=4000")
	flag.StringVar(&proxyHost, "proxy_host", "", "-porxy_host=xxx.xxx.xxx.xxx")
	flag.StringVar(&proxyPort, "proxy_port", "", "-proxy_port=3000")
	flag.StringVar(&(user.username), "user", "", "-user=abc")
	flag.StringVar(&(user.passwd), "passwd", "", "-passwd=def")
	flag.BoolVar(&needAuth, "need_auth", false, "-need_auth=false")

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
	return local, proxy, needAuth, user
}

func ParseServerConfig() (server, caKey, caCrt string, needAuth bool, user *UserInfo) {
	user = &UserInfo{}
	var (
		host string
		port string
	)
	caKey = ""
	caCrt = ""
	flag.StringVar(&host, "host", "localhost", "-host=0.0.0.0")
	flag.StringVar(&port, "port", "", "-port=3000")
	flag.StringVar(&caKey, "cert_key", "", "-cert_key=/root/test.key")
	flag.StringVar(&caCrt, "cert_crt", "", "-cert_crt=/root/test.crt")
	flag.BoolVar(&needAuth, "need_auth", false, "-need_auth=false")

	flag.Parse()

	if host == "" {
		log.Println("host is required")
		flag.Usage()

		os.Exit(1)
	}
	if port == "" {
		log.Println("port is requred")
		flag.Usage()

		os.Exit(1)
	}

	if caKey == "" {
		log.Println("cert_key is requred")
		flag.Usage()

		os.Exit(1)
	}

	if caCrt == "" {
		log.Println("cert_crt is requred")
		flag.Usage()

		os.Exit(1)
	}

	server = fmt.Sprintf("%s:%s", host, port)
	log.Println(server)
	return server, caKey, caCrt, needAuth, user
}
