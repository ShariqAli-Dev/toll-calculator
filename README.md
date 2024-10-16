# toll-calculator

obu generates lat,long
reciever recieves lat,long - puts in a kafka que

### how to start

- `make reciever` (needs to be started first b/c obu connects to its ws)
- `make obu`
- `docker compose -up` to start kafka

## installing protobuf compiler (protoc)

```
sudo apt install protobuf-compiler
```

## installing grpc and protobuffer plugins for golang

protobuff

```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```

grpc

```
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```
