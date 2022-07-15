# build go binary then build all docker images
build:
	./scripts/build_go_node.sh
	docker-compose --env-file .env.e2e -f test_assets/docker-compose.yml up --build
	
# done automatically by docker-compose
init-influx:
	curl "http://localhost:8086/query" --data-urlencode "q=CREATE DATABASE kre"

# build go binary
build_go:
	./scripts/build_go_node.sh

# run tests at test_go main package
test:
	./scripts/run_go_test.sh

# build all proto files related to e2e
proto: 
	./scripts/generate_proto.sh
