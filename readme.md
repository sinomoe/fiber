# Fiber

## Usage

```bash
docker pull redis
docker run --name redis-server -d -p 6379:6379 redis

go env -w GOPROXY=https://goproxy.cn,direct
go env -w GO111MODULE=on
go mod tidy

nohup go run cmd/comet/main.go > /dev/null &
nohup go run cmd/job/main.go > /dev/null &
nohup go run cmd/logic/main.go > /dev/null &
```

## Components

### Comet

* rpc port: 8869
* websocket port: 8879

### Logic

* api port: 8859

### Job