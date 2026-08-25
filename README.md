# cubic-cwnd

cubic-cwnd 是一个 TCP CUBIC 拥塞窗口核算命令行工具。用户以 JSON 给出上次丢包时的窗口 WMax（MSS 段数）、CUBIC 缩放因子 C、RTT 与当前时间 t（或 ACK 计数），工具按 W(t)=C(t−K)³+WMax 与 K=∛(WMax·β/C) 核算当前窗口、立方区间的长度 K、TCP 友好区（RFC 8312 低 BDP 分支）、fast convergence 与恢复状态，并可沿 RTT 逐轮仿真慢启动/拥塞避免阶段转换、对比 Reno 每 RTT +1 的线性增长，以及双流共享瓶颈下的公平性收敛。它是拥塞控制内核核算，不做网络抓包、不做防火墙产品，也不是令牌桶限速工具。

## 用法

```text
go run . win example/after-loss.json
```

`win` 打印 W(t)、K、TCP 友好区与恢复状态。其它子命令：

```text
cubic-cwnd curve <算例.json>    按时间采样 W(t) 输出曲线表
cubic-cwnd sim <算例.json>      逐 RTT 仿真（慢启动/拥塞避免，cubic/reno）
cubic-cwnd fair <算例.json>     双流共享瓶颈的公平性收敛仿真
cubic-cwnd check <算例.json>    交叉规则自查
cubic-cwnd version / help
```

算例文件路径写 `-` 时从 stdin 读取：

```text
cat example/after-loss.json | cubic-cwnd win -
```

### 算例格式

`example/after-loss.json` 是预置的丢包恢复算例（w_max=16 段、c=0.4、RTT=0.1s、t=0.4s），可直接运行：

```json
{
  "name": "after-loss reference",
  "c": 0.4,
  "w_max": 16.0,
  "rtt_seconds": 0.1,
  "t_seconds": 0.4,
  "horizon_seconds": 8.0,
  "samples": 41,
  "sim": { "mode": "cubic", "rounds": 60, "start": "after-loss" },
  "fair": { "capacity_segments": 60, "rounds": 200, "flow_a_cwnd": 45.0, "flow_b_cwnd": 15.0 }
}
```

- `w_max`、`rtt_seconds` 必填；`t_seconds` 与 `acks` 二选一（`acks` 时 t = acks×RTT）。
- `c` 可选，默认 0.4；`name`、`previous_w_max` 可选。
- `horizon_seconds`（默认 1.0）与 `samples`（默认 41）控制 `curve` 采样。
- `sim` 可选：`mode`（cubic/reno）、`rounds`、`start`（after-loss/fresh）、`initial_cwnd`、`ssthresh`。
- `fair` 可选：`capacity_segments`、`rounds`、`flow_a_cwnd`、`flow_b_cwnd`。
- `example/high-bdp.json` 是高 BDP 算例（w_max=200 段、RTT=0.2s），窗口主要由立方曲线爬回。

## 关键约定

- β 钉 0.7；K = ∛(WMax·β/C)；在 t=K 处 W=WMax 且斜率为 0；t>K 后超立方超过 WMax；C 加倍则 K 变短为 1/∛2。
- 有效窗口 W(t) = max(W_cubic(t), W_est(t))，其中 W_est = WMax·β + 3(1−β)/(1+β)·t/RTT；立方曲线低于 Reno 估计时走 TCP-friendly 分支。
- fast convergence：当 `previous_w_max` 大于 `w_max` 时，参考窗口被压到 (1+β)/2 倍，K 相应变短。
- 丢包后窗口先落到 β·WMax，再沿立方曲线爬回 WMax；Reno 模式在拥塞避免期每 RTT 恰好 +1 段。
- 迭代仿真均有轮数上限（100000），超限报错而非挂死。

## 失败行为

`w_max`≤0、`c`≤0、`rtt_seconds`≤0、`t_seconds`<0、`t_seconds` 与 `acks` 同时或都不给出、`previous_w_max` 为负、未知 JSON 字段、JSON 非法、文件不存在、仿真轮数超上限——一律向 stderr 输出错误信息并以非零退出码结束，绝不静默给出数值。

## 构建与测试

```text
go build ./...
go test ./...
go run . win example/after-loss.json
```

## 目录

- `internal/cubic/` — 立方曲线 W(t)、K、TCP-friendly 估计、fast convergence 与恢复状态
- `internal/sim/` — 逐 RTT 仿真（慢启动/拥塞避免阶段转换，cubic/reno 增长律）
- `internal/fair/` — 双流共享瓶颈的 AIMD 公平性收敛仿真
- `internal/check/` — 交叉规则自查（t=K、WMax 加大、C 加倍、友好区、Reno +1）
- `internal/input/` — 算例 JSON 的严格解析与校验
- `example/` — 离线小算例
