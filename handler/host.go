/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2025-12-18 11:20:36
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2026-01-24 15:33:50
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
	"zabbixMcp/utils"

	"github.com/mark3labs/mcp-go/mcp"
)

// 通过主机名查询主机信息
func GetHostForNameHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	instance := ""
	hostname := ""
	searchWildcardsEnabled := true
	var selectGraphs interface{}
	limit := 15
	limitSelects := 100
	count := false
	mcpToolName := req.Params.Name
	logger.L().Infof("GetHostForNameHandler: mcpToolName=%s", mcpToolName)
	if args, ok := req.Params.Arguments.(map[string]interface{}); ok {
		if v, ok2 := args["instance"].(string); ok2 {
			instance = v
		}
		if v, ok2 := args["hostname"].(string); ok2 {
			hostname = v
		}
		if v, ok2 := args["searchWildcardsEnabled"].(bool); ok2 {
			searchWildcardsEnabled = v
		}
		// selectGraphs 支持两种形式: string "extend"(默认) 或 []string
		if raw, ok2 := args["selectGraphs"]; ok2 {
			// 先尝试处理 []interface{} (常见于 JSON 解码)
			if arr, ok3 := raw.([]interface{}); ok3 {
				tmp := make([]string, 0, len(arr))
				for _, it := range arr {
					if s, ok4 := it.(string); ok4 {
						tmp = append(tmp, s)
					}
				}
				if len(tmp) == 0 {
					selectGraphs = "extend"
				} else {
					// 如果数组第一个元素为 "extend"，按照 Zabbix API 约定使用字符串 "extend"
					if tmp[0] == "extend" {
						selectGraphs = "extend"
					} else {
						selectGraphs = tmp
					}
				}
			} else if arrS, ok3 := raw.([]string); ok3 {
				if len(arrS) == 0 {
					selectGraphs = "extend"
				} else {
					if arrS[0] == "extend" {
						selectGraphs = "extend"
					} else {
						selectGraphs = arrS
					}
				}
			} else if s, ok3 := raw.(string); ok3 {
				// 允许直接传入字符串 "extend"
				selectGraphs = s
			} else {
				// 未知类型，回退为默认
				selectGraphs = "extend"
			}
		}
		if v, ok2 := args["limit"].(int); ok2 {
			limit = v
		}
		if v, ok2 := args["limitSelects"].(int); ok2 {
			limitSelects = v
		}
		if v, ok2 := args["count"].(bool); ok2 {
			count = v
		}
	}
	logger.L().Infof("GetHostForNameHandler: instance=%s, hostname=%s, searchWildcardsEnabled=%v, limit=%d", instance, hostname, searchWildcardsEnabled, limit)
	logger.L().Infof("GetHostForNameHandler: selectGraphs=%v, limitSelects=%d, count=%v", selectGraphs, limitSelects, count)
	spec := models.HostParams{
		Output:                 "extend",
		SelectInterfaces:       "extend",
		Search:                 map[string]interface{}{"host": hostname},
		Limit:                  limit,
		LimitSelects:           limitSelects,
		SearchWildcardsEnabled: searchWildcardsEnabled,
	}
	if mcpToolName == "get_graph_by_hostname" {
		if count {
			spec.SelectGraphs = "count"
		} else {
			spec.SelectGraphs = selectGraphs
		}
	}
	hosts, err := server.GetHosts(ctx, clientPool, spec, instance)
	if err != nil {
		return nil, fmt.Errorf("调用 host.get 失败: %w", err)
	}
	return mcp.NewToolResultStructuredOnly(makeResult(hosts)), nil
}

// UpdateHostHandler 通过注入的 ClientProvider 调用 host.update 并返回结果
// func UpdateHostHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
// 	instanceName := ""
// 	hostids := []string{}
// 	host := ""
// 	name := ""
// 	groupsReplace := make([]map[string]interface{}, 0)
// 	interfacesReplace := make([]map[string]interface{}, 0)
// 	templatesReplace := make([]map[string]interface{}, 0)
// 	templatesClear := make([]map[string]interface{}, 0)
// 	tags := make([]models.Tag, 0)
// 	macrosReplace := make([]map[string]interface{}, 0)
// 	var inventory map[string]interface{}

