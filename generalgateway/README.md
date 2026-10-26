# generalgateway

`generalgateway` 是 HTTP 网关，对应 `server-utils` 中的同名微服务目录。

## 监听端口

| 端口 | 用途 |
| --- | --- |
| `8080` | HTTP API |

当前接口：

```text
POST /v1/hello
GET  /status.html
```

调用链：

```text
HTTP 请求
  -> controller
  -> service
  -> grpc_client
  -> greeter:50051
```

## 编译方式

所有命令都在 `generalgateway` 目录执行：

```bash
cd generalgateway
```

| 场景 | 命令 | 使用配置 | 产物 |
| --- | --- | --- | --- |
| 本地开发 | `go run .` | `config/config.yml` | 不生成发布目录 |
| 测试环境编译 | `bash deploy.sh -b test` | `config/config_test.yml` | `docker/release` |
| 正式环境编译 | `bash deploy.sh -b online` | `config/config.yml` | `docker/release` |
| 测试镜像 | `bash deploy.sh -d VERSION test` | `config/config_test.yml` | 构建并推送测试镜像 |
| 正式镜像 | `bash deploy.sh -d VERSION online` | `config/config.yml` | 构建并推送正式镜像 |

## 测试环境编译

```bash
cd generalgateway
bash deploy.sh -b test
```

测试环境使用 `config/config_test.yml`，通过 Docker 服务名访问 gRPC 服务：

```yaml
rpcclient:
  greeter_address: "greeter:50051"
```

执行过程：

1. 调用 `../pbdata/greeterpb/pb.sh` 重新生成 protobuf 代码。
2. 使用 `CGO_ENABLED=0 GOOS=linux GOARCH=amd64` 编译 Linux 二进制。
3. 清空并重新生成 `docker/release`。
4. 把 `config/config_test.yml` 复制为 `docker/release/conf/config.yml`。
5. 把 `start.sh` 和二进制放入发布目录。

产物：

```text
generalgateway/docker/release/
├── generalgateway
├── start.sh
└── conf/config.yml
```

## 正式环境编译

```bash
cd generalgateway
bash deploy.sh -b online
```

正式环境使用 `config/config.yml`：

```yaml
server:
  address: ":8080"

gin:
  mode: "release"

rpcclient:
  greeter_address: "greeter.example.internal:50051"
```

正式编译前必须把 `rpcclient.greeter_address` 改为正式环境实际可访问的 gRPC 地址。

## 构建和推送镜像

测试镜像：

```bash
cd generalgateway
IMAGE_REPO=harbor.example.com/grpc-demo \
  bash deploy.sh -d 1.0.0 test
```

镜像名：

```text
harbor.example.com/grpc-demo/generalgateway-test:1.0.0
```

正式镜像：

```bash
cd generalgateway
IMAGE_REPO=harbor.example.com/grpc-demo \
  bash deploy.sh -d 1.0.0 online
```

镜像名：

```text
harbor.example.com/grpc-demo/generalgateway:1.0.0
```

推送私有仓库时配置登录信息：

```bash
export DOCKER_REGISTRY=harbor.example.com
export DOCKER_USERNAME=your-user
export DOCKER_PASSWORD=your-password
```

镜像构建使用 `generalgateway/docker/Dockerfile`，构建上下文是 `generalgateway/docker`。Dockerfile 直接复制 `release` 目录。

## 本地验证

在项目根目录执行：

```bash
go test ./...
go vet ./...
```

启动 `greeter` 后运行网关：

```bash
go run ./generalgateway
```

调用接口：

```bash
curl -X POST http://localhost:8080/v1/hello \
  -H "Content-Type: application/json" \
  -d '{"name":"Codex"}'
```

## 清理发布产物

```bash
cd generalgateway
bash deploy.sh -b test delete
```

注意：带 `delete` 参数时，除了删除 `docker/release`，还会删除公共 protobuf 生成代码。下一次编译会重新生成。
