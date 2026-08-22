# BUG_REPRO

## Bug 是什么
预警工单流转接口的错误链断裂：`OpsService.Transition` 用 `%v` 包装仓储错误，`ErrorCode` 用哨兵直接比较，handler 丢失 transition 映射。导致不存在的工单流转返回 400（应为 404）、版本冲突返回 400（应为 409）、非法流转返回 400（应为 422）。

## 如何触发
1. 启动服务（`cd backend && go run .`）。
2. `POST /api/alerts/al-not-exist/transition` body `{"expected_revision":1,"target_status":"active"}` → 返回 400 而不是 404。
3. 创建工单后用错误的 `expected_revision` 流转 → 返回 400 而不是 409。
4. 把 active 工单流转回 queued → 返回 400 而不是 422。

## 真实错误信息
`TestTransitionUnknownAlert`：`expected 404, got 400`；`TestConflictStaleRevision`：`expected 409, got 400`；`TestIllegalTransition`：`expected 422, got 400`。