// 	if args, ok := req.Params.Arguments.(map[string]interface{}); ok {
// 		if v, ok2 := args["instance"].(string); ok2 {
// 			instanceName = v
// 		}
// 		if arr, ok2 := args["hostids"].([]interface{}); ok2 {
// 			for _, it := range arr {
// 				if s, ok3 := it.(string); ok3 && s != "" {
// 					hostids = append(hostids, s)
// 				}
// 			}
// 		}
// 		if v, ok2 := args["hostid"].(string); ok2 && v != "" {
// 			hostids = append(hostids, v)
// 		}
// 		if v, ok2 := args["host"].(string); ok2 {
// 			host = v
// 		}
// 		if v, ok2 := args["name"].(string); ok2 {
// 			name = v
// 		}

// 		// groups (用于替换当前主机组)。接受 groups 或 groups_replace 两种形式，优先使用 groups
// 		if arr, ok2 := args["groups"].([]interface{}); ok2 {
// 			for _, it := range arr {
// 				if m, ok3 := it.(map[string]interface{}); ok3 {
// 					if gid, ok4 := m["groupid"]; ok4 {
// 						groupsReplace = append(groupsReplace, map[string]interface{}{"groupid": gid})
// 					}
// 				}
// 			}
// 		}
// 		if arr, ok2 := args["groups_replace"].([]interface{}); ok2 {
// 			for _, it := range arr {
// 				if m, ok3 := it.(map[string]interface{}); ok3 {
// 					if gid, ok4 := m["groupid"]; ok4 {
// 						groupsReplace = append(groupsReplace, map[string]interface{}{"groupid": gid})
// 					}
// 				}
// 			}
// 		}

// 		// interfaces 用于替换当前主机接口（接受 interfaces 或 interfaces_replace）
// 		if arr, ok2 := args["interfaces"].([]interface{}); ok2 {
// 			for _, it := range arr {
// 				if m, ok3 := it.(map[string]interface{}); ok3 {
// 					interfacesReplace = append(interfacesReplace, m)
// 				}
// 			}
// 		}
// 		if arr, ok2 := args["interfaces_replace"].([]interface{}); ok2 {
// 			for _, it := range arr {
// 				if m, ok3 := it.(map[string]interface{}); ok3 {
// 					interfacesReplace = append(interfacesReplace, m)
// 				}
// 			}
// 		}

// 		// templates: 在 update 场景中 templates 用于替换关联模板（templates_clear 用于取消关联并 clear）
// 		if arr, ok2 := args["templates"].([]interface{}); ok2 {
// 			for _, it := range arr {
// 				if m, ok3 := it.(map[string]interface{}); ok3 {
// 					if tid, ok4 := m["templateid"]; ok4 {
// 						templatesReplace = append(templatesReplace, map[string]interface{}{"templateid": tid})
// 					}
// 				}
// 			}
// 		}
// 		if arr, ok2 := args["templateids"].([]interface{}); ok2 {
// 			for _, it := range arr {
// 				if s, ok3 := it.(string); ok3 && s != "" {
// 					templatesReplace = append(templatesReplace, map[string]interface{}{"templateid": s})
// 				}
// 			}
// 		}
// 		if arr, ok2 := args["templates_replace"].([]interface{}); ok2 {
// 			for _, it := range arr {
// 				if m, ok3 := it.(map[string]interface{}); ok3 {
// 					if tid, ok4 := m["templateid"]; ok4 {
// 						templatesReplace = append(templatesReplace, map[string]interface{}{"templateid": tid})
// 					}
// 				}
// 			}
// 		}
// 		if arr, ok2 := args["templates_clear"].([]interface{}); ok2 {
// 			for _, it := range arr {
// 				if m, ok3 := it.(map[string]interface{}); ok3 {
// 					if tid, ok4 := m["templateid"]; ok4 {
// 						templatesClear = append(templatesClear, map[string]interface{}{"templateid": tid})
// 					}
// 				}
// 			}
// 		}

// 		// tags 用于替换当前主机标签（tags 或 tags_replace）
// 		if raw, ok := args["tags"]; ok {
// 			if arr, ok2 := raw.([]interface{}); ok2 {
// 				for _, it := range arr {
// 					if m, ok3 := it.(map[string]interface{}); ok3 {
// 						tag, _ := m["tag"].(string)
// 						value, _ := m["value"].(string)
// 						tags = append(tags, models.Tag{Tag: tag, Value: value})
// 					}
// 				}
// 			}
// 		}

// 		// macros 用于替换当前用户宏（macros 或 macros_replace）
// 		if arr, ok2 := args["macros"].([]interface{}); ok2 {
// 			for _, it := range arr {
// 				if m, ok3 := it.(map[string]interface{}); ok3 {
// 					macrosReplace = append(macrosReplace, m)
// 				}
// 			}
// 		}
// 		if arr, ok2 := args["macros_replace"].([]interface{}); ok2 {
// 			for _, it := range arr {
// 				if m, ok3 := it.(map[string]interface{}); ok3 {
// 					macrosReplace = append(macrosReplace, m)
// 				}
// 			}
// 		}

