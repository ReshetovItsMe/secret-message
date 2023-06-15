# Secret Assistant

To generate proto run

export GO_PATH=~/go
export PATH=$PATH:/$GO_PATH/bin

protoc --proto_path=../proto \
 --go_out=./proto --go_opt=paths=source_relative \
 --go-grpc_out=./proto --go-grpc_opt=paths=source_relative \
 message.proto
