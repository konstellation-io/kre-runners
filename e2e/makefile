build:
	docker-compose --env-file .env.e2e -f test_assets/docker-compose.yml up --build
	
init-influx:
	curl "http://localhost:8086/query" --data-urlencode "q=CREATE DATABASE kre"

build2:
	docker-compose --env-file .env.e2e -f test_assets/docker-compose.yml up --build \
  curl "http://localhost:8086/query" --data-urlencode "q=CREATE DATABASE kre"

#command: sh -c "apk add curl && sleep 5 && curl 'http://localhost:8086/query' --data-urlencode 'q=CREATE DATABASE kre'"
