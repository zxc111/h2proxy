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
		Log.Error(err)
	}
}

func logger() *zap.SugaredLogger {
	var tmp *zap.Logger
	lock := new(sync.Mutex)
	lock.Lock()
	defer lock.Unlock()
	var err error
	if Debug {
		tmp, err = zap.NewDevelopment()
	} else {
		tmp, err = zap.NewProduction()
	}
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}

	sugar := tmp.Sugar()
	return sugar
}

func InitLogger() {
	Log = logger()
}
