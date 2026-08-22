# Carbon Capture Monitor Backend

Go 模块：`example.com/carbon-capture-monitor-service`。

```bash
go test ./...
go build ./...
go run .
```

入口为 `main.go`，默认监听 `:8080`，可用 `PORT` 调整。健康检查为 `GET /healthz`；API 为 `GET /api/capture-units` 和 `POST /api/capture-units/status`。
