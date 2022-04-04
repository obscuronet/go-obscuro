The `.go` files in this folder are auto-generated from the `.proto` service definition using the `protoc` Protocol 
Buffer compiler.

Install Protobuf and Protoc-gen-go with:

    brew install protobuf

    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

The files were generated using the following command:

    cd path/to/generated

    protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative enclave.proto
