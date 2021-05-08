module github.com/waittttting/cRPC-common

go 1.14

require (
	github.com/bwmarrin/snowflake v0.3.0 // indirect
	github.com/coreos/etcd v3.3.22+incompatible // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-00010101000000-000000000000 // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/uuid v1.2.0 // indirect
	github.com/sirupsen/logrus v1.7.0
	go.etcd.io/etcd v3.3.22+incompatible
	go.uber.org/zap v1.16.0 // indirect
	google.golang.org/grpc v1.35.0 // indirect
)

replace (
	github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.0.0
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
)
