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
	"sync"
	"time"
)

// create http request with connect method
func CreateTunnel(from net.Conn, remoteAddr string, config *ClientConfig) {

	defer cost(time.Now().UnixNano(), remoteAddr)

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

	// if category = http, return 200 for connect method Established
	if config.Category == HTTP {
		fmt.Fprint(from, "HTTP/1.1 200 Connection Established\r\n\r\n")
	}

	if Debug {
		var wg sync.WaitGroup
		wg.Add(1)
		wg.Add(1)
		go copy(from, resp.Body, &wg, 1)
		go copy(w, from, &wg, 2)
		wg.Wait()
		log.Println("copy finish")
	} else {
		go io.Copy(w, from)
		io.Copy(from, resp.Body)
	}

}

// for debug
func copy(dst io.Writer, src io.Reader, wg *sync.WaitGroup, num int) {
	defer wg.Done()
	res := make([]byte, 65535)

	for {
		n, err := src.Read(res)
		log.Printf("read %d", num)
		if err != nil {
			log.Println("eof")
			break
		}
		if n != 0 {
			log.Println(res[:n])
			log.Println(string(res[:n]))
			dst.Write(res[:n])
		}
	}
}

// make new transport
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

// print request cost
func cost(start int64, path string) {
	t := (time.Now().UnixNano() - start) / 1000000
	log.Printf("%s cost: %d ms", path, t)
}
