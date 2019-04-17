package h2proxy

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
)

const (
	HTTP    = "http"
	SERVER  = "server"
	SOCKSV5 = "socks5"

	LOCAL_HOST = "local_host"
	LOCAL_PORT = "local_port"

	PROXY_HOST = "proxy_host"
	PROXY_PORT = "proxy_port"
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
	Category  string
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
	flag.StringVar(&localHost, LOCAL_HOST, "localhost", "-local_host=127.0.0.1")
	flag.StringVar(&localPort, LOCAL_PORT, "3002", "-local_port=4000")
	flag.StringVar(&proxyHost, PROXY_HOST, "", "-porxy_host=xxx.xxx.xxx.xxx")
	flag.StringVar(&proxyPort, PROXY_PORT, "", "-proxy_port=3000")

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
		PROXY_HOST: proxyHost,
		PROXY_PORT: proxyPort,
		LOCAL_HOST: localHost,
		LOCAL_PORT: localPort,
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
			Proxy:     fmt.Sprintf("%s:%s", proxyHost, proxyPort),
			Local:     fmt.Sprintf("%s:%s", localHost, localPort),
			needAuth:  needAuth,
			user:      user,
			DebugPort: DebugPort,
			Category:  category,
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
			Server:    fmt.Sprintf("%s:%s", host, port),
			CaCrt:     caCrt,
			CaKey:     caKey,
			User:      user,
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
