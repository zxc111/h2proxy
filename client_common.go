package h2proxy

import (
	"context"
	"crypto/tls"
	"fmt"
	"golang.org/x/net/http2"
	"io"
	"net"
	"net/http"
	"sync"
	"time"
)

// create http request with connect method
func CreateTunnel(ctx context.Context, from net.Conn, remoteAddr string, config *ClientConfig) {

	defer cost(time.Now().UnixNano(), remoteAddr)

	tr := NewTransportWithProxy(config.Proxy)
	defer tr.CloseIdleConnections()
	r, w := io.Pipe()
	defer r.Close()
	defer w.Close()

	Log.Info(remoteAddr)
	timeoutCtx, cancel := context.WithCancel(ctx)
	req, err := http.NewRequestWithContext(
		timeoutCtx,
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

	// if category = http, return 200 for connect method Established
	if config.Category == HTTP {
		_, err := fmt.Fprint(from, "HTTP/1.1 200 Connection Established\r\n\r\n")
		Log.Debug(err)
	}
	exit1 := make(chan struct{})

	go func() {
		var wg sync.WaitGroup
		wg.Add(2)
		// read from remote
		timer := time.AfterFunc(60*time.Second, func() {
			cancel()
		})

		go func(dst net.Conn, src io.ReadCloser, group *sync.WaitGroup, nn int) {
			defer wg.Done()
			for {
				timer.Reset(50 * time.Second)
				res := make([]byte, 10240)
				n, err := src.Read(res)
				if err != nil {
					if io.EOF != err {
						if e, ok := err.(net.Error); ok && e.Timeout() {
						} else {
							Log.Info(err)
						}
					}
					return
				}

				dst.SetWriteDeadline(time.Now().Add(50 * time.Second))
				if n != 0 {
					_, err := dst.Write(res[:n])
					if err != nil {
						Log.Error(err)
						return
					}
				}
			}
		}(from, resp.Body, &wg, 1)

		go copyData(w, from, &wg, 2, cancel)
		wg.Wait()
		Log.Info("copyData finish")
		close(exit1)
	}()
	select {
	case <-exit1:
	case <-ctx.Done():
	case <-time.Tick(time.Hour):
	}
}

func copyData(dst *io.PipeWriter, src net.Conn, wg *sync.WaitGroup, nn int, ctxCancel context.CancelFunc) {
	defer wg.Done()
	res := make([]byte, 65535)

	defer ctxCancel()
	for {
		src.SetReadDeadline(time.Now().Add(50 * time.Second))
		n, err := src.Read(res)
		if err != nil {
			if io.EOF != err {
				if e, ok := err.(net.Error); ok && e.Timeout() {
				} else {
					Log.Info(err)
				}
			}
			return
		}
		if n != 0 {
			_, err := dst.Write(res[:n])
			if err != nil {
				Log.Error(err)
				return
			}
		}
	}
}

// make new transport
func NewTransportWithProxy(proxyAddr string) *http2.Transport {
	return &http2.Transport{
		DialTLS: func(network, addr string, config *tls.Config) (net.Conn, error) {
			return tls.DialWithDialer(&net.Dialer{Timeout: time.Minute}, "tcp", proxyAddr, &tls.Config{
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
