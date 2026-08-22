# Carbon Capture Monitor Service

碳捕集设备读数与运行状态监测 HTTP 服务，展示捕集率、压力、溶剂余量，允许值班员更新设备状态，并维护设备预警工单与历史读数汇总。

## Layout

`backend/config`、`backend/domain`、`backend/store`、`backend/validation`、`backend/health`、`backend/ops`、`backend/readings`、`backend/httpapi`、`backend/web` 分别负责配置、模型、存储、校验、健康检查、预警工单领域、读数历史、接口和内嵌浏览器页面；`backend/monitor*.go` 是后台巡检调度。

## Run and API

执行 `cd backend && go run .` 启动，默认端口为 `8080`，可通过 `PORT` 修改。后台巡检间隔由 `MONITOR_INTERVAL_SECONDS` 控制（默认 30 秒）。

### Capture units

- `GET /healthz`
- `GET /api/capture-units`
- `POST /api/capture-units/status`，JSON body：`{"id":"CC-BETA","status":"online"}`
- `GET /`：监测页面

状态可选 `online`、`attention`、`offline`、`maintenance`。

### Alerts（预警工单）

- `GET /api/alerts?subject=&status=&priority=&owner=&page=&page_size=`
- `POST /api/alerts`，body：`{"subject":"...","owner":"...","priority":"high","labels":{"site":"..."}}`
- `POST /api/alerts/{id}/transition`，body：`{"expected_revision":1,"target_status":"active"}`
- `GET /api/alerts/{id}/audit`
- `GET /api/alerts/snapshot`
- `GET /api/alerts/rules?severity=critical`

工单状态：`queued` → `active` → `paused` → `closed`，`closed` 为终态；更新采用乐观锁（`expected_revision`）。

### Readings（历史读数）

- `POST /api/readings`，body：`{"unit_id":"CC-ALPHA","capture_rate_pct":88.1,"pressure_kpa":179.0,"solvent_level_pct":70}` 或批量 `{"unit_id":"CC-ALPHA","readings":[{...}]}`
- `GET /api/readings?unit=CC-ALPHA&from=&to=&limit=`
- `GET /api/readings/summary?unit=CC-ALPHA&window_minutes=60`

读数保留窗口由 `READING_RETENTION_HOURS`（默认 24 小时）与 `MAX_READINGS_PER_UNIT`（默认 2000）控制。

## Directory Tree

```text
backend/       Go 模块、HTTP 服务、领域包和内嵌 web 静态资源
database/      数据库说明
output/        验证记录
README.md      项目说明
prompt.txt     任务提示
runtime_smoke.json  运行冒烟配置
```

## Verification

- `gofmt -w .`: passed
- `go build ./...`: passed
- `go test ./...`: passed
- Runtime smoke: health returned 200 with `status=ok`; collection returned 200 with 3 units; valid status update returned 200; unknown unit returned 404; alert create/list/transition/audit/snapshot returned 200; readings post/list/summary returned 200; `/` and `/app.js` returned 200.
