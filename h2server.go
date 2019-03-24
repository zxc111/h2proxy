package main

import (
	"bytes"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	server := &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodConnect:
				connect(w, r)
			default:
				get(w, r)
			}
		}),
	}
	server.ListenAndServeTLS("localhost.cert", "localhost.key")
}

type flushWriter struct {
	w io.Writer
}

func (fw flushWriter) Write(p []byte) (n int, err error) {
	n, err = fw.w.Write(p)
	if f, ok := fw.w.(http.Flusher); ok {
		f.Flush()
	}
	return
}

func connect(w http.ResponseWriter, r *http.Request) {
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	remoteAddr := r.Host
	if strings.Count(remoteAddr, ":") == 0 {
		remoteAddr += ":443"
	}
	conn, err := net.Dial("tcp", remoteAddr)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	to := flushWriter{w}

	go io.Copy(conn, r.Body)

	io.Copy(to, conn)
}

func get(w http.ResponseWriter, r *http.Request) {
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Println(err)
	}
	if strings.Count(r.Host, ":") == 0 {
		r.Host += ":80"
	}
	conn, err := net.Dial("tcp", r.Host)
	if err != nil {
		log.Println(err)
		return
	}
	to := flushWriter{w}
	dump = bytes.Replace(dump, []byte("HTTP/2.0"), []byte("HTTP/1.1"), -1)
	log.Println(string(dump))

	defer conn.Close()
	go io.Copy(conn, bytes.NewBuffer(dump))
	io.Copy(to, conn)
}
