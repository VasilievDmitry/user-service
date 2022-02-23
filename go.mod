module github.com/lotproject/user-service

go 1.13

require (
	github.com/InVisionApp/go-health v2.1.0+incompatible
	github.com/InVisionApp/go-logger v1.0.1 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-sql-driver/mysql v1.6.0
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.3.0
	github.com/jmoiron/sqlx v1.3.4
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/lotproject/go-helpers v0.0.0-20220222054749-0ae55f93fcf4
	github.com/lotproject/go-helpers/db v0.0.0-20220222054749-0ae55f93fcf4
	github.com/lotproject/go-helpers/log v0.0.0-20220222054749-0ae55f93fcf4
	github.com/lotproject/go-proto/go/user_service v0.0.0-20220223073357-cb57dcd01cf8
	github.com/micro/go-micro v1.18.0
	github.com/natefinch/lumberjack v2.0.0+incompatible
	github.com/prometheus/client_golang v1.7.1
	go.uber.org/zap v1.21.0
	golang.org/x/crypto v0.0.0-20200221231518-2aa609cf4a9d
	gopkg.in/DATA-DOG/go-sqlmock.v1 v1.3.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
)

replace (
	github.com/coreos/etcd => github.com/ozonru/etcd v3.3.20-grpc1.27-origmodule+incompatible
	github.com/micro/go-micro => github.com/paysuper/go-micro v0.0.0-20220210193104-32a80cb1af1c
	google.golang.org/grpc => google.golang.org/grpc v1.27.0
)
