### 1️⃣ Генерация протофайлов
```sh
protoc -I proto proto/sso/sso.proto --go_out=./gen/go --go_opt=paths=source_relative --go-grpc_out=./gen/go/ --go-grpc_opt=paths=source_relative
```

### Запуск программы

```sh
go run cmd/sso/main.go --config=./config/local.yaml
```