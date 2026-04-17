# Secondary Development / 二开说明

这个 fork 用来保存当前线上使用的二开补丁，尽量保持与 upstream 接近，方便后续拉取上游更新后做 diff 对照。

## 当前二开内容

1. 邮箱域名黑名单
   - 新增风控配置字段 `email_domain_blacklist`
   - 支持填写 `example.com`、`@example.com`、`.example.com`
   - 同时拦截根域名和子域名邮箱（例如 `user@sub.example.com`）
2. 注册与游客下单风控
   - 注册接口会校验邮箱域名黑名单
   - 游客下单同样会校验邮箱域名黑名单
3. 文案与测试
   - 补充了对应错误映射和中英文本
   - 增加了风控与注册拦截测试

## 本次主要改动文件

- `internal/service/order_risk_control_setting.go`
- `internal/service/order_risk_control_service.go`
- `internal/service/user_auth_service.go`
- `internal/http/handlers/public/user_auth.go`
- `internal/http/handlers/public/error_mapping.go`
- `internal/i18n/messages.go`
- `internal/service/order_risk_control_test.go`
- `internal/service/user_auth_service_registration_blacklist_test.go`

## 同步上游建议

```bash
git fetch upstream
git log --oneline upstream/main..origin/main
git diff upstream/main..origin/main
```

说明：线上手工清理 `@example.com` 垃圾账号属于运维动作，不属于仓库代码变更，不记录在本 repo 内。
