package server

import (
	"context"
	"fmt"
	"zabbixMcp/logger"
	"zabbixMcp/models"
	"zabbixMcp/zabbix"
)

// 用户组: 所有用户组 获取用户组user
// 获取用户组信息
func GetUserGroups(ctx context.Context, provider zabbix.ClientProvider, spec models.ParamSpec, instance string) ([]map[string]interface{}, error) {
	if provider == nil {
		return nil, fmt.Errorf("no zabbix client")
	}
	var (
		lease zabbix.ClientLease
		err   error
	)
	if instance != "" {
		lease, err = provider.AcquireByInstance(ctx, instance)
	} else {
		lease, err = provider.Acquire(ctx)
	}
	if err != nil {
		return nil, err
	}
	var callErr error
	defer func() { lease.Release(callErr) }()
	client := lease.Client()
	adapted := client.AdaptAPIParams("usergroup.get", spec)
	var userGroups []map[string]interface{}
	callErr = client.Call(ctx, "usergroup.get", adapted, &userGroups)
	if callErr != nil {
		logger.L().Error("get user group error: %s", callErr.Error())
		return nil, callErr
	}
	return userGroups, nil
}
