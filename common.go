package h2proxy

import (
	"bytes"
	"crypto/tls"
	"flag"
	"fmt"
	"golang.org/x/net/http2"
	"log"
	"net"
	"os"
)

const (
	HTTP    = "http"
	SERVER  = "server"
	SOCKSV5 = "socks5"
)

type ServerConfig struct {
	Server    string
	CaKey     string
	CaCrt     string
	NeedAuth  bool
	User      *UserInfo
	DebugPort int
}

type ClientConfig struct {
	Local     string
	Proxy     string
	needAuth  bool
	user      *UserInfo
	DebugPort int
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

		DebugPort int
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
	flag.IntVar(&DebugPort, "debugPort", 9999, "-debug_port=9999")

	flag.Parse()

	serverRequired := map[string]string{
		"host":     host,
		"port":     port,
		"cert_key": caKey,
		"cert_crt": caCrt,
	}
	clientRequired := map[string]string{
		"proxy_host":     proxyHost,
		"proxy_port":     proxyPort,
		"local_host": localHost,
		"local_port": localPort,
	}
	var errorLog bytes.Buffer

	switch category {
	case HTTP, SOCKSV5:
		for k, v := range clientRequired {
			if v == "" {
				errorLog.WriteString(k)
				errorLog.WriteString(" is required\n")
			}
		}
		newClientConfig := &ClientConfig{
			Proxy:    fmt.Sprintf("%s:%s", proxyHost, proxyPort),
			Local:    fmt.Sprintf("%s:%s", localHost, localPort),
			needAuth: needAuth,
			user:     user,
			DebugPort: DebugPort,

		}
		log.Printf("local: %s", newClientConfig.Local)
		log.Printf("proxy: %s", newClientConfig.Proxy)
		return category, newClientConfig
	case SERVER:
		for k, v := range serverRequired {
			if v == "" {
				errorLog.WriteString(k)
				errorLog.WriteString(" is required\n")
			}
		}
		if errorLog.Len() != 0 {
			log.Print(errorLog.String())
			flag.Usage()
			os.Exit(1)
		}
		newServerConfig := &ServerConfig{
			Server: fmt.Sprintf("%s:%s", host, port),
			CaCrt:  caCrt,
			CaKey:  caKey,
			User:   user,
			DebugPort: DebugPort,
		}
		return category, newServerConfig

	default:
		log.Println("category is required")
		flag.Usage()
		os.Exit(1)
	}
	return "", nil
}
