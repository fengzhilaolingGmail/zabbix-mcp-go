package handler

import "zabbixMcp/zabbix"

// 持有可选的客户端池引用，main 初始化后会调用 SetClientPool 注入
// 现在使用 zabbix.ZabbixClientHandler 接口，隐藏底层具体类型
var clientPool zabbix.ZabbixClientHandler

// SetClientPool 注入全局池引用（可为 nil）
func SetClientPool(p zabbix.ZabbixClientHandler) {
	clientPool = p
}
