// Command cubic-cwnd 是 TCP CUBIC 拥塞窗口核算命令行工具。
//
// 用户以 JSON 给出上次丢包时的窗口 WMax（MSS 段数）、CUBIC 缩放
// 因子 C、RTT 与当前时间 t（或 ACK 计数），工具按
// W(t)=C(t-K)^3+WMax、K=cbrt(WMax*Beta/C) 核算当前窗口、TCP 友好
// 区与 fast convergence 状态，并可沿 RTT 逐轮仿真窗口演化。
//
// 子命令：
//
//	cubic-cwnd win <算例.json>    单点核算：W、K、友好区、恢复状态
//	cubic-cwnd curve <算例.json>  按时间采样 W(t) 输出曲线表
//	cubic-cwnd sim <算例.json>    逐 RTT 仿真（慢启动/拥塞避免、cubic/reno）
//	cubic-cwnd fair <算例.json>   双流共享瓶颈的公平性收敛仿真
//	cubic-cwnd check <算例.json>  交叉规则自查（t=K、WMax、C、友好区、Reno）
//	cubic-cwnd version            打印版本
//	cubic-cwnd help               显示帮助
//
// 算例文件路径可写 "-" 表示从 stdin 读取。所有错误写入 stderr 并
// 以非零退出码结束，非法输入绝不静默给出数值。
package main

import (
	"fmt"
	"os"
	"strings"

	"cubic-cwnd/internal/check"
	"cubic-cwnd/internal/cubic"
	"cubic-cwnd/internal/fair"
	"cubic-cwnd/internal/input"
	"cubic-cwnd/internal/sim"
)

// version 是 CLI 版本号。
const version = "1.0.0"

// usageText 是帮助文本。
const usageText = `cubic-cwnd —— TCP CUBIC 拥塞窗口核算

用法:
  cubic-cwnd win <算例.json>      单点核算（W、K、TCP 友好区、恢复状态）
  cubic-cwnd curve <算例.json>    按时间采样 W(t) 输出曲线表
  cubic-cwnd sim <算例.json>      逐 RTT 仿真（慢启动/拥塞避免，cubic/reno）
  cubic-cwnd fair <算例.json>     双流共享瓶颈的公平性收敛仿真
  cubic-cwnd check <算例.json>    交叉规则自查
  cubic-cwnd version              打印版本
  cubic-cwnd help                 显示本帮助

算例示例:
  cubic-cwnd win example/after-loss.json

算例文件可写 "-" 从 stdin 读取：
  cat example/after-loss.json | cubic-cwnd win -

输入 JSON 字段（窗口单位为 MSS 段数，时间单位为秒）:
  必填: w_max、rtt_seconds，以及 t_seconds 与 acks 二选一
  可选: name、c（默认 0.4）、previous_w_max、horizon_seconds、
        samples、sim{mode,rounds,start,initial_cwnd,ssthresh}、
        fair{capacity_segments,rounds,flow_a_cwnd,flow_b_cwnd}

关键约定:
  Beta 钉 0.7；K=cbrt(WMax*Beta/C)；有效窗口 W(t)=max(W_cubic,W_tcp)；
  低 BDP 走 TCP-friendly 分支；fast convergence 在 previous_w_max
  大于 w_max 时把参考窗口压到 (1+Beta)/2。

失败行为:
  w_max<=0、c<=0、rtt_seconds<=0、t<0、unknown 字段、JSON 非法、
  文件不存在、t 与 acks 同时或都不给出、迭代轮数超上限——一律向
  stderr 输出错误信息并以非零退出码结束。
`

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fail("缺少子命令，运行 cubic-cwnd help 查看用法")
	}
	switch args[0] {
	case "win":
		needFile(args)
		runWin(args[1])
	case "curve":
		needFile(args)
		runCurve(args[1])
	case "sim":
		needFile(args)
		runSim(args[1])
	case "fair":
		needFile(args)
		runFair(args[1])
	case "check":
		needFile(args)
		runCheck(args[1])
	case "version", "-v", "--version":
		fmt.Printf("cubic-cwnd %s\n", version)
	case "help", "-h", "--help":
		fmt.Print(usageText)
	default:
		fail("未知子命令 %q，运行 cubic-cwnd help 查看用法", args[0])
	}
}

// needFile 校验子命令携带算例文件参数。
func needFile(args []string) {
	if len(args) < 2 {
		fail("%s 需要一个算例文件参数（或 - 表示 stdin）", args[0])
	}
}

// loadSpec 按路径加载算例；路径为 "-" 时从 stdin 读取。
func loadSpec(path string) (*input.Spec, error) {
	if path == "-" {
		return input.ReadFrom(os.Stdin)
	}
	return input.LoadFile(path)
}

// runWin 执行单点核算并打印报表。
func runWin(path string) {
	spec, err := loadSpec(path)
	if err != nil {
		fail("%v", err)
	}
	res, err := cubic.ComputeFastConv(spec.ToParams(), spec.PrevWMax, true)
	if err != nil {
		fail("%v", err)
	}
	fmt.Print(formatWin(res, spec))
}

