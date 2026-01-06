# zabbix-mcp-go

一个基于 [MCP (Model Context Protocol)](https://github.com/mark3labs/mcp-go) 的 Zabbix 多实例接入端。项目通过连接池复用多个 Zabbix API 客户端，统一暴露为一组 MCP 工具，便于在 IDE Copilot、LLM Agent 或自定义自动化脚本中直接完成用户、用户组、实例等常见操作。

## 🚀 功能总览
> ⚠ 下述工具仅部分功能在zabbix5/zabbix7中进行测试，其他版本未测试，请自行测试。
> 如有问题：欢迎提交 Issue/PR 或 fengzhilaoling@gmail.com中留言，共同完善 Zabbix MCP 能力！

| 领域 | MCP 工具 ID | 能力说明 | 关键参数 | 返回内容 |
|------|--------------|----------|-----------|-----------|
| 实例管理 | `get_instances_info` | 查看客户端池中全部或指定实例的连接方式、版本、占用情况 | `instance`（可选，按名称筛选） | `[]ClientInfo`，包含 URL、登录方式、是否 InUse、版本号等 |
| 用户查询 | `get_users` | 按实例列出用户，可选单个 `username` 精准过滤，并附带用户组与权限信息 | `instance`（必填）、`username`（可选） | `[]map[string]interface{}`，对应 Zabbix `user.get` 结果 |
| 用户创建 | `create_user` | 在指定实例中创建账号，自动生成高强度初始密码，可以指定角色与用户组 | `instance`、`username`、`userGroup`（必填），`name`、`roleID`（可选） | `map[string]interface{}`，附带生成的 `passwd` |
| 用户更新 | `update_user` | 修改用户姓名、所属用户组，支持一键刷新密码 | `instance`、`userid`（必填），`name`、`usrgrps[]`、`updatePasswd`（可选） | 更新后的 `user.update` 结果 |
| 用户禁用 | `disable_user` | 自动查找 "No access to the frontend" 组并把指定用户移入该组，同时重置密码 | `instance`、`userid` | `user.update` 执行结果 |
| 用户删除 | `delete_user` | 直接调用 `user.delete`，支持一次删除多个用户 ID | `instance`、`userids[]` | 删除结果集合 |
| 用户组查询 | `get_groups` | 查询用户组详情，可携带名称过滤、状态筛选，并附带成员/权限/标签过滤器等 | `instance`（必填）、`name`、`status`、`selectUsers`、`selectRights`、`selectTagFilters` | `[]map[string]interface{}`，对应 `usergroup.get` |
| 用户组创建 | `create_group` | 创建用户组并可设置主机组/模板组权限、标签过滤器与初始成员 | `instance`（必填）、`name`（必填）、`hostgroup_rights`、`templategroup_rights`、`tag_filters`、`userids` | `map[string]interface{}`，对应 `usergroup.create` 的返回结果 |
| 主机查询 | `get_hosts` | 按实例查询主机信息，支持精确或模糊匹配、以及丰富的 `select_*` 子查询（接口、模板、监控项等） | `instance`（可选）、`hostnames`/`search`、`active_available`、各类 `select_*` 标志 | `[]map[string]interface{}`，对应 Zabbix `host.get` |
| 通过主机名查找 | `find_host_by_name` | 通过主机名（数组）查询主机，支持开启 `search` 进行模糊匹配，并可控制返回字段 | `instance`、`hostnames[]`、`search`、各类 `select_*` | `[]map[string]interface{}`，对应 Zabbix `host.get` |
| 主机创建 | `create_host` | 在指定实例中创建主机，要求指定主机名与主机组；支持接口、模板链接、标签、用户宏与清单信息 | `instance`、`host`（必填）、`name`、`groups[]`(仅 `groupid`)、`interfaces[]`、`templateids[]`/`templates[]`、`tags[]`、`macros[]`、`inventory` | `map[string]interface{}`，对应 `host.create` 的返回结果（包含 hostid） |
| 主机更新 | `update_host` | 更新主机属性（替换式语义）：传入的 `groups`/`interfaces`/`templates`/`tags`/`macros` 会替换当前关联，未列出的将被移除；要求指定单个 `hostid` 进行更新 | `instance`、`hostid`（必填）、`host`、`name`、`groups[]`(仅 `groupid`)、`interfaces[]`、`templates[]`、`templates_clear[]`、`tags[]`、`macros[]`、`inventory` | `map[string]interface{}`，对应 `host.update` 的返回结果 |
| 监控项查询 | `get_items` | 按实例查询监控项，支持按主机 ID/主机名或监控项 key/name 过滤（至少需提供主机或监控项过滤之一） | `instance`（必填）、`host_ids` / `hostname`、`item_key` / `item_name` | `[]map[string]interface{}`，对应 Zabbix `item.get` 的结果 |
| 历史数据（按时间） | `get_history_by_time` | 按明确开始/结束时间范围获取历史数据，支持是否汇总与历史类型选择 | `instance`（必填）、`host_ids[]`（必填）、`item_ids[]`（必填）、`start_time`（必填）、`end_time`（必填）、`summary`、`history` | 历史数据数组或汇总结果，视 `summary` 与 `history` 而定 |
| 历史数据（按范围） | `get_history_by_range` | 按相对时间范围（例如 `7d 15h`）获取历史数据，适合快速区间查询 | `instance`（必填）、`host_ids[]`（必填）、`item_ids[]`（必填）、`time_range`（必填）、`summary`、`history` | 历史数据数组或汇总结果 |
| 历史同比/环比对比 | `get_history_compare` | 获取当前区间与前一周期（previous）对比的数据，支持按日/小时粒度与百分比格式化 | `instance`（必填）、`host_ids[]`（必填）、`item_ids[]`（必填）、`start_time`/`end_time` 或 `time_range`、`period`、`pct_format`、`timezone`、`history` | 返回包含 current 与 previous 两个时间段数据及汇总/同比变化的结构 |

> ✅ 上述工具均已在 `register/` 下完成注册，可直接通过 MCP Server 暴露给客户端。

> **其他功能补充中** 

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

### Cursor / VS Code 集成配置

> 以下示例均以 Windows 为例，路径可按需替换为自己的工作目录或用户目录。

#### Cursor（支持 stdio / SSE 双模式）

1. 打开 Cursor → `Settings` → `MCP Servers`，或直接编辑 `C:\Users\<you>\AppData\Roaming\Cursor\User\globalStorage\state.mcp.json`。
2. 根据需要添加下列配置：

| 模式 | 运行命令 | Cursor 配置片段 |
|------|-----------|------------------|
| stdio | `D:\go_code\zabbix-mcp-go\zabbixMcp.exe -stdio` | ```json
{
  "name": "zabbix-mcp-stdio",
  "type": "stdio",
  "command": "D:/go_code/zabbix-mcp-go/zabbixMcp.exe",
  "args": ["-stdio"],
  "cwd": "D:/go_code/zabbix-mcp-go"
}
``` |
| SSE/HTTP | `D:\go_code\zabbix-mcp-go\zabbixMcp.exe -http -port 5443` | ```json
{
  "name": "zabbix-mcp-sse",
  "type": "sse",
  "url": "http://127.0.0.1:5443/sse",
  "registrationUrl": "http://127.0.0.1:5443/openapi.json"
}
``` |

3. 保存后在 Cursor 的 “Available MCP Servers” 中启用即可；SSE 模式下需保持服务常驻监听。

#### VS Code / GitHub Copilot Chat（Insiders 构建）

1. 确保安装最新版 VS Code + Copilot Chat，并启用实验性的 MCP 支持（`"github.copilot.chat.enableMcp": true`）。
2. 在 VS Code 用户设置（`settings.json`）中添加：

```json
"github.copilot.mcpServers": {
  "zabbix-mcp-stdio": {
    "type": "stdio",
    "command": "D:/go_code/zabbix-mcp-go/zabbixMcp.exe",
    "args": ["-stdio"],
    "cwd": "D:/go_code/zabbix-mcp-go"
  },
  "zabbix-mcp-sse": {
    "type": "sse",
    "url": "http://127.0.0.1:5443/sse",
    "registrationUrl": "http://127.0.0.1:5443/openapi.json"
  }
}
```

3. 重启 VS Code 或重新加载窗口后，即可在 Copilot 侧边栏的 MCP 工具列表中看到 `zabbix-mcp-*`，并在对话中通过 `@zabbix-mcp-stdio` 等方式直接调用。

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


## ❤️ 支持项目 / 赞助

如果这个项目在你的日常运维或自动化工作中带来了帮助，欢迎通过扫码赞助支持持续维护：

| 渠道 | 收款码 |
|------|--------|
| 微信 | ![WeChat Pay](docs/sponsor-wechat.jpg "将你的微信收款码命名为 sponsor-wechat.png 放入 docs/ 目录") |
| 支付宝 | ![Alipay](docs/sponsor-alipay.jpg "将你的支付宝收款码命名为 sponsor-alipay.png 放入 docs/ 目录") |

> 也可以直接联系作者添加其他渠道，或在 PR/Issue 中留言。感谢支持！


## 📌 后续展望

- 扩展更多 Zabbix API（触发器、模板、媒体等）
- 加入鉴权/审计日志落库
- 引入单元测试与集成测试保障

欢迎提交 Issue/PR，共同完善 Zabbix MCP 能力！
