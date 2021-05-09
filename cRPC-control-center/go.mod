module github.com/waittttting/cRPC-control-center

go 1.14

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/sirupsen/logrus v1.7.0
	github.com/waittttting/cRPC-common v0.0.2
    github.com/gin-gonic/gin v1.6.3
)

replace (
	github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.0.0
	github.com/waittttting/cRPC-common => /Users/changjinsheng/go/go_module/cRPC-common
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
)
