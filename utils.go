package h2proxy

import (
	"go.uber.org/zap"
	"io"
	"log"
	"sync"
)

var (
	Log  = logger()
	once = new(sync.Once)
)


func closeConn(conn io.Closer) {
	err := conn.Close()
	if err != nil {
		log.Print(err)
	}
}

func logger() *zap.Logger {
	var tmp *zap.Logger
	once.Do(func(){
		var err error
		tmp, err = zap.NewProduction()
		if err != nil {
			log.Fatalf("can't initialize zap logger: %v", err)
		}
	})
	return tmp
}
