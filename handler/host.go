/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2025-12-18 11:20:36
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2026-01-08 13:59:51
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

// CreateHostHandler 通过注入的 ClientProvider 调用 host.create 并返回结果
func CreateHostHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	instanceName := ""
	host := ""
	name := ""
	groups := make([]map[string]interface{}, 0)
	interfaces := make([]map[string]interface{}, 0)
	templates := make([]map[string]interface{}, 0)
	tags := make([]map[string]interface{}, 0)
	macros := make([]map[string]interface{}, 0)
	var inventory map[string]interface{}

	if args, ok := req.Params.Arguments.(map[string]interface{}); ok {
		if v, ok2 := args["instance"].(string); ok2 {
			instanceName = v
		}
		if v, ok2 := args["host"].(string); ok2 {
			host = v
		}
		if v, ok2 := args["name"].(string); ok2 {
			name = v
		}
		// groups can be provided as groupids []string or groups []map
		if arr, ok2 := args["groupids"].([]interface{}); ok2 {
			for _, it := range arr {
				if s, ok3 := it.(string); ok3 && s != "" {
					groups = append(groups, map[string]interface{}{"groupid": s})
				}
			}
		}
		if arr, ok2 := args["groups"].([]interface{}); ok2 {
			for _, it := range arr {
				if m, ok3 := it.(map[string]interface{}); ok3 {
					// only keep groupid if present
					if gid, ok4 := m["groupid"]; ok4 {
						groups = append(groups, map[string]interface{}{"groupid": gid})
					}
				}
			}
		}
		// interfaces array
		if arr, ok2 := args["interfaces"].([]interface{}); ok2 {
			for _, it := range arr {
				if m, ok3 := it.(map[string]interface{}); ok3 {
					interfaces = append(interfaces, m)
				}
			}
		}
		// templates can be templateids []string or templates []map
		if arr, ok2 := args["templateids"].([]interface{}); ok2 {
			for _, it := range arr {
				if s, ok3 := it.(string); ok3 && s != "" {
					templates = append(templates, map[string]interface{}{"templateid": s})
				}
			}
		}
		if arr, ok2 := args["templates"].([]interface{}); ok2 {
			for _, it := range arr {
				if m, ok3 := it.(map[string]interface{}); ok3 {
					if tid, ok4 := m["templateid"]; ok4 {
						templates = append(templates, map[string]interface{}{"templateid": tid})
					}
				}
			}
		}
		// tags
		if arr, ok2 := args["tags"].([]interface{}); ok2 {
			for _, it := range arr {
				if m, ok3 := it.(map[string]interface{}); ok3 {
					tags = append(tags, m)
				}
			}
		}
		// macros
		if arr, ok2 := args["macros"].([]interface{}); ok2 {
			for _, it := range arr {
				if m, ok3 := it.(map[string]interface{}); ok3 {
					macros = append(macros, m)
				}
			}
		}
		if inv, ok2 := args["inventory"].(map[string]interface{}); ok2 {
			inventory = inv
		}
	}

	if clientPool == nil {
		return mcp.NewToolResultStructuredOnly(makeResult([]map[string]interface{}{})), nil
	}

	spec := models.HostParams{
		Host:            host,
		Name:            name,
		Groups:          groups,
		Interfaces:      interfaces,
		TemplatesToLink: templates,
		TagsToCreate:    tags,
		MacrosToCreate:  macros,
		Inventory:       inventory,
	}

	result, err := server.CreateHost(ctx, clientPool, spec, instanceName)
	if err != nil {
		logger.L().Errorf("调用 host.create 失败: %v", err)
		return nil, fmt.Errorf("调用 host.create 失败: %w", err)
	}
	return mcp.NewToolResultStructuredOnly(makeResult(result)), nil
}

