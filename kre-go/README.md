# KRE Golang Runner

This is an implementation in Go for the KRE runner.

## Run Integration Test

Start a NATS and MongoDB servers using:
```
docker-compose up
```

Execute the test running:
```
go test ./...
```

If you want to change the test message, you have to generate the protobuf source code: 
```
protoc -I=test --go_out=. test/test_msg.proto
```
