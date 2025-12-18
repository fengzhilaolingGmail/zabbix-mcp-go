package server

import "zabbixMcp/zabbix"

// GetInstancesInfo 返回由 client 提供的实例信息
// 当 instanceName 为空时返回所有实例，否则仅返回匹配的实例
func GetInstancesInfo(client zabbix.ZabbixClientHandler, instanceName string) ([]zabbix.ClientInfo, error) {
	infos := client.Info(instanceName)
	return infos, nil
}
