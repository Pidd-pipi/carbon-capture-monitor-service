# BUG_REPRO

## Bug 是什么
服务生命周期存在资源与响应缺陷：`newEnterpriseServer` 把 `IdleTimeout` 清零（空闲连接不回收），`opsEnterpriseMiddleware` 用 defer 在响应提交后才写延迟头（头丢失），`web.FS` 只内嵌 `index.html` 漏了 `app.js`。表现为反复重启后内存/连接占用降不下来、响应缺 `X-Operations-Latency-Ms` 头、页面脚本 404。

## 如何触发
1. 启动服务（`cd backend && go run .`）。
2. 检查 `newEnterpriseServer` 返回的 `IdleTimeout`（为 0）。
3. 任意请求，响应头里没有 `X-Operations-Latency-Ms`。
4. `GET /app.js` → 404。

## 真实错误信息
`TestServerTimeoutsConfigured`：`IdleTimeout must be positive, got 0s`；`TestOpsLatencyHeaderPresent`：`response missing X-Operations-Latency-Ms header`；`TestAppJServed`：`expected 200 for /app.js, got 404`。
