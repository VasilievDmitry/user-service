module github.com/lotproject/user-service/pkg

go 1.14

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/golang/protobuf v1.5.2
	github.com/micro/go-micro v1.18.0
	github.com/stretchr/testify v1.4.0
	google.golang.org/protobuf v1.27.1
)

replace github.com/micro/go-micro => github.com/paysuper/go-micro v0.0.0-20220210193104-32a80cb1af1c
