# grpc-server-server
## Install and run
 ```shell
git clone https://github.com/zdarovich/grpc-server-client
cd grpc-server-client
docker build --tag grpc-server:1.0 .
docker run --publish 7000:7000 --name gs grpc-server:1.0
```

## Remove
 ```shell
docker rm --force gs
```

# grpc-server-client
## Run
 ```shell
cd grpc-server-client
go run cmd/client/main.go --url=www.delfi.ee
```
