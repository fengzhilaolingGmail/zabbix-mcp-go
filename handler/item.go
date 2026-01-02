/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2026-01-02 16:17:56
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2026-01-02 17:50:44
 * @FilePath: \zabbix-mcp-go\handler\item.go
 * @Description: 监控项相关功能
 * Copyright (c) 2026 by fengzhilaoling@gmail.com, All Rights Reserved.
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

func GetItemsHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	instance := ""
	spec := models.ParamsItem{}
	// 用于验证二选一必填的参数
	hasHostFilter := false
	hasItemFilter := false
	if args, ok := req.Params.Arguments.(map[string]interface{}); ok {
		// 解析 instance
		if v, ok2 := args["instance"].(string); ok2 {
			instance = v
		}

		// 解析 ID 数组参数
		if arr, ok := args["itemids"].([]interface{}); ok {
			for _, v := range arr {
				if s, ok := v.(string); ok && s != "" {
					spec.ItemIDs = append(spec.ItemIDs, s)
				}
			}
		}
		if arr, ok := args["groupids"].([]interface{}); ok {
			for _, v := range arr {
				if s, ok := v.(string); ok && s != "" {
					spec.GroupIDs = append(spec.GroupIDs, s)
				}
			}
		}
		if arr, ok := args["templateids"].([]interface{}); ok {
			for _, v := range arr {
				if s, ok := v.(string); ok && s != "" {
					spec.TemplateIDs = append(spec.TemplateIDs, s)
				}
			}
		}
		if arr, ok := args["host_ids"].([]interface{}); ok {
			for _, v := range arr {
				if s, ok := v.(string); ok && s != "" {
					spec.HostIDs = append(spec.HostIDs, s)
				}
			}
			hasHostFilter = true
		}
		if arr, ok := args["proxyids"].([]interface{}); ok {
			for _, v := range arr {
				if s, ok := v.(string); ok && s != "" {
					spec.ProxyIDs = append(spec.ProxyIDs, s)
				}
			}
		}
		if arr, ok := args["interfaceids"].([]interface{}); ok {
			for _, v := range arr {
				if s, ok := v.(string); ok && s != "" {
					spec.InterfaceIDs = append(spec.InterfaceIDs, s)
				}
			}
		}
		if arr, ok := args["graphids"].([]interface{}); ok {
			for _, v := range arr {
				if s, ok := v.(string); ok && s != "" {
					spec.GraphIDs = append(spec.GraphIDs, s)
				}
			}
		}
		if arr, ok := args["triggerids"].([]interface{}); ok {
			for _, v := range arr {
				if s, ok := v.(string); ok && s != "" {
					spec.TriggerIDs = append(spec.TriggerIDs, s)
				}
			}
		}

		// 解析布尔参数
		if v, ok := args["webitems"].(bool); ok {
			spec.WebItems = v
		}
		if v, ok := args["inherited"].(bool); ok {
			spec.Inherited = v
		}
		if v, ok := args["templated"].(bool); ok {
			spec.Templated = v
		}
		if v, ok := args["monitored"].(bool); ok {
			spec.Monitored = v
		}
		if v, ok := args["with_triggers"].(bool); ok {
			spec.WithTriggers = v
		}
		if v, ok := args["countOutput"].(bool); ok {
			spec.CountOutput = v
		}
		if v, ok := args["editable"].(bool); ok {
			spec.Editable = v
		}
		if v, ok := args["excludeSearch"].(bool); ok {
			spec.ExcludeSearch = v
		}
		if v, ok := args["preservekeys"].(bool); ok {
			spec.PreserveKeys = v
		}
		if v, ok := args["searchByAny"].(bool); ok {
			spec.SearchByAny = v
		}
		if v, ok := args["searchWildcardsEnabled"].(bool); ok {
			spec.SearchWildcardsEnabled = v
		}
		if v, ok := args["startSearch"].(bool); ok {
			spec.StartSearch = v
		}

		// 解析字符串参数
		if v, ok := args["group"].(string); ok && v != "" {
			spec.Group = v
		}
		if v, ok := args["hostname"].(string); ok && v != "" {
			spec.Host = v
			hasHostFilter = true
		}

		// 解析整数参数
		if v, ok := args["evaltype"].(float64); ok {
			spec.EvalType = int(v)
		}
		if v, ok := args["limitSelects"].(float64); ok {
			spec.LimitSelects = int(v)
		}
		if v, ok := args["limit"].(float64); ok {
			spec.Limit = int(v)
		}

		// 解析标签数组
		if tags, ok := args["tags"].([]interface{}); ok {
			for _, tag := range tags {
				if tagMap, ok := tag.(map[string]interface{}); ok {
					spec.Tags = append(spec.Tags, tagMap)
				}
			}
		}

		// 解析 select 参数
		if v, ok := args["selectHosts"].(bool); ok && v {
			spec.SelectHosts = "extend"
		}
		if v, ok := args["selectInterfaces"].(bool); ok && v {
			spec.SelectInterfaces = "extend"
		}
		if v, ok := args["selectTriggers"].(bool); ok && v {
			spec.SelectTriggers = "extend"
		}
		if v, ok := args["selectGraphs"].(bool); ok && v {
			spec.SelectGraphs = "extend"
		}
		if v, ok := args["selectDiscoveryData"].(bool); ok && v {
			spec.SelectDiscoveryData = "extend"
		}
		if v, ok := args["selectDiscoveryRule"].(bool); ok && v {
			spec.SelectDiscoveryRule = "extend"
		}
		if v, ok := args["selectDiscoveryRulePrototype"].(bool); ok && v {
			spec.SelectDiscoveryRulePrototype = "extend"
		}
		if v, ok := args["selectItemDiscovery"].(bool); ok && v {
			spec.SelectItemDiscovery = "extend"
		}
		if v, ok := args["selectPreprocessing"].(bool); ok && v {
			spec.SelectPreprocessing = "extend"
		}
		if v, ok := args["selectTags"].(bool); ok && v {
			spec.SelectTags = "extend"
		}
		if v, ok := args["selectValueMap"].(bool); ok && v {
			spec.SelectValueMap = "extend"
		}

		// 解析 filter
		if filter, ok := args["filter"].(map[string]interface{}); ok {
			spec.Filter = filter
		}

		// 解析 search - 先初始化 map
		if spec.Search == nil {
			spec.Search = make(map[string]interface{})
		}
		// 按监控项键搜索
		if v, ok := args["item_key"].(string); ok && v != "" {
			spec.Search["key_"] = v
			hasItemFilter = true
			spec.SearchWildcardsEnabled = true
		}
		// 按监控项名称搜索
		if v, ok := args["item_name"].(string); ok && v != "" {
			spec.Search["name"] = v
			hasItemFilter = true
			spec.SearchWildcardsEnabled = true
		}

		// 解析 sortfield
		if v, ok := args["sortfield"].(string); ok && v != "" {
			spec.SortField = v
		} else if arr, ok := args["sortfield"].([]interface{}); ok {
			var fields []string
			for _, v := range arr {
				if s, ok := v.(string); ok && s != "" {
					fields = append(fields, s)
				}
			}
			if len(fields) > 0 {
				spec.SortField = fields
			}
		}

		// 解析 sortorder
		if v, ok := args["sortorder"].(string); ok && v != "" {
			spec.SortOrder = v
		} else if arr, ok := args["sortorder"].([]interface{}); ok {
			var orders []string
			for _, v := range arr {
				if s, ok := v.(string); ok && s != "" {
					orders = append(orders, s)
				}
			}
			if len(orders) > 0 {
				spec.SortOrder = orders
			}
		}

		// 解析 output
		if v, ok := args["output"].(string); ok && v != "" {
			spec.Output = v
		} else if arr, ok := args["output"].([]interface{}); ok {
			var fields []string
			for _, v := range arr {
				if s, ok := v.(string); ok && s != "" {
					fields = append(fields, s)
				}
			}
			if len(fields) > 0 {
				spec.Output = fields
			}
		} else {
			spec.Output = "extend"
		}
	}
	if !hasHostFilter && !hasItemFilter {
		return nil, fmt.Errorf("必须指定 host_ids 或 host 或 item_key 或 item_name")
	}
	logger.L().Infof("获取监控项列表参数: %v, %v", hasHostFilter, hasItemFilter)

	if clientPool == nil {
		return mcp.NewToolResultStructuredOnly(makeResult([]map[string]interface{}{})), nil
	}

	items, err := server.GetItems(ctx, clientPool, spec, instance)
	if err != nil {
		return nil, fmt.Errorf("调用 item.get 失败: %w", err)
	}
	return mcp.NewToolResultStructuredOnly(makeResult(items)), nil
}
