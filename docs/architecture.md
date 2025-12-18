# 架构概览

本项目的 Zabbix 访问层通过 `APIClient` + `ClientProvider` 两级抽象隐藏所有底层细节：

- **APIClient**：最小调用接口，暴露 `Call(ctx, method, params, result)`，并提供版本特性查询与参数适配能力。`ZabbixClient` 是默认实现，内部负责 HTTP 认证模式探测、版本缓存和重试策略。
- **ClientProvider**：负责提供 APIClient 的安全租借句柄。`ClientPool` 实现了容量受控的连接池，`Acquire`/`Release` 语义确保不会发生并发复用；同时 `Info` 提供运行期可观测数据。

在 `handler -> server -> provider` 的调用链中：

1. handler 从 `ClientProvider` 租借客户端，将请求上下文透传到 server 层；
2. server 只关注业务输入输出，通过 `APIClient` 执行具体的 Zabbix API；
3. pool 负责生命周期管理、健康检查及实例元信息暴露。

这样可以：

- 隔离 HTTP/认证/版本兼容逻辑，业务层无需了解 Zabbix 具体协议差异；
- 方便在测试中注入假实现（mock `ClientProvider` 即可）；
- 为后续添加装饰器（指标、重试、熔断等）提供统一挂载点。
