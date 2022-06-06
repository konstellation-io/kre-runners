build:
	./scripts/build_go_node.sh
	docker-compose --env-file .env.e2e -f test_assets/docker-compose.yml up --build
	
# done automatically by docker-compose
init-influx:
	curl "http://localhost:8086/query" --data-urlencode "q=CREATE DATABASE kre"

build_go:
	./scripts/build_go_node.sh

test:
	./scripts/run_go_test.sh

proto: 
	./scripts/generate_proto.sh