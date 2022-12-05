package proto

//go:generate protoc --proto_path=./v1 --micro_out=./v1 --micro_opt=paths=source_relative --go_out=./v1 --go_opt=paths=source_relative ./v1/user_service.proto ./v1/user_service_entity.proto
//go:generate protoc-go-inject-tag -input=./v1/*.pb.go -XXX_skip=bson,json,structure,validate

//go:generate protoc --proto_path=./game-service --micro_out=./game-service --micro_opt=paths=source_relative --go_out=./game-service --go_opt=paths=source_relative ./game-service/service.proto
