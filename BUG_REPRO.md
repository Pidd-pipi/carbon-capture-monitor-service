# BUG_REPRO

## Bug 是什么
并发场景下预警列表查询与工单写入互相干扰：列表偶尔缺条目、两次刷新内容不一致、新建工单要再刷一次才出现。`ops.OpsStore` 读路径不加锁且直接返回内部引用，`SortRecords` 原地排序共享底层数组，`alertsCollection` 又用未失效的列表缓存把旧快照回给前端。

## 如何触发
1. 启动服务（`cd backend && go run .`）。
2. 多个请求并发：一边 `POST /api/alerts` 创建工单，一边 `GET /api/alerts` 刷新列表。
3. 或按序触发缓存问题：先 `GET /api/alerts?page=1&page_size=50`，再创建一条工单，再 `GET /api/alerts?page=1&page_size=50`，第二次列表看不到新工单。

## 真实错误信息
`go run -race` 复现并发场景输出：

```
==================
WARNING: DATA RACE
Write at 0x00c0001980f8 by goroutine 10:
  example.com/carbon-capture-monitor-service/httpapi.(*server).alertsCollection()
      .../httpapi/alerts_handlers.go:43
Previous read at 0x00c0001980f8 by goroutine 8:
  example.com/carbon-capture-monitor-service/httpapi.(*server).alertsCollection()
      .../httpapi/alerts_handlers.go:34
==================
Found 1 data race(s)
```

存储层并发读写同样报 `WARNING: DATA RACE`（`OpsStore.List` 与 `OpsStore.Put`）。
