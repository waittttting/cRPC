module cx-mp-user

go 1.16


require (
	github.com/BurntSushi/toml v0.3.1
	github.com/sirupsen/logrus v1.7.0
	github.com/waittttting/cRPC-client v0.0.4
)

replace (
	github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.0.0
	github.com/waittttting/cRPC-client => ../cRPC-client
	github.com/waittttting/cRPC-common => ../cRPC-common
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
)