# BUG_REPRO

## Bug 是什么
读数历史查询存在切片底层数组共享问题：`readings.Store.List` 用 `items[:0]` 原地压缩共享底层数组，`FilterInPlace`/`Aggregate` 又原地修改调用方切片，容量上限保留的是最旧读数，接口 limit 多返回一条。表现为某台设备的历史里混入/丢失读数、批量上报后数据串场。

## 如何触发
1. 启动服务（`cd backend && go run .`）。
2. `POST /api/readings` 给设备 A 上报若干读数。
3. 用 `GET /api/readings?unit=A&from=...` 做一次窗口查询，再 `GET /api/readings?unit=A` 拉全量，历史条数变少或内容错乱。
4. 超过容量上限后新读数丢失（保留的是最旧读数）；`GET /api/readings?unit=A&limit=2` 返回 3 条。

## 真实错误信息
`TestReadingsListIsolated`：窗口查询后全量历史 `got 2 readings, want 3`；`TestAppendKeepsNewest`：保留的是旧读数；`TestReadingsLimitHonored`：`limit=2 returned 3 items`。
