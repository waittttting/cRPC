module github.com/waittttting/cRPC-client

go 1.14

require (
	github.com/sirupsen/logrus v1.7.0
	github.com/waittttting/cRPC-common v0.0.1 // indirect
)

replace (
	github.com/waittttting/cRPC-common v0.0.1  => ../cRPC-common
	github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.0.0
    google.golang.org/grpc => google.golang.org/grpc v1.26.0
)