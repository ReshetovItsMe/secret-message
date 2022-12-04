# Secret Assistant

To generate proto run

protoc --proto_path=../proto \
  --go_out=../proto/generated/go --go_opt=paths=source_relative \
    --go-grpc_out=../proto/generated/go --go-grpc_opt=paths=source_relative \
    message.proto