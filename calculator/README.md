# calculator

`calculator` 是一个与 `greeter` 结构一致的独立 gRPC 微服务，提供整数加法和乘法计算。

## 接口

| RPC | 请求 | 结果 |
| --- | --- | --- |
| `Calculator.Add` | `{ left: int64, right: int64 }` | `{ result: int64 }` |
| `Calculator.Multiply` | `{ left: int64, right: int64 }` | `{ result: int64 }` |

结果超出 `int64` 范围时返回 gRPC `OutOfRange`。

## 端口

| 端口 | 用途 |
| --- | --- |
| `50053` | gRPC 业务端口 |
| `50054` | HTTP 状态检查，路径为 `/status.html` |

## 本地运行

在项目根目录执行：

```bash
go run ./calculator
```

检查状态：

```bash
curl http://localhost:50054/status.html
```

服务启用了 gRPC reflection，可使用支持 reflection 的客户端调用，例如：

```bash
grpcurl -plaintext -d '{"left":"21","right":"2"}' \
  localhost:50053 calculatorpb.Calculator.Multiply
```

因为 `left` 和 `right` 是 `int64`，JSON 请求中推荐使用字符串表示，避免客户端精度问题。

## 通过 generalgateway 调用

启动 `generalgateway` 后可调用：

```bash
curl -X POST http://localhost:8080/v1/calculate \
  -H "Content-Type: application/json" \
  -d '{"operation":"multiply","left":21,"right":2}'
```

`operation` 支持 `add` 和 `multiply`。HTTP 接口接受数字或字符串形式的 `int64`，结果位于 `data.result`。

## 编译和部署

| 场景 | 命令 | 使用配置 | 产物 |
| --- | --- | --- | --- |
| 本地开发 | `go run ./calculator` | `config/config.yml` | 不生成发布目录 |
| 测试环境编译 | `bash deploy.sh -b test` | `config/config_test.yml` | `docker/release` |
| 正式环境编译 | `bash deploy.sh -b online` | `config/config.yml` | `docker/release` |
| 测试镜像 | `bash deploy.sh -d VERSION test` | `config/config_test.yml` | 构建并推送测试镜像 |
| 正式镜像 | `bash deploy.sh -d VERSION online` | `config/config.yml` | 构建并推送正式镜像 |

测试镜像：

```bash
cd calculator
IMAGE_REPO=harbor.example.com/grpc-demo \
  bash deploy.sh -d 1.0.0 test
```

镜像名为 `harbor.example.com/grpc-demo/calculator-test:1.0.0`。

正式镜像：

```bash
cd calculator
IMAGE_REPO=harbor.example.com/grpc-demo \
  bash deploy.sh -d 1.0.0 online
```

镜像名为 `harbor.example.com/grpc-demo/calculator:1.0.0`。

## 修改协议

修改 `pbdata/calculatorpb/calculator.proto` 后，在项目根目录执行：

```bash
bash pbdata/calculatorpb/pb.sh
```
