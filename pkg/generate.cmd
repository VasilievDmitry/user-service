echo "GENERATE PROTO"
protoc --micro_out=./userpb --micro_opt=paths=source_relative --go_out=./userpb --go_opt=paths=source_relative user_service.proto

echo "INJECTING TAGS"
protoc-go-inject-tag -input=./userpb/user_service.pb.go -XXX_skip=bson,json,structure,validate

echo "GENERATING MOCKS"
mockery --all --dir=./userpb --output=./userpb