// 		if inv, ok2 := args["inventory"].(map[string]interface{}); ok2 {
// 			inventory = inv
// 		}
// 	}

// 	if clientPool == nil {
// 		return mcp.NewToolResultStructuredOnly(makeResult([]map[string]interface{}{})), nil
// 	}

// 	// Validate hostids: require exactly one hostid for host.update
// 	if len(hostids) == 0 {
// 		return nil, fmt.Errorf("host.update 需要提供一个 hostid 或 hostids")
// 	}
// 	if len(hostids) > 1 {
// 		return nil, fmt.Errorf("host.update 不支持一次更新多个主机，请对每个主机分别调用 update_host")
// 	}

// 	// Build spec using replace semantics for update
// 	spec := models.HostParams{
// 		HostID:            hostids[0],
// 		Host:              host,
// 		Name:              name,
// 		GroupsReplace:     groupsReplace,
// 		InterfacesReplace: interfacesReplace,
// 		TemplatesReplace:  templatesReplace,
// 		TemplatesClear:    templatesClear,
// 		Tags:              tags,
// 		MacrosReplace:     macrosReplace,
// 		Inventory:         inventory,
// 	}

// 	result, err := server.UpdateHost(ctx, clientPool, spec, instanceName)
// 	if err != nil {
// 		logger.L().Errorf("调用 host.update 失败: %v", err)
// 		return nil, fmt.Errorf("调用 host.update 失败: %w", err)
// 	}
// 	return mcp.NewToolResultStructuredOnly(makeResult(result)), nil
// }

// // 主机分类查询
// func GetHostsHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
// 	instance := ""
// 	hostids := []string{}
// 	style := ""
// 	is_detailed := false
// 	if args, ok := req.Params.Arguments.(map[string]interface{}); ok {
// 		if v, ok2 := args["instance"].(string); ok2 {
// 			instance = v
// 		}
// 		if arr, ok := args["hostids"].([]interface{}); ok {
// 			for _, v := range arr {
// 				if s, ok := v.(string); ok && s != "" {
// 					hostids = append(hostids, s)
// 				}
// 			}
// 		}
// 		if v, ok := args["style"].(string); ok {
// 			style = v
// 		}
// 		if v, ok := args["is_detailed"].(bool); ok {
// 			is_detailed = v
// 		}
// 	}
// 	if clientPool == nil {
// 		return mcp.NewToolResultStructuredOnly(makeResult([]map[string]interface{}{})), nil
// 	}
// 	spec := models.HostParams{Output: "extend", SelectInterfaces: "extend"}
// 	spec.HostIDs = hostids
// 	switch style {
// 	case "graphs":
// 		if is_detailed {
// 			spec.SelectGraphs = "extend"
// 		} else {
// 			spec.SelectGraphs = []string{"graphid", "name"}
// 		}
// 	case "templates":
// 		if is_detailed {
// 			spec.SelectParentTemplates = "extend"
// 		} else {
// 			spec.SelectParentTemplates = []string{"templateid", "name"}
// 		}
// 	case "items":
// 		if is_detailed {
// 			spec.SelectItems = "extend"
// 		} else {
// 			spec.SelectItems = []string{"itemid", "name"}
// 		}
// 	case "triggers":
// 		if is_detailed {
// 			spec.SelectTriggers = "extend"
// 		} else {
// 			spec.SelectTriggers = []string{"triggerid", "name"}
// 		}
// 	case "macros":
// 		if is_detailed {
// 			spec.SelectMacros = "extend"
// 		} else {
// 			spec.SelectMacros = []string{"macro", "value"}
// 		}
// 	case "value_maps":
// 		if is_detailed {
// 			spec.SelectValueMaps = "extend"
// 		} else {
// 			spec.SelectValueMaps = []string{"valuemapid", "name"}
// 		}
// 	case "groups":
// 		if is_detailed {
// 			spec.SelectHostGroups = "extend"
// 		} else {
// 			spec.SelectHostGroups = []string{"groupid", "name", "flags", "uuid"}
// 		}
// 	case "tags":
// 		if is_detailed {
// 			spec.SelectTags = "extend"
// 		} else {
// 			spec.SelectTags = []string{"tag", "value"}
// 		}
// 	case "dashboards":
// 		if is_detailed {
// 			spec.SelectDashboards = "extend"
// 		} else {
// 			spec.SelectDashboards = []string{"dashboardid", "name"}
// 		}
// 	case "discoveries":
// 		if is_detailed {
// 			spec.SelectDiscoveries = "extend"
// 		} else {
// 			spec.SelectDiscoveries = []string{"name", "delay", "key_", "status"}
// 		}
// 	case "discovery_rule":
// 		if is_detailed {
// 			spec.SelectDiscoveryRule = "extend"
// 		} else {
// 			spec.SelectDiscoveryRule = []string{"itemid", "name", "delay"}
// 		}
// 	case "discovery_rule_prototype":
// 		if is_detailed {
// 			spec.SelectDiscoveryRulePrototype = "extend"
// 		} else {
// 			spec.SelectDiscoveryRulePrototype = []string{"itemid", "name", "delay"}
// 		}
// 	case "discovery_data":
// 		if is_detailed {
// 			spec.SelectDiscoveryData = "extend"
// 		} else {
// 			spec.SelectDiscoveryData = []string{"itemid", "name", "delay"}
// 		}
// 	case "discovery":
// 		if is_detailed {
// 			spec.SelectHostDiscovery = "extend"
// 		} else {
// 			spec.SelectHostDiscovery = []string{"discoveryid", "name", "delay"}
// 		}
// 	case "http_tests":
// 		if is_detailed {
// 			spec.SelectHTTPTests = "extend"
// 		} else {
// 			spec.SelectHTTPTests = []string{"httptestid", "name", "agent", "delay"}
// 		}
// 	case "interfaces":
// 		if is_detailed {
// 			spec.SelectInterfaces = "extend"
// 		} else {
// 			spec.SelectInterfaces = []string{"interfaceid", "available", "type", "ip", "dns", "useip"}
// 		}
// 	case "inventory":
// 		if is_detailed {
// 			spec.SelectInventory = "extend"
// 		} else {
// 			spec.SelectInventory = []string{"hostid", "name", "description"}
// 		}
// 	case "inherited_tags":
// 		spec.SelectInheritedTags = "extend"
// 	default:
// 		return nil, fmt.Errorf("不支持的 style: %s", style)
// 	}
// 	hosts, err := server.GetHosts(ctx, clientPool, spec, instance)
// 	if err != nil {
// 		return nil, fmt.Errorf("调用 host.get 失败: %w", err)
// 	}
// 	return mcp.NewToolResultStructuredOnly(makeResult(hosts)), nil
// }

