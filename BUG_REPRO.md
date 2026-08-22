# BUG_REPRO

## Bug 是什么
缺省配置路径存在 typed-nil 与 nil map 问题：`config.NewLabelsProvider` 在未配置 `ALERT_DEFAULT_LABELS` 时返回 typed-nil 接口，`configureDefaultLabels` 判空恒真后把 nil map 写进 `ops.defaultLabels`，`NormalizeRecord` 对无标签记录写入 nil map 触发 panic。

## 如何触发
1. 不设置 `ALERT_DEFAULT_LABELS` 环境变量启动服务（`cd backend && go run .`）。
2. `POST /api/alerts` body `{"subject":"x","owner":"a","priority":"low"}`（不带 labels）。
3. 请求直接 panic 崩溃。

## 真实错误信息
```
panic: assignment to entry in nil map

goroutine 1 [running]:
example.com/carbon-capture-monitor-service/ops.NormalizeRecord(...)
	.../ops/ops_model.go:119
main.main()
	.../cmd_repro.go:11
exit status 2
```
