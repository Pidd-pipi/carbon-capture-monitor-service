# BUG_REPRO

## Bug 是什么
后台巡检循环存在并发生命周期缺陷：`WaitGroup` 计数在 goroutine 内执行、worker 错误分支不关闭结果通道、周期内单元失败被静默吞掉、审计事件无上限增长。表现为巡检跑一阵子后预警不再更新、链路卡住，且内存持续增长。

## 如何触发
1. 启动服务（`cd backend && go run .`）。
2. 让某个巡检周期内单元评估失败（例如请求上下文被取消），再观察后续周期是否还能产生预警。
3. 长期运行后观察审计事件（`GET /api/alerts/{id}/audit`）数量是否无上限增长。

## 真实错误信息
并发巡检在 `-race` 下复现：

```
WARNING: DATA RACE
...
panic: send on closed channel
```

`TestMonitorCycleSurfacesUnitErrors` 复现错误被吞：`runCycle` 对已取消的上下文返回 `nil`，`LastError()` 也为 `nil`。
