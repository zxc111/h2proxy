package h2proxy

import (
	"bytes"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
)

type Http2Server struct {
	Config *ServerConfig
}

func (h Http2Server) Start() {
	config := h.Config
	server := &http.Server{
		Addr:    config.Server,
		Handler: http.HandlerFunc(handle(config)),
	}

	// require cert.
	// generate cert for test:
	// openssl req -new -x509 -days 365 -key test1.key -out test1.crt
	if err := server.ListenAndServeTLS(config.CaCrt, config.CaKey); err != nil {
		log.Fatal(err)
	}
}

func handle(config *ServerConfig) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if config.NeedAuth && !CheckAuth(config.User, r) {
			// TODO check auth
			w.WriteHeader(400)
			return
		}
		switch r.Method {
		case http.MethodConnect:
			connectMethod(w, r)
		default:
			get(w, r)
		}
	}
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

func connectMethod(w http.ResponseWriter, r *http.Request) {
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
	defer r.Body.Close()

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

	to := flushWriter{w}

	req, _ := http.NewRequest(
		r.Method,
		"http://"+r.Host+r.RequestURI,
		r.Body,
	)
	cli := http.Client{}
	req.Header = r.Header

	resp, err := cli.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	res, _ := httputil.DumpResponse(resp, true)
	io.Copy(to, bytes.NewBuffer(res))
}
