# greeter

`greeter` 是本项目中的 gRPC 服务，对应 `server-utils` 中的独立微服务目录。

## 端口

| 端口 | 用途 |
| --- | --- |
| `50051` | gRPC 业务端口 |
| `50052` | HTTP 状态检查，路径为 `/status.html` |

## 编译方式

所有命令都在 `greeter` 目录执行：

```bash
cd greeter
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
cd greeter
bash deploy.sh -b test
```

执行过程：

1. 调用 `../pbdata/greeterpb/pb.sh` 重新生成 protobuf 代码。
2. 使用 `CGO_ENABLED=0 GOOS=linux GOARCH=amd64` 编译 Linux 二进制。
3. 清空并重新生成 `docker/release`。
4. 把 `config/config_test.yml` 复制为 `docker/release/conf/config.yml`。
5. 把 `start.sh` 和二进制放入发布目录。

产物：

```text
greeter/docker/release/
├── greeter
├── start.sh
└── conf/config.yml
```

## 正式环境编译

```bash
cd greeter
bash deploy.sh -b online
```

正式环境与测试环境使用同一编译过程，区别是配置来源：

```text
config/config.yml -> docker/release/conf/config.yml
```

正式编译前，应先确认 `config/config.yml` 中以下地址符合正式环境：

```yaml
rpcserver:
  protocol: "tcp"
  address: ":50051"

healthserver:
  address: ":50052"
```

## 构建和推送镜像

测试镜像：

```bash
cd greeter
IMAGE_REPO=harbor.example.com/grpc-demo \
  bash deploy.sh -d 1.0.0 test
```

镜像名：

```text
harbor.example.com/grpc-demo/greeter-test:1.0.0
```

正式镜像：

```bash
cd greeter
IMAGE_REPO=harbor.example.com/grpc-demo \
  bash deploy.sh -d 1.0.0 online
```

镜像名：

```text
harbor.example.com/grpc-demo/greeter:1.0.0
```

推送私有仓库时配置登录信息：

```bash
export DOCKER_REGISTRY=harbor.example.com
export DOCKER_USERNAME=your-user
export DOCKER_PASSWORD=your-password
```

镜像构建使用 `greeter/docker/Dockerfile`，构建上下文是 `greeter/docker`。Dockerfile 直接复制 `release` 目录。

## 本地验证

在项目根目录执行：

```bash
go test ./...
go vet ./...
```

启动服务：

```bash
go run ./greeter
```

检查状态：

```bash
curl http://localhost:50052/status.html
```

## 清理发布产物

```bash
cd greeter
bash deploy.sh -b test delete
```

注意：带 `delete` 参数时，除了删除 `docker/release`，还会删除自动生成的：

- `../pbdata/greeterpb/greeter.pb.go`
- `../pbdata/greeterpb/greeter_grpc.pb.go`

下一次编译会重新生成这些文件。
