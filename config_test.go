package h2proxy

import (
	"github.com/BurntSushi/toml"
	"log"
	"testing"
)

func TestParseToml(t *testing.T) {
	tomlData := `
[server]
	Server   ='a'
	CaKey    ='b'
	CaCrt    ='c'
	NeedAuth = true
	Pprof   = 1234
[server.User]
Username = "123"
Passwd = "abc"

`
	conf := &FileConfig{}
	_, err := toml.Decode(tomlData, &conf)
	if err != nil {
		log.Fatal(err)
	}
	if conf == nil || conf.Server == nil || conf.Server.User == nil {
		t.Fatal()
	}

	if conf.Server.Server != "a" || conf.Server.CaKey != "b" || conf.Server.CaCrt != "c" ||
		conf.Server.NeedAuth != true || conf.Server.Pprof != 1234 || conf.Server.User.Username != "123" ||
		conf.Server.User.Passwd != "abc" {
		t.Fatal()
	}
}

func TestParseConfigFile(t *testing.T) {
	conf := parseFile("config.toml.example")

	if conf == nil || conf.Server == nil || conf.Server.User == nil {
		t.Fatal()
	}
	if conf.Category != "server" {
		t.Fatal()
	}
	if conf.Server.Server != "a" || conf.Server.CaKey != "b" || conf.Server.CaCrt != "c" ||
		conf.Server.NeedAuth != true || conf.Server.Pprof != 1234 || conf.Server.User.Username != "123" ||
		conf.Server.User.Passwd != "abc" {
		t.Fatal()
	}

	if conf.Client.Local != "a" || conf.Client.Proxy != "c" ||
		conf.Client.NeedAuth != true || conf.Client.Pprof != 12345 || conf.Client.User.Username != "abc" ||
		conf.Client.User.Passwd != "321" {
		t.Fatal()
	}
}
