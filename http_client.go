package h2proxy

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
)

type HttpProxy struct {
	Config *ClientConfig
}

func handler(w http.ResponseWriter, r *http.Request, config *ClientConfig) {
	switch r.Method {
	case http.MethodConnect:
		hijacker, _ := w.(http.Hijacker)
		clientConn, _, err := hijacker.Hijack()
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		defer closeConn(clientConn)

		remote := "http://" + r.URL.Host
		ctx, _ := context.WithTimeout(context.Background(), time.Hour)

		CreateTunnel(ctx, clientConn, remote, config)
	default:
		remote := r.URL.Scheme + "://" + r.URL.Host

		hijacker, _ := w.(http.Hijacker)
		clientConn, _, err := hijacker.Hijack()
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		defer closeConn(clientConn)
		ctx, _ := context.WithTimeout(context.Background(), time.Hour)

		GetMethod(ctx, r, remote, clientConn, config)
	}
}

// not connectMethod method (http not https,don't need tunnel)
func GetMethod(ctx context.Context, from *http.Request, remote string, to net.Conn, config *ClientConfig) {

	dump, err := httputil.DumpRequest(from, true)
	if err != nil {
		Log.Error(err)
	}
	tr := NewTransport(config.Proxy)

	remoteAddr := remote
	Log.Info(remoteAddr)

	req, err := http.NewRequest(
		http.MethodGet,
		remoteAddr,
		bytes.NewBuffer(dump),
	)
	if err != nil {
		Log.Error(err)
		return
	}

	req.Header = from.Header

	resp, err := tr.RoundTrip(req)
	if err != nil {
		Log.Error(err)
		return
	}
	defer closeConn(resp.Body)

	if resp.StatusCode != 200 {
		Log.Info(resp.StatusCode)
		io.Copy(os.Stdout, resp.Body)
		fmt.Fprint(to, resp.StatusCode)
		Log.Info("Connect Proxy Server Error")
		return
	}
	exit := make(chan struct{})
	go func() {
		io.Copy(to, resp.Body)
		close(exit)
	}()
	select {
	case <-exit:
	case <-ctx.Done():
	}
}

func (h HttpProxy) Start() {
	config := h.Config
	Log.Info("local: %s", config.Local)
	Log.Info("remote: %s", config.Proxy)
	server := &http.Server{
		Addr: config.Local,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodConnect:
				handler(w, r, config)
			default:
				handler(w, r, config)
			}
		}),
	}
	Log.Fatal(server.ListenAndServe())
}
