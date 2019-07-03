package h2proxy

import (
	"go.uber.org/zap"
	"io"
	"log"
	"sync"
)

var (
	Log  *zap.SugaredLogger
	once = new(sync.Once)
)

func closeConn(conn io.Closer) {
	err := conn.Close()
	if err != nil {
		log.Print(err)
	}
}

func logger() *zap.SugaredLogger {
	var tmp *zap.Logger
	once.Do(func() {
		var err error
		if Debug {
			tmp, err = zap.NewDevelopment()
		} else {
			tmp, err = zap.NewProduction()
		}
		if err != nil {
			log.Fatalf("can't initialize zap logger: %v", err)
		}

	})
	sugar := tmp.Sugar()
	return sugar
}

func InitLogger() {
	Log = logger()
}