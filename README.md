# zabbix-mcp-go

一个基于 [MCP (Model Context Protocol)](https://github.com/mark3labs/mcp-go) 的 Zabbix 多实例接入端。项目通过连接池复用多个 Zabbix API 客户端，统一暴露为一组 MCP 工具，便于在 IDE Copilot、LLM Agent 或自定义自动化脚本中直接完成用户、用户组、实例等常见操作。

## 🚀 功能总览

| 领域 | MCP 工具 ID | 能力说明 | 关键参数 | 返回内容 |
|------|--------------|----------|-----------|-----------|
| 实例管理 | `get_instances_info` | 查看客户端池中全部或指定实例的连接方式、版本、占用情况 | `instance`（可选，按名称筛选） | `[]ClientInfo`，包含 URL、登录方式、是否 InUse、版本号等 |
| 用户查询 | `get_users` | 按实例列出用户，可选单个 `username` 精准过滤，并附带用户组与权限信息 | `instance`（必填）、`username`（可选） | `[]map[string]interface{}`，对应 Zabbix `user.get` 结果 |
| 用户创建 | `create_user` | 在指定实例中创建账号，自动生成高强度初始密码，可以指定角色与用户组 | `instance`、`username`、`userGroup`（必填），`name`、`roleID`（可选） | `map[string]interface{}`，附带生成的 `passwd` |
| 用户更新 | `update_user` | 修改用户姓名、所属用户组，支持一键刷新密码 | `instance`、`userid`（必填），`name`、`usrgrps[]`、`updatePasswd`（可选） | 更新后的 `user.update` 结果 |
| 用户禁用 | `disable_user` | 自动查找 "No access to the frontend" 组并把指定用户移入该组，同时重置密码 | `instance`、`userid` | `user.update` 执行结果 |
| 用户删除 | `delete_user` | 直接调用 `user.delete`，支持一次删除多个用户 ID | `instance`、`userids[]` | 删除结果集合 |
| 用户组查询 | `get_groups` | 查询用户组详情，可携带名称过滤、状态筛选，并附带成员/权限/标签过滤器等 | `instance`（必填）、`name`、`status`、`selectUsers`、`selectRights`、`selectTagFilters` | `[]map[string]interface{}`，对应 `usergroup.get` |

> ✅ 上述工具均已在 `register/` 下完成注册，可直接通过 MCP Server 暴露给客户端。

## 🧩 架构速览

- **配置解析 (`config.go`)**：从 `config.yml` 读取多个 Zabbix 实例，支持密码/Token 双认证以及默认实例标记。
- **客户端池 (`zabbix/pool.go`)**：按实例构建可重用客户端，具备按名称借用、健康检查与版本缓存能力。
- **适配层 (`models/` + `zabbix/version.go`)**：通过 `ParamSpec` + `AdaptAPIParams` 自动适配不同 Zabbix 版本的字段差异，并在 delete 场景下输出原生 `[]string`。
- **业务服务 (`server/`)**：封装 user/host/instance 等领域方法，负责租借客户端、调用 API、记录日志。
- **MCP Handler (`handler/` + `register/`)**：解析工具入参、组合参数结构，最后以统一 JSON 结构输出。
- **日志与密码工具 (`logger/`, `utils/proc.go`)**：Zap 日志，附带高强度密码生成器，确保用户创建/禁用时始终可用。

## ⚙️ 配置

在根目录创建或编辑 `config.yml`：

```yaml
instances:
  - name: "demo-prod"
    url: "https://zabbix.example.com/api_jsonrpc.php"
    auth_type: "password"
    username: "admin"
    password: "s3cr3t"
  - name: "demo-token"
    url: "https://zbx-token.example.com/api_jsonrpc.php"
    auth_type: "token"
    token: "<your_token_here>"
    default: true
```

> `auth_type` 可选 `password` / `token`；如果配置 `default: true`，在客户端池信息查询时会标记该实例。

## 🏃‍♂️ 运行

```bash
# 安装依赖（首次）
go mod tidy

# 构建
+go build -o zabbixMcp.exe

# 以 stdio 模式运行（适合集成至编辑器插件）
./zabbixMcp.exe -stdio

# 以 HTTP/SSE 模式启动（默认端口 5443）
./zabbixMcp.exe -http -port 5443 -loglevel debug
```

程序启动后会：
1. 读取 `config.yml`、初始化客户端池并检测版本；
2. 创建 MCP Server，并注册全部工具；
3. 根据命令行参数选择 stdio / HTTP / 双通道运行方式。

## 🧪 开发与调试

```bash
# 运行所有单元测试（当前以编译通过为主）
go test ./...

# gofmt 格式化
+gofmt -w ./handler ./models ./server ./zabbix
```

### 日志定位
- 日志默认输出在控制台，如需文件输出可扩展 `logger/logger.go`。
- 所有 API 调用均带有“调用方法 + 参数 + 错误”日志，便于追踪。

## 📁 目录概览

```
├── handler/            # MCP 工具处理器
├── register/           # MCP 工具注册入口
├── server/             # 业务服务层（user/host/instances）
├── models/             # ParamSpec 定义 & 构造器
├── zabbix/             # 客户端、连接池、版本探测
├── utils/              # 辅助工具（如密码生成）
├── logger/             # zap 日志包装
├── config.go|yml       # 多实例配置加载
├── main.go             # 程序入口，负责启动 MCP server
└── README.md           # 当前文档
```

## 📌 后续展望

- 扩展更多 Zabbix API（触发器、模板、媒体等）
- 加入鉴权/审计日志落库
- 引入单元测试与集成测试保障

欢迎提交 Issue/PR，共同完善 Zabbix MCP 能力！
