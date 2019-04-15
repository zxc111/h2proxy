module github.com/zxc111/h2proxy

go 1.12

require golang.org/x/net v0.0.0-20190320064053-1272bf9dcd53

require golang.org/x/text v0.3.0

replace (
	golang.org/x/net => github.com/golang/net v0.0.0-20190320064053-1272bf9dcd53
	golang.org/x/text => github.com/golang/text v0.3.0
)
