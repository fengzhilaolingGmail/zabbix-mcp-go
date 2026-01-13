/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2025-12-18 11:20:36
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2026-01-13 09:36:18
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

// GetHostForNameHandler 通过主机名查询主机信息
func GetHostForNameHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	instance := ""
	hostnames := []string{}
	activeAvailable := ""
	typezbx := "1"
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
		if s, ok := args["hostnames"].(string); ok && s != "" {
			hostnames = append(hostnames, s)
		}
		if v, ok := args["active_available"].(string); ok {
			activeAvailable = v
		}
		if v, ok := args["type"].(string); ok {
			typezbx = v
		}
	}
	if clientPool == nil {
		return mcp.NewToolResultStructuredOnly(makeResult([]map[string]interface{}{})), nil
	}
	spec := models.HostParams{Output: "extend", SelectInterfaces: "extend"}
	if len(hostnames) > 0 {
		spec.Search = map[string]interface{}{"host": hostnames}
	}
	if activeAvailable != "" {
		spec.Filter = map[string]interface{}{"active_available": activeAvailable, "type": typezbx}
	}
	hosts, err := server.GetHosts(ctx, clientPool, spec, instance)
	if err != nil {
		return nil, fmt.Errorf("调用 host.get 失败: %w", err)
	}
	return mcp.NewToolResultStructuredOnly(makeResult(hosts)), nil
}

// 主机查询接口
func GetHostsHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
	case "templates":
		if is_detailed {
			spec.SelectParentTemplates = "extend"
		} else {
			spec.SelectParentTemplates = []string{"templateid", "name"}
		}
	case "items":
		if is_detailed {
			spec.SelectItems = "extend"
		} else {
			spec.SelectItems = []string{"itemid", "name"}
		}
	case "triggers":
		if is_detailed {
			spec.SelectTriggers = "extend"
		} else {
			spec.SelectTriggers = []string{"triggerid", "name"}
		}
	case "macros":
		if is_detailed {
			spec.SelectMacros = "extend"
		} else {
			spec.SelectMacros = []string{"macro", "value"}
		}
	case "value_maps":
		if is_detailed {
			spec.SelectValueMaps = "extend"
		} else {
			spec.SelectValueMaps = []string{"valuemapid", "name"}
		}
	case "groups":
		if is_detailed {
			spec.SelectHostGroups = "extend"
		} else {
			spec.SelectHostGroups = []string{"groupid", "name", "flags", "uuid"}
		}
	case "tags":
		if is_detailed {
			spec.SelectTags = "extend"
		} else {
			spec.SelectTags = []string{"tag", "value"}
		}
	case "dashboards":
		if is_detailed {
			spec.SelectDashboards = "extend"
		} else {
			spec.SelectDashboards = []string{"dashboardid", "name"}
		}
	case "discoveries":
		if is_detailed {
			spec.SelectDiscoveries = "extend"
		} else {
			spec.SelectDiscoveries = []string{"name", "delay", "key_", "status"}
		}
	case "discovery_rule":
		if is_detailed {
			spec.SelectDiscoveryRule = "extend"
		} else {
			spec.SelectDiscoveryRule = []string{"itemid", "name", "delay"}
		}
	case "discovery_rule_prototype":
		if is_detailed {
			spec.SelectDiscoveryRulePrototype = "extend"
		} else {
			spec.SelectDiscoveryRulePrototype = []string{"itemid", "name", "delay"}
		}
	case "discovery_data":
		if is_detailed {
			spec.SelectDiscoveryData = "extend"
		} else {
			spec.SelectDiscoveryData = []string{"itemid", "name", "delay"}
		}
	case "discovery":
		if is_detailed {
			spec.SelectHostDiscovery = "extend"
		} else {
			spec.SelectHostDiscovery = []string{"discoveryid", "name", "delay"}
		}
	case "http_tests":
		if is_detailed {
			spec.SelectHTTPTests = "extend"
		} else {
			spec.SelectHTTPTests = []string{"httptestid", "name", "agent", "delay"}
		}
	case "interfaces":
		if is_detailed {
			spec.SelectInterfaces = "extend"
		} else {
			spec.SelectInterfaces = []string{"interfaceid", "available", "type", "ip", "dns", "useip"}
		}
	case "inventory":
		if is_detailed {
			spec.SelectInventory = "extend"
		} else {
			spec.SelectInventory = []string{"hostid", "name", "description"}
		}
	case "inherited_tags":
		spec.SelectInheritedTags = "extend"
	default:
		return nil, fmt.Errorf("不支持的 style: %s", style)
	}
	hosts, err := server.GetHosts(ctx, clientPool, spec, instance)
	if err != nil {
		return nil, fmt.Errorf("调用 host.get 失败: %w", err)
	}
	return mcp.NewToolResultStructuredOnly(makeResult(hosts)), nil
}

func CreateHostHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	instanceName := ""
	host := ""
	name := ""
	hostType := ""
	main := ""
	useip := ""
	ip := ""
	dns := ""
	port := ""
	groups := []string{}
	templates := []string{}
	tags := make([]map[string]interface{}, 0)
	macros := make([]map[string]interface{}, 0)
	inventory := make(map[string]interface{})
	groupMap := make([]map[string]interface{}, 0)
	templateMap := make([]map[string]interface{}, 0)
	interfaces := make([]map[string]interface{}, 0)
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
		if v, ok2 := args["type"].(string); ok2 {
			hostType = v
		}
		if v, ok2 := args["main"].(string); ok2 {
			main = v
		}
		if v, ok2 := args["useip"].(string); ok2 {
			useip = v
		}
		if v, ok2 := args["ip"].(string); ok2 {
			ip = v
		}
		if v, ok2 := args["dns"].(string); ok2 {
			dns = v
		}
		if v, ok2 := args["port"].(string); ok2 {
			port = v
		}
		if arr, ok2 := args["groups"].([]interface{}); ok2 {
			for _, v := range arr {
				if s, ok := v.(string); ok && s != "" {
					groups = append(groups, s)
				}
			}
		}
		if arr, ok2 := args["groupids"].([]interface{}); ok2 {
			for _, v := range arr {
				if s, ok := v.(string); ok && s != "" {
					groups = append(groups, s)
				}
			}
		}
		// 兼容map格式
		if arr, ok2 := args["groupids"].([]interface{}); ok2 {
			for _, it := range arr {
				if s, ok3 := it.(string); ok3 && s != "" {
					groupMap = append(groupMap, map[string]interface{}{"groupid": s})
				}
			}
		}
		if arr, ok2 := args["groups"].([]interface{}); ok2 {
			for _, it := range arr {
				if m, ok3 := it.(map[string]interface{}); ok3 {
					// only keep groupid if present
					if gid, ok4 := m["groupid"]; ok4 {
						groupMap = append(groupMap, map[string]interface{}{"groupid": gid})
					}
				}
			}
		}
		if arr, ok2 := args["templates"].([]interface{}); ok2 {
			for _, v := range arr {
				if s, ok := v.(string); ok && s != "" {
					templates = append(templates, s)
				}
			}
		}
		if arr, ok2 := args["templateids"].([]interface{}); ok2 {
			for _, v := range arr {
				if s, ok := v.(string); ok && s != "" {
					templates = append(templates, s)
				}
			}
		}
		// 兼容map格式
		if arr, ok2 := args["templateids"].([]interface{}); ok2 {
			for _, it := range arr {
				if s, ok3 := it.(string); ok3 && s != "" {
					templateMap = append(templateMap, map[string]interface{}{"templateid": s})
				}
			}
		}
		if arr, ok2 := args["templates"].([]interface{}); ok2 {
			for _, it := range arr {
				if m, ok3 := it.(map[string]interface{}); ok3 {
					if tid, ok4 := m["templateid"]; ok4 {
						templateMap = append(templateMap, map[string]interface{}{"templateid": tid})
					}
				}
			}
		}
		if v, ok2 := args["tags"].([]map[string]interface{}); ok2 {
			tags = v
		}
		if v, ok2 := args["macros"].([]map[string]interface{}); ok2 {
			macros = v
		}
		if v, ok2 := args["inventory"].(map[string]interface{}); ok2 {
			inventory = v
		}
		// interfaces 兼容AI会传输
		if arr, ok2 := args["interfaces"].([]interface{}); ok2 {
			for _, it := range arr {
				if m, ok3 := it.(map[string]interface{}); ok3 {
					interfaces = append(interfaces, m)
				}
			}
		}
	}
	logger.L().Infof("groups: %s, templates: %s", groups, templates)
	if clientPool == nil {
		return mcp.NewToolResultStructuredOnly(makeResult([]map[string]interface{}{})), nil
	}

	if len(interfaces) == 0 {
		interfaces = append(interfaces, map[string]interface{}{
			"type":  hostType,
			"main":  main,
			"useip": useip,
			"ip":    ip,
			"dns":   dns,
			"port":  port,
		})
	}

	spec := models.HostParams{
		Host:           host,
		Name:           name,
		Interfaces:     interfaces,
		TagsToCreate:   tags,
		MacrosToCreate: macros,
		Inventory:      inventory,
	}

	for _, group := range groups {
		logger.L().Infof("group: %s", group)
		spec.Groups = append(spec.Groups, map[string]interface{}{"groupid": group})
	}
	for _, template := range templates {
		logger.L().Infof("template: %s", template)
		spec.TemplatesToLink = append(spec.TemplatesToLink, map[string]interface{}{"templateid": template})
	}
	if len(groupMap) > 0 {
		spec.Groups = append(spec.Groups, groupMap...)
	}
	if len(templateMap) > 0 {
		spec.TemplatesToLink = append(spec.TemplatesToLink, templateMap...)
	}

	result, err := server.CreateHost(ctx, clientPool, spec, instanceName)
	if err != nil {
		logger.L().Errorf("调用 host.create 失败: %v", err)
		return nil, fmt.Errorf("调用 host.create 失败: %w", err)
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
		if arr, ok := args["hostids"].([]interface{}); ok {
			for _, v := range arr {
				if s, ok := v.(string); ok && s != "" {
					hostIDs = append(hostIDs, s)
				}
			}
		}
		if v, ok := args["hostid"].(string); ok && v != "" {
			hostIDs = append(hostIDs, v)
		}
	}
	if clientPool == nil {
		return mcp.NewToolResultStructuredOnly(makeResult([]map[string]interface{}{})), nil
	}
	spec := models.HostParams{HostIDs: hostIDs}
	hosts, err := server.DeleteHosts(ctx, clientPool, spec, instanceName)
	if err != nil {
		logger.L().Errorf("调用 host.delete 失败: %w", err)
		return nil, fmt.Errorf("调用 host.delete 失败: %w", err)
	}
	return mcp.NewToolResultStructuredOnly(makeResult(hosts)), nil
}
