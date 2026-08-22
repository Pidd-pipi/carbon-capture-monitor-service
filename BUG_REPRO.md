# BUG_REPRO

## Bug 是什么
预警状态机与其消费层状态集合不同步：`opsTransitionTable` 缺少 `paused → queued` 迁移边，`Move` 用硬编码规则同样漏掉该边，`validation.AlertStatus` 与 `routes.go` 的 `knownAlertStatuses` 都不认 `queued`。表现为暂停工单无法退回排队态、按状态筛不到排队单、带空格的合法状态被拒。

## 如何触发
1. 启动服务（`cd backend && go run .`）。
2. 创建工单并流转到 `active`，再尝试流转回 `queued` → 报「transition is not allowed」。
3. `GET /api/alerts/status/queued` → 400 unknown alert status。
4. `AlertStatus("  queued ")` → 被拒绝（不 trim）。

## 真实错误信息
`TestRewindPausedToQueued`：`paused->queued transition should be allowed: ... transition is not allowed`；`TestQueuedStatusView`：`expected 200, got 400`；`TestAlertStatusTrimsSpaces`：`whitespace not tolerated`。