// UpdateHostHandler 通过注入的 ClientProvider 调用 host.update 并返回结果
func UpdateHostHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	instanceName := ""
	hostids := []string{}
	host := ""
	name := ""
	groupsReplace := make([]map[string]interface{}, 0)
	interfacesReplace := make([]map[string]interface{}, 0)
	templatesReplace := make([]map[string]interface{}, 0)
	templatesClear := make([]map[string]interface{}, 0)
	tagsReplace := make([]map[string]interface{}, 0)
	macrosReplace := make([]map[string]interface{}, 0)
	var inventory map[string]interface{}

	if args, ok := req.Params.Arguments.(map[string]interface{}); ok {
		if v, ok2 := args["instance"].(string); ok2 {
			instanceName = v
		}
		if arr, ok2 := args["hostids"].([]interface{}); ok2 {
			for _, it := range arr {
				if s, ok3 := it.(string); ok3 && s != "" {
					hostids = append(hostids, s)
				}
			}
		}
		if v, ok2 := args["hostid"].(string); ok2 && v != "" {
			hostids = append(hostids, v)
		}
		if v, ok2 := args["host"].(string); ok2 {
			host = v
		}
		if v, ok2 := args["name"].(string); ok2 {
			name = v
		}

		// groups (用于替换当前主机组)。接受 groups 或 groups_replace 两种形式，优先使用 groups
		if arr, ok2 := args["groups"].([]interface{}); ok2 {
			for _, it := range arr {
				if m, ok3 := it.(map[string]interface{}); ok3 {
					if gid, ok4 := m["groupid"]; ok4 {
						groupsReplace = append(groupsReplace, map[string]interface{}{"groupid": gid})
					}
				}
			}
		}
		if arr, ok2 := args["groups_replace"].([]interface{}); ok2 {
			for _, it := range arr {
				if m, ok3 := it.(map[string]interface{}); ok3 {
					if gid, ok4 := m["groupid"]; ok4 {
						groupsReplace = append(groupsReplace, map[string]interface{}{"groupid": gid})
					}
				}
			}
		}

		// interfaces 用于替换当前主机接口（接受 interfaces 或 interfaces_replace）
		if arr, ok2 := args["interfaces"].([]interface{}); ok2 {
			for _, it := range arr {
				if m, ok3 := it.(map[string]interface{}); ok3 {
					interfacesReplace = append(interfacesReplace, m)
				}
			}
		}
		if arr, ok2 := args["interfaces_replace"].([]interface{}); ok2 {
			for _, it := range arr {
				if m, ok3 := it.(map[string]interface{}); ok3 {
					interfacesReplace = append(interfacesReplace, m)
				}
			}
		}

		// templates: 在 update 场景中 templates 用于替换关联模板（templates_clear 用于取消关联并 clear）
		if arr, ok2 := args["templates"].([]interface{}); ok2 {
			for _, it := range arr {
				if m, ok3 := it.(map[string]interface{}); ok3 {
					if tid, ok4 := m["templateid"]; ok4 {
						templatesReplace = append(templatesReplace, map[string]interface{}{"templateid": tid})
					}
				}
			}
		}
		if arr, ok2 := args["templateids"].([]interface{}); ok2 {
			for _, it := range arr {
				if s, ok3 := it.(string); ok3 && s != "" {
					templatesReplace = append(templatesReplace, map[string]interface{}{"templateid": s})
				}
			}
		}
		if arr, ok2 := args["templates_replace"].([]interface{}); ok2 {
			for _, it := range arr {
				if m, ok3 := it.(map[string]interface{}); ok3 {
					if tid, ok4 := m["templateid"]; ok4 {
						templatesReplace = append(templatesReplace, map[string]interface{}{"templateid": tid})
					}
				}
			}
		}
		if arr, ok2 := args["templates_clear"].([]interface{}); ok2 {
			for _, it := range arr {
				if m, ok3 := it.(map[string]interface{}); ok3 {
					if tid, ok4 := m["templateid"]; ok4 {
						templatesClear = append(templatesClear, map[string]interface{}{"templateid": tid})
					}
				}
			}
		}

		// tags 用于替换当前主机标签（tags 或 tags_replace）
		if arr, ok2 := args["tags"].([]interface{}); ok2 {
			for _, it := range arr {
				if m, ok3 := it.(map[string]interface{}); ok3 {
					tagsReplace = append(tagsReplace, m)
				}
			}
		}
		if arr, ok2 := args["tags_replace"].([]interface{}); ok2 {
			for _, it := range arr {
				if m, ok3 := it.(map[string]interface{}); ok3 {
					tagsReplace = append(tagsReplace, m)
				}
			}
		}

		// macros 用于替换当前用户宏（macros 或 macros_replace）
		if arr, ok2 := args["macros"].([]interface{}); ok2 {
			for _, it := range arr {
				if m, ok3 := it.(map[string]interface{}); ok3 {
					macrosReplace = append(macrosReplace, m)
				}
			}
		}
		if arr, ok2 := args["macros_replace"].([]interface{}); ok2 {
			for _, it := range arr {
				if m, ok3 := it.(map[string]interface{}); ok3 {
					macrosReplace = append(macrosReplace, m)
				}
			}
		}

		if inv, ok2 := args["inventory"].(map[string]interface{}); ok2 {
			inventory = inv
		}
	}

	if clientPool == nil {
		return mcp.NewToolResultStructuredOnly(makeResult([]map[string]interface{}{})), nil
	}

	// Validate hostids: require exactly one hostid for host.update
	if len(hostids) == 0 {
		return nil, fmt.Errorf("host.update 需要提供一个 hostid 或 hostids")
	}
	if len(hostids) > 1 {
		return nil, fmt.Errorf("host.update 不支持一次更新多个主机，请对每个主机分别调用 update_host")
	}

	// Build spec using replace semantics for update
	spec := models.HostParams{
		HostID:            hostids[0],
		Host:              host,
		Name:              name,
		GroupsReplace:     groupsReplace,
		InterfacesReplace: interfacesReplace,
		TemplatesReplace:  templatesReplace,
		TemplatesClear:    templatesClear,
		TagsReplace:       tagsReplace,
		MacrosReplace:     macrosReplace,
		Inventory:         inventory,
	}

	result, err := server.UpdateHost(ctx, clientPool, spec, instanceName)
	if err != nil {
		logger.L().Errorf("调用 host.update 失败: %v", err)
		return nil, fmt.Errorf("调用 host.update 失败: %w", err)
	}
	return mcp.NewToolResultStructuredOnly(makeResult(result)), nil
}

func GetHostsRefineHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	instance := ""
	hostids := []string{}
	style := ""
	is_detailed := false
	if args, ok := req.Params.Arguments.(map[string]interface{}); ok {
		if v, ok2 := args["instance"].(string); ok2 {
			instance = v
		}
		if arr, ok := args["hostids"].([]interface{}); ok {
			for _, v := range arr {
				if s, ok := v.(string); ok && s != "" {
					hostids = append(hostids, s)
				}
			}
		}
		if v, ok := args["style"].(string); ok {
			style = v
		}
		if v, ok := args["is_detailed"].(bool); ok {
			is_detailed = v
		}
	}
	if clientPool == nil {
		return mcp.NewToolResultStructuredOnly(makeResult([]map[string]interface{}{})), nil
	}
	spec := models.HostParams{Output: "extend", SelectInterfaces: "extend"}
	spec.HostIDs = hostids
	switch style {
	case "graphs":
		if is_detailed {
			spec.SelectGraphs = "extend"
		} else {
			spec.SelectGraphs = []string{"graphid", "name"}
		}

	default:
		return nil, fmt.Errorf("不支持的 style: %s", style)
	}
	hosts, err := server.GetHosts(ctx, clientPool, spec, instance)
	if err != nil {
		return nil, fmt.Errorf("调用 host.get 失败: %w", err)
	}
	return mcp.NewToolResultStructuredOnly(makeResult(hosts)), nil
}
