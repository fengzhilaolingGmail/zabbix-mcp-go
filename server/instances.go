package server

import (
	"context"

	"zabbixMcp/zabbix"
)

// GetInstancesInfo 返回由 provider 提供的实例信息
// 当 instanceName 为空时返回所有实例，否则仅返回匹配的实例
func GetInstancesInfo(ctx context.Context, provider zabbix.ClientProvider, instanceName string) ([]zabbix.ClientInfo, error) {
	if provider == nil {
		return []zabbix.ClientInfo{}, nil
	}
	if ctx != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
	}
	return provider.Info(instanceName), nil
}
