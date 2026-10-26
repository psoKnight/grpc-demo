# grpc-demo

一个参考 `server-utils` 架构裁剪出的最小 Go gRPC 项目。

## 架构

```text
.
├── pbdata/greeterpb/             # gRPC 协议、生成代码和 pb.sh
├── greeter/                      # 具体 gRPC 服务，对应参考项目中的 message/
│   ├── config/
│   ├── service/
│   ├── main.go
│   ├── start.sh
│   ├── deploy.sh
│   └── docker/Dockerfile
├── generalgateway/               # HTTP 入口与 gRPC 客户端
│   ├── config/
│   ├── controller/
│   ├── grpc_client/
│   ├── service/
│   ├── main.go
│   ├── start.sh
│   ├── deploy.sh
│   └── docker/Dockerfile
├── common/healthing/             # 通用健康检查
└── go.mod
```

调用链与参考项目一致：

```text
HTTP 请求
  -> generalgateway/controller
  -> generalgateway/service
  -> generalgateway/grpc_client
  -> pbdata/greeterpb
  -> greeter/service
```

## 微服务文档

- [greeter 构建与部署](greeter/README.md)
- [generalgateway 构建与部署](generalgateway/README.md)

## 本地运行

先启动 gRPC 服务：

```bash
go run ./greeter
```

再启动 HTTP 网关：

```bash
go run ./generalgateway
```

调用：

```bash
curl -X POST http://localhost:8080/v1/hello \
  -H "Content-Type: application/json" \
  -d '{"name":"Codex"}'
```

## 服务打包

每个微服务拥有独立的 `start.sh`、`deploy.sh` 和 `docker/Dockerfile`，与参考项目保持一致。

只生成 Linux 发布目录，不构建镜像：

```bash
cd greeter
./deploy.sh -b test

cd ../generalgateway
./deploy.sh -b test
```

发布目录结构：

```text
docker/release/
├── greeter
├── start.sh
└── conf/config.yml
```

构建并推送测试镜像：

```bash
cd greeter
IMAGE_REPO=harbor.example.com/grpc-demo ./deploy.sh -d 1.0.0 test

cd ../generalgateway
IMAGE_REPO=harbor.example.com/grpc-demo ./deploy.sh -d 1.0.0 test
```

构建并推送线上镜像：

```bash
cd greeter
IMAGE_REPO=harbor.example.com/grpc-demo ./deploy.sh -d 1.0.0 online

cd ../generalgateway
IMAGE_REPO=harbor.example.com/grpc-demo ./deploy.sh -d 1.0.0 online
```

需要登录私有仓库时：

```bash
export DOCKER_REGISTRY=harbor.example.com
export DOCKER_USERNAME=your-user
export DOCKER_PASSWORD=your-password

cd greeter
IMAGE_REPO=harbor.example.com/grpc-demo ./deploy.sh -d 1.0.0 test
```

默认生成：

- 测试镜像：`grpc-demo/greeter-test:VERSION`
- 线上镜像：`grpc-demo/greeter:VERSION`

网关测试配置通过 Docker 服务名访问 gRPC：

```yaml
rpcclient:
  greeter_address: "greeter:50051"
```

## 修改协议

修改 `pbdata/greeterpb/greeter.proto` 后重新生成：

```bash
protoc -I . --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  pbdata/greeterpb/greeter.proto
```

`deploy.sh` 构建前也会自动调用 `pbdata/greeterpb/pb.sh`。

## 验证

```bash
go test ./...
go vet ./...
```
# grpc-demo
