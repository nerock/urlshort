# _Nerock URL Shortener_

Simple URL shortener that allows to create, retrieve and delete shortened URLs

## Demo
Try the [demo](https://nerock.dev/api/docs)

## How to run
```
go build -o urlshort cmd/main.go
./urlshort
```

## Rebuild gRPC definitions
### Requirements
- [Protoc compiler v3](https://grpc.io/docs/protoc-installation/)
- [protoc-gen-go](google.golang.org/protobuf/cmd/protoc-gen-go@v1.26)
- [protoc-gen-go-grpc](oogle.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1)
### How to
`./generate_grpc.sh`

## Documentation
The API documentation is available at `/docs` endpoint and can the file can be edited in `docs/swagger.json`

## Environment variables
|ENV VAR|SUMMARY|DEFAULT|
|-------|-------|-------|
|PORT|HTTP Server port|8080|
|GRPC_PORT|gRPC Server port|50051|
|DBCONN|Sqlite DB connection string|urlshort.db|
|DOMAIN|Domain where the app is deployed to build short URLs|localhost|