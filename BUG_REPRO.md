# BUG_REPRO

## Bug 是什么
读数汇总链路的 context 取消未正确传播：`readingsSummary` handler 用 `context.Background()` 替代请求上下文，`readings.Store.Summary` 用 `context.Background()` 调用 `List`，`monitor.applyResult` 在取消后仍用后台上下文重试创建告警。表现为客户端取消后请求仍继续计算、取消后仍产生多余告警。

## 如何触发
1. 启动服务（`cd backend && go run .`）。
2. 发送一个请求上下文已取消的 `GET /api/readings/summary?unit=CC-ALPHA&window_minutes=60`，接口返回 200 而不是 499。
3. 对监控器传入已取消的上下文运行周期，仍会创建告警。

## 真实错误信息
`TestSummaryPropagatesClientCancel`：`expected 499, got 200`；`TestSummaryHonorsCancel`：`Summary completed despite a cancelled context`；`TestApplyResultHonorsCancel`：`alert was created despite cancelled context: 1 alerts`。
