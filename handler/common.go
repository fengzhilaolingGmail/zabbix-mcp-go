package handler

import "zabbixMcp/zabbix"

// 持有可选的客户端池引用，main 初始化后会调用 SetClientPool 注入
// 现在使用 zabbix.ZabbixClientHandler 接口，隐藏底层具体类型
var clientPool zabbix.ZabbixClientHandler

// SetClientPool 注入全局池引用（可为 nil）
func SetClientPool(p zabbix.ZabbixClientHandler) {
	clientPool = p
}

// makeResult 将任意数据包装为统一的 JSON 构造，便于所有 MCP 工具返回一致的格式。
// 返回结构为: {"ok": true, "data": <payload>}。
// makeResult 将任意数据包装为统一的 JSON 构造（map[string]interface{}），便于所有 MCP 工具返回一致的格式。
// 返回结构为: {"ok": true, "data": <payload>}。
func makeResult(data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"ok":   true,
		"data": data,
	}
}
