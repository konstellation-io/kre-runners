# KRE Golang Runner

This is an implementation in Go for the KRE runner.

## Run Integration Test

Start a NATS and MongoDB servers using:
```bash
# On kre-go/test
docker-compose up
```

Execute the test running:
```bash
# On kre-go
go test ./...
```

If you want to change the test message, you have to generate the protobuf source code: 
```bash
# On kre-go
protoc -I=test --go_out=. --go_opt=paths=source_relative \test/test_msg.proto
```
