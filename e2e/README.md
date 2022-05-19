

```sh
docker-compose --env-file .env.e2e -f test_assets/docker-compose.yml up --build

python -m grpc_tools.protoc \
       -I=test_assets/entrypoint-vol/krt-files/ \
       --python_out=. \
       --grpc_python_out=. \
       test_assets/entrypoint-vol/krt-files/public_input.proto

protoc -I=test_assets/entrypoint-vol/krt-files/ --go_out=test_assets/nodeC-vol/src/node test_assets/entrypoint-vol/krt-files/public_input.proto

protoc -I=./test_assets/entrypoint-vol/krt-files/ \
  --go_out=./test_assets/nodeC-vol/src/node \
  --go_opt=paths=source_relative ./test_assets/entrypoint-vol/krt-files/*.proto

```