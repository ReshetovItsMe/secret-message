# Secret Assistant

To generate proto run in `proto` folder

`protoc --go_out=./generated/go --go_opt=paths=source_relative \
    --go-grpc_out=./generated/go --go-grpc_opt=paths=source_relative \
    message.proto`


protoc --go_out=./proto/generated/go --go_opt=Mproto/message.proto=github.com/ReshetovItsMe/one-time-messaging-exchange-be/message \
    --go-grpc_out=./proto/generated/go --go-grpc_opt=Mproto/message.proto=github.com/ReshetovItsMe/one-time-messaging-exchange-be/message \
    proto/message.proto