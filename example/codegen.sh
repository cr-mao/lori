find proto -type f -name "*.pb.go" -delete

# https://developers.google.com/protocol-buffers/docs/reference/go-generated#package
# 默认就是import 模式
protoc --go_out=.  --go_opt=paths=import  \
    --go-grpc_out=.  --go-grpc_opt=paths=import \
    proto/helloworld.proto

