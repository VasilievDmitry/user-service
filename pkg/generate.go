package pkg

//go:generate protoc -I=. --micro_out=. --micro_opt=paths=source_relative --go_out=. --go_opt=paths=source_relative user_service.proto user_service_entity.proto
//go:generate protoc-go-inject-tag -input=./user_service.pb.go -XXX_skip=bson,json,structure,validate
//go:generate protoc-go-inject-tag -input=./user_service_entity.pb.go -XXX_skip=bson,json,structure,validate
