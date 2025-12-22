/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2025-12-18 11:10:11
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2025-12-22 13:35:15
 * @FilePath: \zabbix-mcp-go\server\user.go
 * @Description: 用户相关功能
 * @Copyright: Copyright (c) 2025 by fengzhilaoling@gmail.com, All Rights Reserved.
 */

package server

import (
	"context"
	"fmt"

	"zabbixMcp/logger"
	"zabbixMcp/models"
	"zabbixMcp/utils"
	"zabbixMcp/zabbix"
)

// GetUsers 调用底层 ClientProvider 执行 user.get，并返回解析后的列表。
// instanceName 为空时使用任意可用客户端，否则强制选择指定实例。
func GetUsers(ctx context.Context, provider zabbix.ClientProvider, spec models.ParamSpec, instance string) ([]map[string]interface{}, error) {
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
	adapted := client.AdaptAPIParams("user.get", spec)
	var users []map[string]interface{}
	callErr = client.Call(ctx, "user.get", adapted, &users)
	if callErr != nil {
		logger.L().Error("get user error: %s", callErr.Error())
		return nil, callErr
	}
	return users, nil
}

// 创建用户
func CreateUsers(ctx context.Context, provider zabbix.ClientProvider, spec models.ParamSpec, instance, passwd string) (map[string]interface{}, error) {
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
	adapted := client.AdaptAPIParams("user.create", spec)
	var users map[string]interface{}
	callErr = client.Call(ctx, "user.create", adapted, &users)
	if callErr != nil {
		logger.L().Error("create user error: %s", callErr.Error())
		return nil, callErr
	}
	users["passwd"] = passwd
	return users, nil
}

func UpdateUser(ctx context.Context, provider zabbix.ClientProvider, spec models.ParamSpec, instance, passwd string) (map[string]interface{}, error) {
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
	adapted := client.AdaptAPIParams("user.update", spec)
	var users map[string]interface{}
	callErr = client.Call(ctx, "user.update", adapted, &users)
	if callErr != nil {
		logger.L().Error("update user error: %s", callErr.Error())
		return nil, callErr
	}
	if passwd != "" {
		users["passwd"] = passwd
	}
	return users, nil
}

// 禁用用户
func DisableUser(ctx context.Context, provider zabbix.ClientProvider, userId, instance string) (map[string]interface{}, error) {
	// 群组设置为：No access to the frontend
	// 查找No access to the frontend 群组id
	groupSpec := models.UserGroup{
		Output: "extend",
		Status: 0,
		Filter: map[string]interface{}{"name": "No access to the frontend"},
	}
	groups, err := GetUserGroups(ctx, provider, groupSpec, instance)
	if err != nil {
		logger.L().Error("获取\"No access to the frontend\"用户组失败: %s", err.Error())
		return nil, err
	}
	if len(groups) == 0 {
		return nil, fmt.Errorf("未找到 \"No access to the frontend\" 用户组")
	}
	fmt.Println(groups)
	var targetGroupID string
	for _, g := range groups {
		if id, ok := g["usrgrpid"].(string); ok && id != "" {
			targetGroupID = id
			break
		}
	}
	if targetGroupID == "" {
		return nil, fmt.Errorf("用户组数据缺少 usrgrpid")
	}
	logger.L().Infof("禁用用户: %s, 加入用户组: %s", userId, targetGroupID)
	var userSpec models.UserParams
	userSpec.Userid = userId
	userSpec.Usrgrps = []string{targetGroupID}
	pwd, err := utils.GenerateSecurePassword(12) // 密码无需回传
	if err != nil {
		logger.L().Error("生成密码失败: %s", err.Error())
		return nil, err
	}
	// TODO 设置
	users, err := UpdateUser(ctx, provider, userSpec, instance, pwd)
	if err != nil {
		logger.L().Error("禁用用户失败: %s", err.Error())
		return nil, err
	}
	return users, nil
}
