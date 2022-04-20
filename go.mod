module github.com/lotproject/user-service

go 1.13

require (
	github.com/InVisionApp/go-health/v2 v2.1.2
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-sql-driver/mysql v1.6.0
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.3.0
	github.com/jmoiron/sqlx v1.3.4
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/lotproject/game-service/pkg v0.0.0-20220419045839-a78768287ad4
	github.com/lotproject/go-helpers/db v0.0.0-20220223094055-21a0af1e4859
	github.com/lotproject/go-helpers/hash v0.0.0-20220223094055-21a0af1e4859
	github.com/lotproject/go-helpers/log v0.0.0-20220223094055-21a0af1e4859
	github.com/lotproject/go-helpers/random v0.0.0-20220223094055-21a0af1e4859
	github.com/lotproject/user-service/pkg v0.0.0-20220420184741-98c89215a2a6
	github.com/micro/go-micro v1.18.0
	github.com/natefinch/lumberjack v2.0.0+incompatible
	github.com/prometheus/client_golang v1.7.1
	github.com/stretchr/testify v1.7.0
	go.uber.org/zap v1.21.0
	golang.org/x/crypto v0.0.0-20200221231518-2aa609cf4a9d
	google.golang.org/protobuf v1.27.1
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
)

replace (
	github.com/coreos/etcd => github.com/ozonru/etcd v3.3.20-grpc1.27-origmodule+incompatible
	github.com/micro/go-micro => github.com/paysuper/go-micro v0.0.0-20220210193104-32a80cb1af1c
	google.golang.org/grpc => google.golang.org/grpc v1.27.0
)
