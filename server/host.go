/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2025-12-18 11:13:06
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2025-12-25 11:03:09
 * @FilePath: \zabbix-mcp-go\server\host.go
 * @Description: 主机相关功能
 * @Copyright: Copyright (c) 2025 by fengzhilaoling@gmail.com, All Rights Reserved.
 */

package server

import (
	"context"
	"fmt"

	"zabbixMcp/logger"
	"zabbixMcp/models"
	"zabbixMcp/zabbix"
)

// TODO 获取主机  host.create
func GetHosts(ctx context.Context, provider zabbix.ClientProvider, spec models.ParamSpec, instance string) ([]map[string]interface{}, error) {
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
	adapted := client.AdaptAPIParams("host.get", spec)
	var hosts []map[string]interface{}
	callErr = client.Call(ctx, "host.get", adapted, &hosts)
	if callErr != nil {
		logger.L().Error("get hosts error: %s", callErr.Error())
		return nil, callErr
	}
	return hosts, nil
}

// CreateHost 调用底层 ClientProvider 执行 host.create，并返回创建结果。
// instance 为空时使用任意可用客户端，否则强制选择指定实例。
func CreateHost(ctx context.Context, provider zabbix.ClientProvider, spec models.ParamSpec, instance string) (map[string]interface{}, error) {
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
	adapted := client.AdaptAPIParams("host.create", spec)
	var hosts map[string]interface{}
	callErr = client.Call(ctx, "host.create", adapted, &hosts)
	if callErr != nil {
		logger.L().Error("create host error: %s", callErr.Error())
		return nil, callErr
	}
	return hosts, nil
}

// UpdateHost 调用底层 ClientProvider 执行 host.update，并返回更新结果。
// instance 为空时使用任意可用客户端，否则强制选择指定实例。
func UpdateHost(ctx context.Context, provider zabbix.ClientProvider, spec models.ParamSpec, instance string) (map[string]interface{}, error) {
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
	adapted := client.AdaptAPIParams("host.update", spec)
	var hosts map[string]interface{}
	callErr = client.Call(ctx, "host.update", adapted, &hosts)
	if callErr != nil {
		logger.L().Error("update host error: %s", callErr.Error())
		return nil, callErr
	}
	return hosts, nil
}

// TODO 新建 主机 host.delete
// TODO 删除 主机 host.get
// TODO 更新 主机 host.update
// TODO 将相关objects添加到主机  host.massadd
// TODO 将相关objects从主机中移除 host.massremove
// TODO 从主机中替换或移除相关objects host.massupdate
