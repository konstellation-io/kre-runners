.PHONY: build build_go init-influx test proto

# build go binary then build all docker images
build: build_go
	docker-compose -f test_assets/docker-compose.yml up --build
	
# build go binary
build_go:
	./scripts/build_go_node.sh

# run tests at test_go main package
test:
	./scripts/run_go_test.sh

# build all proto files related to e2e
proto: 
	./scripts/generate_proto.sh
