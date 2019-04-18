package h2proxy


import (
	"io"
	"log"
)

func closeConn(conn io.Closer) {
	err := conn.Close()
	if err != nil {
		log.Println(err)
	}
}