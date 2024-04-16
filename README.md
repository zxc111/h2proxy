# h2proxy

[![Go Report Card](https://goreportcard.com/badge/github.com/zxc111/h2proxy)](https://goreportcard.com/report/github.com/zxc111/h2proxy)
[![Build Status](https://travis-ci.com/zxc111/h2proxy.svg?branch=master)](https://travis-ci.com/zxc111/h2proxy)

http2 proxy server &amp;&amp; client


[Android Client](https://github.com/zxc111/SmartProxy)

## 如何运行
从 [release](https://github.com/zxc111/h2proxy/releases) 下载对应系统编译好的
```bash
./linux -conf conf.toml
./mac -conf conf.toml
或是
./win.exe -conf conf.toml
```
### 服务端
```toml
category = "server"            # 执行模式-服务端
[server]
	Server   = "0.0.0.0:1234"  # 监听地址和端口
	CaKey    = "xxx.key"       # 证书
	CaCrt    = "xxx.cert"      # 证书私钥
	NeedAuth = true            # 是否身份认证 
	Pprof    = 12345           # pprof 端口
	debug = true               # 是否开启 debug

[server.user]
    username = "aaa"        # 用户名
    passwd   = "bbb"        # 密码
```

### 客户端
#### sock5
```toml
category = "socks5"             
[client]
	Local    = '0.0.0.0:30000' 
	Proxy    = 'host:30000'
	needAuth = false
	Pprof    = 12345
```
#### http
```toml
category = "http"             
[client]
	Local    = '0.0.0.0:30000' 
	Proxy    = 'host:30000'
	needAuth = true
	Pprof    = 12345
[server.user]
    username = "aaa"
    passwd   = "bbb"
```