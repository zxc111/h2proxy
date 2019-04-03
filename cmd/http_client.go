package main

import (
	"bytes"
	"fmt"
	"github.com/zxc111/h2proxy"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	//local, proxy, needAuth, user = h2proxy.ParseClientConfig()
}

func handler(w http.ResponseWriter, r *http.Request, config *h2proxy.ClientConfig) {
	switch r.Method {
	case http.MethodConnect:
		hijacker, _ := w.(http.Hijacker)
		clientConn, _, err := hijacker.Hijack()
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		defer clientConn.Close()

		remote := "http://" + r.URL.Host
		ConnectMethod(clientConn, remote, config)
	default:
		remote := r.URL.Scheme + "://" + r.URL.Host

		hijacker, _ := w.(http.Hijacker)
		clientConn, _, err := hijacker.Hijack()
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		defer clientConn.Close()

		GetMethod(r, remote, clientConn, config)
	}
}

// connectMethod method (for https, create tunnel)
func ConnectMethod(from net.Conn, remoteAddr string, config *h2proxy.ClientConfig) {

	tr := h2proxy.NewTransport(config.Proxy)

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
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println(resp.StatusCode)
		io.Copy(os.Stdout, resp.Body)
		log.Println("Connect Proxy Server Error")
		return
	}
	fmt.Fprint(from, "HTTP/1.1 200 Connection Established\r\n\r\n")

	go io.Copy(w, from)
	io.Copy(from, resp.Body)
}

// not connectMethod method (http not https,don't need tunnel)
func GetMethod(from *http.Request, remote string, to net.Conn, config *h2proxy.ClientConfig) {

	dump, err := httputil.DumpRequest(from, true)
	if err != nil {
		log.Println(err)
	}
	tr := h2proxy.NewTransport(config.Proxy)

	remoteAddr := remote
	log.Println(remoteAddr)

	req, err := http.NewRequest(
		http.MethodGet,
		remoteAddr,
		bytes.NewBuffer(dump),
	)

	req.Header = from.Header
	if err != nil {
		log.Println(err)
	}
	resp, err := tr.RoundTrip(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println(resp.StatusCode)
		io.Copy(os.Stdout, resp.Body)
		log.Println("Connect Proxy Server Error")
		return
	}
	io.Copy(to, resp.Body)

}

func startHttp(config *h2proxy.ClientConfig) {

	log.Printf("local: %s", config.Local)
	log.Printf("remote: %s", config.Proxy)
	server := &http.Server{
		Addr: config.Local,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodConnect:

				fmt.Println("connectMethod")
				handler(w, r, config)
			default:

				handler(w, r, config)
			}
		}),
	}
	server.ListenAndServe()
}
