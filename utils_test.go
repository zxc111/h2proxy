package h2proxy

import (
	"bufio"
	"bytes"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"strings"
	"testing"
)

func TestDebug(t *testing.T) {
	//InitLogger()
	Log = logger()
	res := captureOutput(print, false)
	if strings.Contains(res, "debug") {
		t.Fatal()
	}

	res = captureOutput(print, true)
	if !strings.Contains(res, "debug") {
		t.Fatal()
	}
}

func print() {
	Log.Debugf("test debug", )
	Log.Info("test info", )
}

func captureStdout(f func(), ) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}
func captureOutput(funcToRun func(), debug bool) string {
	var buffer bytes.Buffer

	oldLogger := Log

	encoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	writer := bufio.NewWriter(&buffer)

	var level zapcore.Level
	if debug {
		level = zapcore.DebugLevel
	} else {
		level = zapcore.InfoLevel
	}
	Log = zap.New(
		zapcore.NewCore(encoder, zapcore.AddSync(writer), level)).
		Sugar()

	funcToRun()
	writer.Flush()

	Log = oldLogger

	return buffer.String()
}
