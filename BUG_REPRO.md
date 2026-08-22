# BUG_REPRO

## Bug 是什么
巡检评估链路存在超时与更新缺陷：`opsContext` 忽略传入的 timeout 参数（不产生 deadline），`evaluateUnit` 在无读数时提前返回跳过设备状态告警，`store.UpdateReadings` 只更新捕集率与压力、漏掉溶剂余量与更新时间。表现为巡检超时不生效、缺读数的异常设备不再产生状态告警、设备状态停留旧值。

## 如何触发
1. 启动服务（`cd backend && go run .`）。
2. 给 `opsContext` 传一个短超时，检查返回的 ctx 是否带 deadline（不带）。
3. 让一台设备处于 `attention` 状态但无任何读数，运行 `evaluateUnit`，不产生状态告警。
4. 调用 `store.UpdateReadings` 更新溶剂余量，返回对象里溶剂余量仍是旧值。

## 真实错误信息
`TestTimeoutContextEnforced`：`opsContext returned a context without a deadline`；`TestEvaluateUnitStatusWithoutReadings`：`expected a status alert, got {trigger:false}`；`TestUpdateReadingsComplete`：`solvent level not updated`。
