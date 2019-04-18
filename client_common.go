package h2proxy

import (
	"crypto/tls"
	"fmt"
	"golang.org/x/net/http2"
	"io"
	"log"
	"net"
	"net/http"
	"os"
)

// create http request with connect method
func CreateTunnel(from net.Conn, remoteAddr string, config *ClientConfig) {

	tr := NewTransport(config.Proxy)

	r, w := io.Pipe()

	log.Println(remoteAddr)

	req, err := http.NewRequest(
		http.MethodConnect,
		remoteAddr,
		r,
	)
	if err != nil {
		log.Println(err)
	}

	resp, err := tr.RoundTrip(req)
	if err != nil {
		log.Println(err)
		return
	}

	defer closeConn(resp.Body)

	if resp.StatusCode != 200 {
		log.Println(resp.StatusCode)
		// TODO
		io.Copy(os.Stdout, resp.Body)
		log.Println("Connect Proxy Server Error")
		return
	}
	if config.Category == HTTP {
		fmt.Fprint(from, "HTTP/1.1 200 Connection Established\r\n\r\n")
	}

	go io.Copy(w, from)
	io.Copy(from, resp.Body)
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
