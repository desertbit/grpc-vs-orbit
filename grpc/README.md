## Source

- https://github.com/grpc/grpc-go/blob/master/examples/helloworld
- https://grpc.io/docs/quickstart/go/

## Protobuf

```
go get -u github.com/golang/protobuf/protoc-gen-go
protoc -I api/ api/api.proto --go_out=plugins=grpc:api
```