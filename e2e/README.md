

```sh
docker-compose --env-file .env -f test_assets/docker-compose.yml up --build

python -m grpc_tools.protoc \
       -I=test_assets/entrypoint-vol/krt-files/ \
       --python_out=. \
       --grpc_python_out=. \
       test_assets/entrypoint-vol/krt-files/public_input.proto
```
```sh