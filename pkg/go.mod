module github.com/lotproject/user-service/pkg

go 1.13

require (
	github.com/golang/protobuf v1.5.2
	github.com/micro/go-micro v1.18.0
	google.golang.org/protobuf v1.27.1
)

replace (
	github.com/coreos/etcd => github.com/ozonru/etcd v3.3.20-grpc1.27-origmodule+incompatible
	github.com/micro/go-micro => github.com/paysuper/go-micro v0.0.0-20220210193104-32a80cb1af1c
	google.golang.org/grpc => google.golang.org/grpc v1.27.0
)
