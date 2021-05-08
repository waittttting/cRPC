module github.com/waittttting/cRPC-config-center

go 1.14

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/gin-gonic/gin v1.6.3
	github.com/mattn/go-sqlite3 v1.14.6 // indirect
	github.com/sirupsen/logrus v1.7.0
	github.com/waittttting/cRPC-common v0.0.2
	gorm.io/driver/mysql v1.0.3
	gorm.io/driver/sqlite v1.1.4 // indirect
	gorm.io/gorm v1.20.9
	github.com/go-redis/redis v6.15.9+incompatible
)

replace (
	github.com/waittttting/cRPC-common => ../cRPC-common
	github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.0.0
    google.golang.org/grpc => google.golang.org/grpc v1.26.0
)
