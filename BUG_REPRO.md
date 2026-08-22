# BUG_REPRO

## Bug 是什么
设备状态更新链路的错误链断裂：`handlers.status` 用 `%v` 包装仓储错误导致 `errors.Is` 失效，`domain.ValidateStatusUpdate` 用 `%v` 包装领域错误，`health` 的依赖探测接错开关。表现为更新不存在的设备返回 500（应为 404）、维护态直接转在线返回 500（应为 400）、健康检查把依赖状态误报为 unknown。

## 如何触发
1. 启动服务（`cd backend && go run .`）。
2. `POST /api/capture-units/status` body `{"id":"CC-NOPE","status":"online"}` → 返回 500 而不是 404。
3. 先把 CC-ALPHA 置为 `maintenance`，再置为 `online` → 返回 500 而不是 400。
4. `GET /healthz` → 响应里 `"store":"unknown"` 而不是 `"store":"ok"`。

## 真实错误信息
`TestUnknownUnitLookup`：`expected 404, got 500`；`TestMaintenanceToOnlineRejected`：`expected 400, got 500`；`TestHealthReportsStoreOK`：`health body: {"store":"unknown",...}`。
