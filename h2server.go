package h2proxy

import (
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

type Http2Server struct {
	Config *ServerConfig
}

var noAuthBody = []byte("Proxy Authentication Required")

func (h Http2Server) Start() {
	config := h.Config
	//http.HandleFunc("/test", handle(config))

	server := &http.Server{
		Addr:    config.Server,
		Handler: http.HandlerFunc(handle(config)),
	}
	// require cert.
	// generate cert for test:
	// openssl req -new -x509 -days 365 -key test1.key -out test1.crt
	Log.Fatal(server.ListenAndServeTLS(config.CaCrt, config.CaKey))
}

func handle(config *ServerConfig) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		if config.NeedAuth && !CheckAuth(config.User, r) {
			Log.Debug("auth failed")

			w.Header().Set("Proxy-Authenticate", `Basic realm="Access to internal site"`)
			w.WriteHeader(407)

			_, err := w.Write(noAuthBody)
			if err != nil {
				Log.Error(err)
			}
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
	defer cost(time.Now().UnixNano(), r.URL.RequestURI())

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
		Log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer closeConn(conn)

	to := flushWriter{w}
	defer closeConn(r.Body)

	go io.Copy(conn, r.Body)

	io.Copy(to, conn)
}

func get(w http.ResponseWriter, r *http.Request) {
	defer cost(time.Now().UnixNano(), r.URL.RequestURI())
	defer closeConn(r.Body)

	f, ok := w.(http.Flusher)
	if !ok {
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
	req.Header.Del("Proxy-Authorization")

	resp, err := cli.Do(req)
	if err != nil {
		Log.Fatal(err)
	}
	for k, v := range resp.Header {
		if len(v) == 0 {
			continue
		}
		w.Header().Set(k, v[0])
	}
	f.Flush()

	defer closeConn(resp.Body)
	io.Copy(to, resp.Body)
}
