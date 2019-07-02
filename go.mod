module github.com/zxc111/h2proxy

go 1.12

require golang.org/x/net v0.0.0-20190320064053-1272bf9dcd53

require (
	github.com/google/pprof v0.0.0-20190515194954-54271f7e092f // indirect
	github.com/ianlancetaylor/demangle v0.0.0-20181102032728-5e5cf60278f6 // indirect
	github.com/pkg/errors v0.8.1 // indirect
	github.com/stretchr/testify v1.3.0 // indirect
	go.uber.org/atomic v1.4.0 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/zap v1.10.0
	golang.org/x/arch v0.0.0-20190312162104-788fe5ffcd8c // indirect
	golang.org/x/tools v0.0.0-20190606050223-4d9ae51c2468 // indirect
)

replace (
	golang.org/x/net => github.com/golang/net v0.0.0-20190320064053-1272bf9dcd53
	golang.org/x/text => github.com/golang/text v0.3.0
)
