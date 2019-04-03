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

type ServerConfig struct {
	Server   string
	CaKey    string
	CaCrt    string
	NeedAuth bool
	User     *UserInfo
}

type ClientConfig struct {
	Local    string
	Proxy    string
	needAuth bool
	user     *UserInfo
}

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

func ParseConfig() (category string, config interface{}) {
	flag.StringVar(&category, "category", "client", "-category=http/server/socks5")

	user := &UserInfo{}
	var (
		host  string
		port  string
		caKey string
		caCrt string

		needAuth bool

		localHost string
		localPort string
		proxyHost string
		proxyPort string
	)

	// server
	flag.StringVar(&host, "host", "localhost", "-host=0.0.0.0")
	flag.StringVar(&port, "port", "", "-port=3000")
	flag.StringVar(&caKey, "cert_key", "", "-cert_key=/root/test.key")
	flag.StringVar(&caCrt, "cert_crt", "", "-cert_crt=/root/test.crt")

	// client
	flag.StringVar(&localHost, "local_host", "localhost", "-local_host=127.0.0.1")
	flag.StringVar(&localPort, "local_port", "3002", "-local_port=4000")
	flag.StringVar(&proxyHost, "proxy_host", "", "-porxy_host=xxx.xxx.xxx.xxx")
	flag.StringVar(&proxyPort, "proxy_port", "", "-proxy_port=3000")

	// common
	flag.BoolVar(&needAuth, "need_auth", false, "-need_auth=false")
	flag.StringVar(&(user.username), "user", "", "-user=abc")
	flag.StringVar(&(user.passwd), "passwd", "", "-passwd=def")

	flag.Parse()

	switch category {
	case "http", "socks5":
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
			log.Println("local_port is required")
			flag.Usage()

			os.Exit(1)
		}
		newClientConfig := &ClientConfig{
			Proxy:    fmt.Sprintf("%s:%s", proxyHost, proxyPort),
			Local:    fmt.Sprintf("%s:%s", localHost, localPort),
			needAuth: needAuth,
			user:     user,
		}
		log.Printf("local: %s", newClientConfig.Local)
		log.Printf("proxy: %s", newClientConfig.Proxy)
		return category, newClientConfig
	case "server":
		if host == "" {
			log.Println("host is required")
			flag.Usage()

			os.Exit(1)
		}
		if port == "" {
			log.Println("port is required")
			flag.Usage()

			os.Exit(1)
		}
		if caKey == "" {
			log.Println("cert_key is required")
			flag.Usage()

			os.Exit(1)
		}
		if caCrt == "" {
			log.Println("cert_crt is required")
			flag.Usage()

			os.Exit(1)
		}
		newServerConfig := &ServerConfig{
			Server: fmt.Sprintf("%s:%s", host, port),
			CaCrt:  caCrt,
			CaKey:  caKey,
			User:   user,
		}
		return category, newServerConfig

	default:
		log.Println("category is required")
		flag.Usage()
		os.Exit(1)
	}
	return "", nil
}