// CreateHostHandler 创建主机
func CreateHostHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	instanceName := ""
	host := ""
	name := ""
	groups := []string{}
	templates := []models.Templates{}
	tags := make([]models.Tag, 0)
	macros := []models.Macros{}
	// inventory := make(map[string]interface{})
	interfaces := []models.ZabbixInterface{}
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
		if v, ok2 := args["groups"].([]string); ok2 {
			groups = v
		}
		if v, ok2 := args["templates"].([]models.Templates); ok2 {
			templates = v
		}
		if v, ok := args["tags"].([]models.Tag); ok {
			tags = v
		}
		if v, ok2 := args["macros"].([]models.Macros); ok2 {
			macros = v
		}
		// if v, ok2 := args["inventory"].(map[string]interface{}); ok2 {
		// 	inventory = v
		// }
		if v, ok2 := args["interfaces"].([]models.ZabbixInterface); ok2 {
			interfaces = v
		}
	}
	logger.L().Infof("groups: %s, templates: %s", groups, templates)
	if clientPool == nil {
		return mcp.NewToolResultStructuredOnly(makeResult([]map[string]interface{}{})), nil
	}

	spec := models.HostParams{
		Host:       host,
		Name:       name,
		Interfaces: interfaces,
		Tags:       tags,
		Macros:     macros,
		// Inventory:  inventory,
	}
	result, err := server.CreateHost(ctx, clientPool, spec, instanceName)
	if err != nil {
		logger.L().Errorf("调用 host.create 失败: %v", err)
		return nil, fmt.Errorf("调用 host.create 失败: %w", err)
	}
	return mcp.NewToolResultStructuredOnly(makeResult(result)), nil
}

func UpdateNewHostHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	instanceName := ""
	// hostid := ""
	groups := []map[string]interface{}{}
	templates := []map[string]interface{}{}
	tags := []models.Tag{}
	interfaces := []models.ZabbixInterface{}
	// macros := []map[string]interface{}{}
	// inventory := map[string]interface{}{}
	style := ""
	if args, ok := req.Params.Arguments.(map[string]interface{}); ok {
		if v, ok2 := args["instance"].(string); ok2 {
			instanceName = v
		}
		// if v, ok := args["hostid"].(string); ok && v != "" {
		// 	hostid = v
		// }
		if v, ok := args["style"].(string); ok && v != "" {
			style = v
		}
		if arr, ok2 := args["groups"].([]interface{}); ok2 {
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
		if arr, ok2 := args["templates"].([]interface{}); ok2 {
			for _, it := range arr {
				if s, ok3 := it.(string); ok3 && s != "" {
					templates = append(templates, map[string]interface{}{"templateid": s})
				}
			}
		}
		if arr, ok2 := args["templates"].([]interface{}); ok2 {
			for _, it := range arr {
				if m, ok3 := it.(map[string]interface{}); ok3 {
					// only keep templateid if present
					if tid, ok4 := m["templateid"]; ok4 {
						templates = append(templates, map[string]interface{}{"templateid": tid})
					}
				}
			}
		}
		// tags 兼容AI会传输
		if raw, ok := args["tags"]; ok {
			if arr, ok2 := raw.([]interface{}); ok2 {
				for _, it := range arr {
					if m, ok3 := it.(map[string]interface{}); ok3 {
						tag, _ := m["tag"].(string)
						value, _ := m["value"].(string)
						tags = append(tags, models.Tag{Tag: tag, Value: value})
					}
				}
			}
		}
		// interfaces 兼容AI会传输
		if raw, ok := args["interfaces"]; ok {
			if arr, ok2 := raw.([]interface{}); ok2 {
				for _, it := range arr {
					if m, ok3 := it.(map[string]interface{}); ok3 {
						interfaceid := m["interfaceid"].(string)
						interfaceType := utils.JsonInt(m["type"], 1)
						hostid := m["hostid"].(string)
						main := utils.JsonInt(m["main"], 1)
						useip := utils.JsonInt(m["useip"], 1)
						ip, _ := m["ip"].(string)
						dns, _ := m["dns"].(string)
						port, _ := m["port"].(string)
						interfaces = append(interfaces, models.ZabbixInterface{
							InterfaceID: interfaceid,
							HostID:      hostid,
							Type:        interfaceType,
							Main:        main,
							UseIP:       useip,
							IP:          ip,
							DNS:         dns,
							Port:        port,
						})
					}
				}
			}
		}
		// if v, ok2 := args["macros"].([]interface{}); ok2 {
		// 	macros = v
		// }
		// if v, ok2 := args["inventory"].([]interface{}); ok2 {
		// 	inventory = v
		// }
	}
	logger.L().Infof("groups %v, templates %v, interfaces %v", groups, templates, interfaces)
	// Build spec using replace semantics for update
	spec := models.HostParams{}
	// spec.HostID = hostid
	switch style {
	// case "groups":
	// 	spec.GroupsReplace = groups
	// case "templates":
	// 	spec.TemplatesReplace = templates
	// case "templates_clear":
	// 	spec.TemplatesClear = templates
	// case "tags":
	// 	spec.Tags = tags
	// case "interfaces":
	// 	spec.ZabbixInterfaces = interfaces
	default:
		return nil, fmt.Errorf("style %s 不支持", style)
	}
	result, err := server.UpdateHost(ctx, clientPool, spec, instanceName)
	if err != nil {
		logger.L().Errorf("调用 host.update 失败: %v", err)
		return nil, fmt.Errorf("调用 host.update 失败: %w", err)
	}
	return mcp.NewToolResultStructuredOnly(makeResult(result)), nil
}

// DeleteHostsHandler 删除主机
func DeleteHostsHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	instanceName := ""
	hostIDs := []string{}
	if args, ok := req.Params.Arguments.(map[string]interface{}); ok {
		if v, ok2 := args["instance"].(string); ok2 {
			instanceName = v
		}
		if arr, ok := args["hostids"].([]string); ok {
			hostIDs = arr

		}
	}
	logger.L().Infof("hostids %v", hostIDs)
	if clientPool == nil {
		return mcp.NewToolResultStructuredOnly(makeResult([]map[string]interface{}{})), nil
	}
	spec := models.HostParams{HostIds: hostIDs}
	hosts, err := server.DeleteHosts(ctx, clientPool, spec, instanceName)
	if err != nil {
		logger.L().Errorf("调用 host.delete 失败: %w", err)
		return nil, fmt.Errorf("调用 host.delete 失败: %w", err)
	}
	return mcp.NewToolResultStructuredOnly(makeResult(hosts)), nil
}
