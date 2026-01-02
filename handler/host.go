/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2025-12-18 11:20:36
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2025-12-31 16:44:42
 * @FilePath: \zabbix-mcp-go\handler\host.go
 * @Description: 主机相关功能
 * @Copyright: Copyright (c) 2025 by fengzhilaoling@gmail.com, All Rights Reserved.
 */
package handler

import (
	"context"
	"fmt"

	"zabbixMcp/logger"
	"zabbixMcp/models"
	"zabbixMcp/server"

	"github.com/mark3labs/mcp-go/mcp"
)

// GetHostsHandler 通过注入的 ClientProvider 调用 host.get 并返回结果
func GetHostsHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	instance := ""
	hostnames := []string{}
	activeAvailable := ""
	typezbx := "1"
	selectParams := map[string]bool{
		"select_discoveries":              false,
		"select_discovery_data":           false,
		"select_discovery_rule":           false,
		"select_discovery_rule_prototype": false,
		"select_graphs":                   false,
		"select_host_discovery":           false,
		"select_host_groups":              false,
		"select_http_tests":               false,
		"select_interfaces":               false,
		"select_inventory":                false,
		"select_items":                    false,
		"select_macros":                   false,
		"select_parent_templates":         false,
		"select_dashboards":               false,
		"select_tags":                     false,
		"select_inherited_tags":           false,
		"select_triggers":                 false,
		"select_value_maps":               false,
		"search":                          false,
	}

	if args, ok := req.Params.Arguments.(map[string]interface{}); ok {
		if v, ok2 := args["instance"].(string); ok2 {
			instance = v
		}
		if arr, ok := args["hostnames"].([]interface{}); ok {
			for _, v := range arr {
				if s, ok := v.(string); ok && s != "" {
					hostnames = append(hostnames, s)
				}
			}
		}
		if v, ok := args["active_available"].(string); ok {
			activeAvailable = v
		}
		if v, ok := args["type"].(string); ok {
			typezbx = v
		}
		for key := range selectParams {
			if val, ok := args[key].(bool); ok {
				selectParams[key] = val
			}
		}
	}
	if clientPool == nil {
		return mcp.NewToolResultStructuredOnly(makeResult([]map[string]interface{}{})), nil
	}
	spec := models.HostParams{Output: "extend", SelectInterfaces: "extend"}
	if activeAvailable != "" {
		spec.Filter = map[string]interface{}{"active_available": activeAvailable, "type": typezbx}
	}
	logger.L().Infof("instance: %s, hostname: %v", instance, hostnames)
	if len(hostnames) > 0 {
		if selectParams["search"] {
			spec.Search = map[string]interface{}{"host": hostnames}
		} else {
			spec.Filter = map[string]interface{}{"host": hostnames}
		}
		if selectParams["select_discoveries"] {
			spec.SelectDiscoveries = "extend"
		}
		if selectParams["select_discovery_data"] {
			spec.SelectDiscoveryData = "extend"
		}
		if selectParams["select_discovery_rule"] {
			spec.SelectDiscoveryRule = "extend"
		}
		if selectParams["select_discovery_rule_prototype"] {
			spec.SelectDiscoveryRulePrototype = "extend"
		}
		if selectParams["select_graphs"] {
			spec.SelectGraphs = "extend"
		}
		if selectParams["select_host_discovery"] {
			spec.SelectHostDiscovery = "extend"
		}
		if selectParams["select_host_groups"] {
			// selectHostGroups does not support the string value "extend".
			// Instead request the specific host group fields that are supported by the API.
			spec.SelectHostGroups = []string{
				"groupid",
				"name",
				"flags",
				"uuid",
			}
		}
		if selectParams["select_http_tests"] {
			spec.SelectHTTPTests = "extend"
		}
		if selectParams["select_interfaces"] {
			spec.SelectInterfaces = "extend"
		}
		if selectParams["select_inventory"] {
			spec.SelectInventory = "extend"
		}
		if selectParams["select_items"] {
			spec.SelectItems = "extend"
		}
		if selectParams["select_macros"] {
			spec.SelectMacros = "extend"
		}
		if selectParams["select_parent_templates"] {
			// selectParentTemplates does not support the string value "extend".
			// Instead request the specific template fields that are supported by the API.
			spec.SelectParentTemplates = []string{
				"templateid",
				"host",
				"description",
				"name",
				"uuid",
				"vendor_name",
				"vendor_version",
			}
		}
		if selectParams["select_dashboards"] {
			spec.SelectDashboards = "extend"
		}
		if selectParams["select_tags"] {
			spec.SelectTags = "extend"
		}
		if selectParams["select_inherited_tags"] {
			spec.SelectInheritedTags = "extend"
		}
		if selectParams["select_triggers"] {
			spec.SelectTriggers = "extend"
		}
		if selectParams["select_value_maps"] {
			spec.SelectValueMaps = "extend"
		}
	}
	hosts, err := server.GetHosts(ctx, clientPool, spec, instance)
	if err != nil {
		return nil, fmt.Errorf("调用 host.get 失败: %w", err)
	}
	return mcp.NewToolResultStructuredOnly(makeResult(hosts)), nil
}

// 通过主机组查询
// 通过主机名查询 详细信息 ()
