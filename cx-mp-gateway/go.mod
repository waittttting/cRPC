module cx-mp-gateway

go 1.14

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.6.3
	github.com/sirupsen/logrus v1.7.0
	github.com/waittttting/cRPC-client v0.0.4
	github.com/waittttting/cRPC-common v0.0.2
)

replace (
	github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.0.0
	github.com/waittttting/cRPC-client => ../cRPC-client
    github.com/waittttting/cRPC-common => ../cRPC-common
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
)
