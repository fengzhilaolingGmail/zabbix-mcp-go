/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2025-12-18 11:20:36
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2026-01-27 19:26:10
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
	var selectHostGroups interface{}
	var selectParentTemplates interface{}
	var selectItems interface{}
	var selectMacros interface{}
	var selectTags interface{}
	var selectTriggers interface{}
	var selectDashboards interface{}
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
		selectGraphs = utils.ParseZabbixSelectParam(args, "selectGraphs", "extend")
		selectHostGroups = utils.ParseZabbixSelectParam(args, "selectHostGroups", "extend")
		selectParentTemplates = utils.ParseZabbixSelectParam(args, "selectParentTemplates", "extend")
		selectItems = utils.ParseZabbixSelectParam(args, "selectItems", "extend")
		selectMacros = utils.ParseZabbixSelectParam(args, "selectMacros", "extend")
		selectTags = utils.ParseZabbixSelectParam(args, "selectTags", "extend")
		selectTriggers = utils.ParseZabbixSelectParam(args, "selectTriggers", "extend")
		selectDashboards = utils.ParseZabbixSelectParam(args, "selectDashboards", "extend")
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
	logger.L().Infof("GetHostForNameHandler: instance=%s, hostname=%s, searchWildcardsEnabled=%v, limit=%d",
		instance, hostname, searchWildcardsEnabled, limit)
	logger.L().Infof("GetHostForNameHandler: selectGraphs=%v, limitSelects=%d, count=%v, selectHostGroups=%v",
		selectGraphs, limitSelects, count, selectHostGroups)
	logger.L().Infof("GetHostForNameHandler: selectParentTemplates=%v, selectItems=%v, selectMacros=%v, selectTags=%v, selectTriggers=%v, selectDashboards=%v",
		selectParentTemplates, selectItems, selectMacros, selectTags, selectTriggers, selectDashboards)
	spec := models.HostParams{
		Output:                 "extend",
		SelectInterfaces:       "extend",
		Search:                 map[string]interface{}{"host": hostname},
		Limit:                  limit,
		LimitSelects:           limitSelects,
		SearchWildcardsEnabled: searchWildcardsEnabled,
	}
	// 不同工具启用不同参数
	if mcpToolName == "get_graph_by_hostname" {
		if count {
			spec.SelectGraphs = "count"
		} else {
			spec.SelectGraphs = selectGraphs
		}
	}
	if mcpToolName == "get_groups_by_hostname" {
		spec.SelectHostGroups = selectHostGroups
	}
	if mcpToolName == "get_templates_by_hostname" {
		if count {
			spec.SelectParentTemplates = "count"
		} else {
			spec.SelectParentTemplates = selectParentTemplates
		}
	}
	if mcpToolName == "get_items_by_hostname" {
		if count {
			spec.SelectItems = "count"
		} else {
			spec.SelectItems = selectItems
		}
	}
	if mcpToolName == "get_macros_by_hostname" {
		spec.SelectMacros = selectMacros
	}
	if mcpToolName == "get_tags_by_hostname" {
		spec.SelectTags = selectTags
	}
	if mcpToolName == "get_triggers_by_hostname" {
		if count {
			spec.SelectTriggers = "count"
		} else {
			spec.SelectTriggers = selectTriggers
		}
	}
	if mcpToolName == "get_dashboards_by_hostname" {
		if count {
			spec.SelectDashboards = "count"
		} else {
			spec.SelectDashboards = selectDashboards
		}
	}
	hosts, err := server.GetHosts(ctx, clientPool, spec, instance)
	if err != nil {
		return nil, fmt.Errorf("调用 host.get 失败: %w", err)
	}
	return mcp.NewToolResultStructuredOnly(makeResult(hosts)), nil
}

