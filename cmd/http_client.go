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

var (
	proxy, local string
	user *h2proxy.UserInfo
	needAuth bool
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	local, proxy, needAuth, user = h2proxy.ParseConfig()
}

func handler(w http.ResponseWriter, r *http.Request) {
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
		ConnectMethod(clientConn, remote)
	default:
		remote := r.URL.Scheme + "://" + r.URL.Host

		hijacker, _ := w.(http.Hijacker)
		clientConn, _, err := hijacker.Hijack()
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		defer clientConn.Close()

		GetMethod(r, remote, clientConn)
	}
}

// connect method (for https, create tunnel)
func ConnectMethod(from net.Conn, remoteAddr string) {

	tr := h2proxy.NewTransport(proxy)

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

// not connect method (http not https,don't need tunnel)
func GetMethod(from *http.Request, remote string, to net.Conn) {

	dump, err := httputil.DumpRequest(from, true)
	if err != nil {
		log.Println(err)
	}
	tr := h2proxy.NewTransport(proxy)

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

func main() {

	log.Printf("local: %s", local)
	log.Printf("remote: %s", proxy)
	server := &http.Server{
		Addr: local,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodConnect:

				fmt.Println("connect")
				handler(w, r)
			default:

				handler(w, r)
			}
		}),
	}
	server.ListenAndServe()
}
