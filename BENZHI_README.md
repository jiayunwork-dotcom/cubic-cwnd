# cubic-cwnd：Go TCP CUBIC 拥塞窗口核算命令行与 HTTP API 服务

cubic-cwnd 读入丢包后的 WMax、C、RTT 与时间（JSON），按 W(t)=C(t−K)³+WMax 与 K=∛(WMax·β/C) 核算当前窗口、TCP 友好区与 fast convergence，并可逐 RTT 仿真 cubic/reno 增长与双流公平性收敛。

## 构建 / 运行 / 测试

```text
go build ./...
go run . win example/after-loss.json
go test ./...
```

其他子命令见项目 `README.md`：`curve`、`sim`、`fair`、`check`、`version`、`help`。算例文件可写 `-` 从 stdin 读取。

## 评测镜像

本目录评测专用文件（勿覆盖项目自带 Dockerfile/README）：

- `benzhi.Dockerfile`
- `build_benzhi_docker.sh`
- `BENZHI_README.md`（本文件）

两种架构都要构建并进容器验证：

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh <image-name> linux/arm64
./build_benzhi_docker.sh <image-name> linux/amd64
docker run -it <image-name>:latest
```
