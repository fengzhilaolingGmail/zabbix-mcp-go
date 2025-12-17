package handler

import "zabbixMcp/zabbix"

// 持有可选的客户端池引用，main 初始化后会调用 SetClientPool 注入
var clientPool ZabbixClientPool

// SetClientPool 注入全局池引用（可为 nil）
func SetClientPool(p ZabbixClientPool) {
	clientPool = p
}

// ZabbixClient 是 handler 层对客户端的抽象（便于测试与替换实现）
// type ZabbixClient interface {
// 	// Call 用于执行任意 Zabbix API 方法
// 	Call(method string, params interface{}) (interface{}, error)
// 	// GetUsers 是 handler 期望的高层便利方法（实现可通过 Call 完成）
// 	GetUsers(params map[string]interface{}) ([]map[string]interface{}, error)
// }

// ZabbixClientPool 是 handler 层的连接池抽象
type ZabbixClientPool interface {
	Get() (*zabbix.ZabbixClient, error)
	Release(client *zabbix.ZabbixClient) error
	Info() []zabbix.ClientInfo
}