// runCurve 采样 W(t) 输出曲线表。
func runCurve(path string) {
	spec, err := loadSpec(path)
	if err != nil {
		fail("%v", err)
	}
	p := spec.ToParams()
	horizon := spec.HorizonSeconds
	n := spec.Samples
	if n < 2 {
		n = 2
	}
	step := horizon / float64(n-1)
	fmt.Printf("curve: %s (W_max=%.3f segments, K=%.6f s)\n", spec.Name, p.WMax, cubic.K(p))
	fmt.Println("t(s)      W_cubic   W_tcp     W(t)      region")
	for i := 0; i < n; i++ {
		q := p
		q.T = step * float64(i)
		region := "cubic"
		if cubic.IsTCPFriendly(q) {
			region = "tcp-friendly"
		}
		fmt.Printf("%8.4f  %9.3f  %9.3f  %9.3f  %s\n",
			step*float64(i), cubic.WCubic(q), cubic.WEst(q), cubic.WEffective(q), region)
	}
}

// runSim 逐 RTT 仿真窗口演化。
func runSim(path string) {
	spec, err := loadSpec(path)
	if err != nil {
		fail("%v", err)
	}
	mode, err := sim.ParseMode(spec.Sim.Mode)
	if err != nil {
		fail("%v", err)
	}
	start, err := sim.ParseStart(spec.Sim.Start)
	if err != nil {
		fail("%v", err)
	}
	cfg := sim.Config{
		Mode:        mode,
		Start:       start,
		Rounds:      spec.Sim.Rounds,
		InitialCwnd: spec.Sim.InitialCwnd,
		Ssthresh:    spec.Sim.Ssthresh,
		WMax:        spec.WMax,
		C:           spec.C,
		RTT:         spec.RTT,
	}
	res, err := sim.Run(cfg)
	if err != nil {
		fail("%v", err)
	}
	fmt.Printf("sim: mode=%s start=%s rounds=%d w_max=%.3f\n", mode, start, cfg.Rounds, spec.WMax)
	fmt.Print(sim.FormatTrace(res))
	fmt.Print(sim.FormatSummary(res))
}

// runFair 运行双流公平性收敛仿真。
func runFair(path string) {
	spec, err := loadSpec(path)
	if err != nil {
		fail("%v", err)
	}
	cfg := fair.Config{
		Capacity: spec.Fair.Capacity,
		Rounds:   spec.Fair.Rounds,
		FlowA:    spec.Fair.FlowA,
		FlowB:    spec.Fair.FlowB,
	}
	res, err := fair.Run(cfg)
	if err != nil {
		fail("%v", err)
	}
	fmt.Printf("fair: capacity=%.3f segments, rounds=%d, A=%.3f B=%.3f\n",
		cfg.Capacity, cfg.Rounds, cfg.FlowA, cfg.FlowB)
	fmt.Print(fair.FormatFrames(res))
	fmt.Println(res.String())
}

// runCheck 执行交叉规则自查；存在未通过规则时非零退出。
func runCheck(path string) {
	spec, err := loadSpec(path)
	if err != nil {
		fail("%v", err)
	}
	results, allPass := check.Verify(spec)
	for _, r := range results {
		fmt.Println(r.String())
	}
	if !allPass {
		fail("check: 存在未通过的规则（详见上方 FAIL 行）")
	}
	fmt.Println("check: 全部规则通过")
}

// formatWin 渲染单点核算报表。
func formatWin(res cubic.Result, spec *input.Spec) string {
	var b strings.Builder
	name := spec.Name
	if name == "" {
		name = "(unnamed)"
	}
	fmt.Fprintf(&b, "case: %s\n", name)
	fmt.Fprintf(&b, "W_max      : %.3f segments\n", res.WMaxRef)
	fmt.Fprintf(&b, "C          : %.3f\n", res.Params.C)
	fmt.Fprintf(&b, "RTT        : %.6f s\n", res.Params.RTT)
	fmt.Fprintf(&b, "t          : %.6f s (%s)\n", res.Params.T, spec.TimeOrigin())
	fmt.Fprintf(&b, "K          : %.6f s\n", res.K)
	fmt.Fprintf(&b, "W_cubic(t) : %.3f segments\n", res.WCubic)
	fmt.Fprintf(&b, "W_tcp(t)   : %.3f segments\n", res.WEst)
	fmt.Fprintf(&b, "W(t)       : %.3f segments\n", res.W)
	fmt.Fprintf(&b, "TCP-friendly region: %v\n", yesno(res.Friendly))
	fmt.Fprintf(&b, "window status       : %s\n", res.Status)
	fmt.Fprintf(&b, "fast convergence    : %v\n", yesno(res.FastConv))
	return b.String()
}

// yesno 把布尔值渲染成 yes/no。
func yesno(v bool) string {
	if v {
		return "yes"
	}
	return "no"
}

// fail 把错误写入 stderr 并以退出码 1 结束。
func fail(format string, a ...any) {
	fmt.Fprintf(os.Stderr, "cubic-cwnd: "+format+"\n", a...)
	os.Exit(1)
}
