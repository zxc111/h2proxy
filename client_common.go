package h2proxy

import (
	"crypto/tls"
	"fmt"
	"golang.org/x/net/http2"
	"io"
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

	Log.Info(remoteAddr)

	req, err := http.NewRequest(
		http.MethodConnect,
		remoteAddr,
		r,
	)
	if err != nil {
		Log.Error(err)
	}
	if config.NeedAuth {
		SetAuthInHeader(config.User, req)
	}

	resp, err := tr.RoundTrip(req)
	if err != nil {
		Log.Error(err)
		return
	}

	defer closeConn(resp.Body)

	if resp.StatusCode != 200 {
		Log.Info(resp.StatusCode)
		// TODO
		io.Copy(os.Stdout, resp.Body)
		Log.Info("Connect Proxy Server Error")
		return
	}

	// if category = http, return 200 for connect method Established
	if config.Category == HTTP {
		_, err := fmt.Fprint(from, "HTTP/1.1 200 Connection Established\r\n\r\n")
		Log.Debug(err)
	}

	if Debug {
		var wg sync.WaitGroup
		wg.Add(1)
		wg.Add(1)
		go copyData(from, resp.Body, &wg, 1)
		go copyData(w, from, &wg, 2)
		wg.Wait()
		Log.Info("copyData finish")
	} else {
		go io.Copy(w, from)
		io.Copy(from, resp.Body)
	}

}

// for debug
func copyData(dst io.Writer, src io.Reader, wg *sync.WaitGroup, num int) {
	defer wg.Done()
	res := make([]byte, 65535)

	for {
		n, err := src.Read(res)
		Log.Debugf("read %d", num)
		if err != nil {
			Log.Info("eof")
			break
		}
		if n != 0 {
			//Log.Info(res[:n])
			//Log.Info(string(res[:n]))
			_, err := dst.Write(res[:n])
			Log.Error(err)
		}
	}
}

// make new transport
func NewTransport(proxyAddr string) *http2.Transport {
	return &http2.Transport{
		DialTLS: func(network, addr string, config *tls.Config) (net.Conn, error) {
			return tls.DialWithDialer(&net.Dialer{Timeout: 10*time.Minute}, "tcp", proxyAddr, &tls.Config{
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
	Log.Infof("%s cost: %d ms", path, t)
}