// 创建主机
func CreateHostHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	instance := ""
	host := ""
	name := ""
	groups := []models.Groups{}
	templates := []models.Templates{}
	tags := make([]models.Tag, 0)
	macros := []models.Macros{}
	// inventory := make(map[string]interface{})
	interfaces := []models.ZabbixInterface{}
	if args, ok := req.Params.Arguments.(map[string]interface{}); ok {
		if v, ok2 := args["instance"].(string); ok2 {
			instance = v
		}
		if v, ok2 := args["host"].(string); ok2 {
			host = v
		}
		if v, ok2 := args["name"].(string); ok2 {
			name = v
		}
		// 解析 groups（支持 []models.Groups 或 []map[string]interface{}）
		if err := utils.ParseSliceFromMap(args, "groups", &groups); err != nil {
			if arr, ok2 := args["groups"].([]interface{}); ok2 {
				for _, v2 := range arr {
					switch vv := v2.(type) {
					case string:
						groups = append(groups, models.Groups{GroupID: vv})
					case map[string]interface{}:
						if gid, ok3 := vv["groupid"].(string); ok3 {
							groups = append(groups, models.Groups{GroupID: gid})
						}
					}
				}
			}
		}

		// 解析 templates（支持 []models.Templates 或 []map[string]interface{}）
		if err := utils.ParseSliceFromMap(args, "templates", &templates); err != nil {
			if arr, ok2 := args["templates"].([]interface{}); ok2 {
				for _, v2 := range arr {
					switch vv := v2.(type) {
					case string:
						templates = append(templates, models.Templates{TemplateID: vv})
					case map[string]interface{}:
						if tid, ok3 := vv["templateid"].(string); ok3 {
							templates = append(templates, models.Templates{TemplateID: tid})
						}
					}
				}
			}
		}
		// 解析 tags，支持多种输入格式（[]models.Tag、[]map[string]interface{} 等）
		if err := utils.ParseSliceFromMap(args, "tags", &tags); err != nil {
			// 回退到手动解析以保持兼容性
			if arr, ok2 := args["tags"].([]interface{}); ok2 {
				for _, v2 := range arr {
					if m, ok3 := v2.(map[string]interface{}); ok3 {
						tagStr, _ := m["tag"].(string)
						valStr, _ := m["value"].(string)
						tags = append(tags, models.Tag{Tag: tagStr, Value: valStr})
					}
				}
			}
		}

		// 解析 macros，支持 []models.Macros 或 []map[string]interface{} 等
		if err := utils.ParseSliceFromMap(args, "macros", &macros); err != nil {
			if arr, ok2 := args["macros"].([]interface{}); ok2 {
				for _, v2 := range arr {
					switch vv := v2.(type) {
					case string:
						macros = append(macros, models.Macros{Macro: vv})
					case map[string]interface{}:
						macroStr, _ := vv["macro"].(string)
						valueStr, _ := vv["value"].(string)
						descStr, _ := vv["description"].(string)
						macros = append(macros, models.Macros{Macro: macroStr, Value: valueStr, Description: descStr})
					}
				}
			}
		}
		// 仅在参数中存在 interfaces 时才解析，避免当请求中没有该字段时报错
		if _, has := args["interfaces"]; has {
			if err := utils.ParseSliceFromMap(args, "interfaces", &interfaces); err != nil {
				// 业务层错误处理（日志+返回/终止）
				logger.L().Errorf("解析interfaces失败: %v", err)
				return nil, fmt.Errorf("解析interfaces失败: %w", err)
			}
		}
	}
	logger.L().Infof("CreateHostHandler: instance=%s, host=%s, name=%s, groups=%v, templates=%v, tags=%v, macros=%v, interfaces=%v",
		instance, host, name, groups, templates, tags, macros, interfaces)
	if clientPool == nil {
		return mcp.NewToolResultStructuredOnly(makeResult([]map[string]interface{}{})), nil
	}
	spec := models.HostParams{
		Host:       host,
		Name:       name,
		Interfaces: interfaces,
		Tags:       tags,
		Macros:     macros,
		Groups:     groups,
		Templates:  templates,
	}
	result, err := server.CreateHost(ctx, clientPool, spec, instance)
	if err != nil {
		logger.L().Errorf("调用 host.create 失败: %v", err)
		return nil, fmt.Errorf("调用 host.create 失败: %w", err)
	}
	return mcp.NewToolResultStructuredOnly(makeResult(result)), nil
}

func UpdateHostHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	instance := ""
	hostid := ""
	host := ""
	name := ""
	groups := []models.Groups{}
	templates := []models.Templates{}
	tags := make([]models.Tag, 0)
	macros := []models.Macros{}
	// inventory := make(map[string]interface{})
	interfaces := []models.ZabbixInterface{}
	options := ""
	if args, ok := req.Params.Arguments.(map[string]interface{}); ok {
		if v, ok2 := args["instance"].(string); ok2 {
			instance = v
		}
		if v, ok2 := args["hostid"].(string); ok2 {
			hostid = v
		}
		if v, ok2 := args["name"].(string); ok2 {
			name = v
		}
		if v, ok2 := args["host"].(string); ok2 {
			host = v
		}
		if v, ok2 := args["options"].(string); ok2 {
			options = v
		}
		// 解析 groups（支持 []models.Groups 或 []map[string]interface{}）
		if err := utils.ParseSliceFromMap(args, "groups", &groups); err != nil {
			if arr, ok2 := args["groups"].([]interface{}); ok2 {
				for _, v2 := range arr {
					switch vv := v2.(type) {
					case string:
						groups = append(groups, models.Groups{GroupID: vv})
					case map[string]interface{}:
						if gid, ok3 := vv["groupid"].(string); ok3 {
							groups = append(groups, models.Groups{GroupID: gid})
						}
					}
				}
			}
		}

		// 解析 templates（支持 []models.Templates 或 []map[string]interface{}）
		if err := utils.ParseSliceFromMap(args, "templates", &templates); err != nil {
			if arr, ok2 := args["templates"].([]interface{}); ok2 {
				for _, v2 := range arr {
					switch vv := v2.(type) {
					case string:
						templates = append(templates, models.Templates{TemplateID: vv})
					case map[string]interface{}:
						if tid, ok3 := vv["templateid"].(string); ok3 {
							templates = append(templates, models.Templates{TemplateID: tid})
						}
					}
				}
			}
		}
		// 解析 tags，支持多种输入格式（[]models.Tag、[]map[string]interface{} 等）
		if err := utils.ParseSliceFromMap(args, "tags", &tags); err != nil {
			// 回退到手动解析以保持兼容性
			if arr, ok2 := args["tags"].([]interface{}); ok2 {
				for _, v2 := range arr {
					if m, ok3 := v2.(map[string]interface{}); ok3 {
						tagStr, _ := m["tag"].(string)
						valStr, _ := m["value"].(string)
						tags = append(tags, models.Tag{Tag: tagStr, Value: valStr})
					}
				}
			}
		}

		// 解析 macros，支持 []models.Macros 或 []map[string]interface{} 等
		if err := utils.ParseSliceFromMap(args, "macros", &macros); err != nil {
			if arr, ok2 := args["macros"].([]interface{}); ok2 {
				for _, v2 := range arr {
					switch vv := v2.(type) {
					case string:
						macros = append(macros, models.Macros{Macro: vv})
					case map[string]interface{}:
						macroStr, _ := vv["macro"].(string)
						valueStr, _ := vv["value"].(string)
						descStr, _ := vv["description"].(string)
						macros = append(macros, models.Macros{Macro: macroStr, Value: valueStr, Description: descStr})
					}
				}
			}
		}
		// 仅在参数中存在 interfaces 时才解析，避免当请求中没有该字段时报错
		if _, has := args["interfaces"]; has {
			if err := utils.ParseSliceFromMap(args, "interfaces", &interfaces); err != nil {
				// 业务层错误处理（日志+返回/终止）
				logger.L().Errorf("解析interfaces失败: %v", err)
				return nil, fmt.Errorf("解析interfaces失败: %w", err)
			}
		}
	}
	logger.L().Infof("UpdateHostHandler: instance=%s, hostid=%s, host=%s, name=%s, groups=%v, templates=%v",
		instance, hostid, host, name, groups, templates)
	logger.L().Infof("UpdateHostHandler: tags=%v, macros=%v, interfaces=%v", tags, macros, interfaces)
	if clientPool == nil {
		return mcp.NewToolResultStructuredOnly(makeResult([]map[string]interface{}{})), nil
	}
	spec := models.HostParams{
		HostId: hostid,
	}
	if name != "" {
		spec.Name = name
	}
	if host != "" {
		spec.Host = host
	}
	switch options {
	case "templates":
		spec.Templates = templates
	case "groups":
		spec.Groups = groups
	case "tags":
		spec.Tags = tags
	case "macros":
		spec.Macros = macros
	case "interfaces":
		spec.Interfaces = interfaces
	case "clear":
		spec.TemplatesClear = templates
	default:
		return nil, fmt.Errorf("未知的更新选项: %s", options)
	}
	result, err := server.UpdateHost(ctx, clientPool, spec, instance)
	if err != nil {
		logger.L().Errorf("调用 host.update 失败: %v", err)
		return nil, fmt.Errorf("调用 host.update 失败: %w", err)
	}
	return mcp.NewToolResultStructuredOnly(makeResult(result)), nil
}

// 删除主机
func DeleteHostsHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	instance := ""
	hostIDs := []string{}
	if args, ok := req.Params.Arguments.(map[string]interface{}); ok {
		if v, ok2 := args["instance"].(string); ok2 {
			instance = v
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
	hosts, err := server.DeleteHosts(ctx, clientPool, spec, instance)
	if err != nil {
		logger.L().Errorf("调用 host.delete 失败: %w", err)
		return nil, fmt.Errorf("调用 host.delete 失败: %w", err)
	}
	return mcp.NewToolResultStructuredOnly(makeResult(hosts)), nil
}
