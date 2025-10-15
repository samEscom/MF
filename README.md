# MF
A modular financial backend built with Go, gRPC, and Docker. Includes loan and user services, database layer with GORM, and automated testing.


In dev process

go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

brew install protobuf 

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/user.proto