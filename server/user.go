/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2025-12-18 11:10:11
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2025-12-18 11:19:19
 * @FilePath: \zabbix-mcp-go\server\user.go
 * @Description: 用户相关功能
 * @Copyright: Copyright (c) 2025 by fengzhilaoling@gmail.com, All Rights Reserved.
 */

package server

import (
	"fmt"

	"zabbixMcp/zabbix"
)

// GetUsers 调用底层 ZabbixClientHandler 执行 user.get，并返回解析后的列表
func GetUsers(client zabbix.ZabbixClientHandler, params map[string]interface{}) ([]map[string]interface{}, error) {
	if client == nil {
		return nil, fmt.Errorf("no zabbix client")
	}

	res, err := client.Call("user.get", params)
	if err != nil {
		return nil, err
	}

	users := make([]map[string]interface{}, 0)
	switch v := res.(type) {
	case []interface{}:
		for _, item := range v {
			if m, ok := item.(map[string]interface{}); ok {
				users = append(users, m)
			}
		}
	case []map[string]interface{}:
		users = v
	default:
		if m, ok := v.(map[string]interface{}); ok {
			users = append(users, m)
		}
	}

	return users, nil
}
