/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2025-12-18 11:13:06
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2025-12-18 11:13:07
 * @FilePath: \zabbix-mcp-go\server\host.go
 * @Description: 主机相关功能
 * @Copyright: Copyright (c) 2025 by fengzhilaoling@gmail.com, All Rights Reserved.
 */

package server

import (
	"fmt"
	"zabbixMcp/zabbix"
)

// GetHosts 调用底层 ZabbixClientHandler 执行 host.get，并返回解析后的列表
func GetHosts(client zabbix.ZabbixClientHandler, params map[string]interface{}) ([]map[string]interface{}, error) {
	if client == nil {
		return nil, fmt.Errorf("no zabbix client")
	}

	res, err := client.Call("host.get", params)
	if err != nil {
		return nil, err
	}

	hosts := make([]map[string]interface{}, 0)
	switch v := res.(type) {
	case []interface{}:
		for _, item := range v {
			if m, ok := item.(map[string]interface{}); ok {
				hosts = append(hosts, m)
			}
		}
	case []map[string]interface{}:
		hosts = v
	default:
		if m, ok := v.(map[string]interface{}); ok {
			hosts = append(hosts, m)
		}
	}

	return hosts, nil
}
