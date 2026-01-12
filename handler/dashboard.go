/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2026-01-06 18:59:21
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2026-01-09 16:38:55
 * @FilePath: \zabbix-mcp-go\handler\dashboard.go
 * @Description: 仪表盘相关功能
 * @Copyright: Copyright (c) 2025 by fengzhilaoling@gmail.com, All Rights Reserved.
 */
package handler

import (
	"context"
	"fmt"
	"strconv"

	"zabbixMcp/logger"
	"zabbixMcp/models"
	"zabbixMcp/server"

	"github.com/mark3labs/mcp-go/mcp"
)

// CreateDashboardHandler 通过注入的 ClientProvider 调用 dashboard.create 并返回结果
func CreateDashboardHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	instanceName := ""
	if args, ok := req.Params.Arguments.(map[string]interface{}); ok {
		if v, ok2 := args["instance"].(string); ok2 {
			instanceName = v
		}
	}

	var argsMap map[string]interface{}
	if a, ok := req.Params.Arguments.(map[string]interface{}); ok {
		argsMap = a
	}

	spec, err := buildDashboardSpecFromArgs(argsMap)
	if err != nil {
		logger.L().Warnf("解析 dashboard 参数失败: %v", err)
		return nil, fmt.Errorf("解析 dashboard 参数失败: %w", err)
	}

	// override users/userGroups if explicit arrays passed (maintain backward compat)
	if argsMap != nil {
		if arr, ok := argsMap["users"].([]interface{}); ok && len(arr) > 0 {
			users := make([]models.DashboardUser, 0, len(arr))
			for _, u := range arr {
				if um, ok := u.(map[string]interface{}); ok {
					du := models.DashboardUser{}
					if id, ok2 := um["userid"].(string); ok2 {
						du.UserID = id
					}
					if p, ok2 := um["permission"].(float64); ok2 {
						du.Permission = int(p)
					}
					users = append(users, du)
				}
			}
			spec.Users = users
		}
		if arr, ok := argsMap["userGroups"].([]interface{}); ok && len(arr) > 0 {
			ugs := make([]models.DashboardUserGroup, 0, len(arr))
			for _, g := range arr {
				if gm, ok := g.(map[string]interface{}); ok {
					dug := models.DashboardUserGroup{}
					if id, ok2 := gm["usrgrpid"].(string); ok2 {
						dug.Usrgrpid = id
					}
					if p, ok2 := gm["permission"].(float64); ok2 {
						dug.Permission = int(p)
					}
					ugs = append(ugs, dug)
				}
			}
			spec.UserGroups = ugs
		}
	}

	if clientPool == nil {
		return mcp.NewToolResultStructuredOnly(makeResult(map[string]interface{}{})), nil
	}

	// 批量解析 widget 内传入的 hostnames（如果有），并把值替换为数字 hostid 或数字数组
	// 收集需要解析的主机名
	hostNamesSet := map[string]struct{}{}
	for pi := range spec.Pages {
		pg := &spec.Pages[pi]
		for wi := range pg.Widgets {
			w := &pg.Widgets[wi]
			for fi := range w.Fields {
				f := &w.Fields[fi]
				if f.Name == "hostids" {
					switch v := f.Value.(type) {
					case string:
						if v != "" {
							hostNamesSet[v] = struct{}{}
						}
					case []interface{}:
						for _, e := range v {
							if s, ok := e.(string); ok && s != "" {
								hostNamesSet[s] = struct{}{}
							}
						}
					}
				}
			}
		}
	}

	if len(hostNamesSet) > 0 {
		hostNames := make([]string, 0, len(hostNamesSet))
		for hn := range hostNamesSet {
			hostNames = append(hostNames, hn)
		}
		specHosts := models.HostParams{Output: []string{"hostid", "host"}}
		specHosts.Filter = map[string]interface{}{"host": hostNames}
		hostList, err := server.GetHosts(ctx, clientPool, specHosts, instanceName)
		hostNameToID := map[string]string{}
		if err == nil {
			for _, h := range hostList {
				hn := ""
				hid := ""
				if v, ok := h["host"].(string); ok {
					hn = v
				}
				if v, ok := h["hostid"].(string); ok {
					hid = v
				}
				if hn != "" && hid != "" {
					hostNameToID[hn] = hid
				}
			}
		} else {
			logger.L().Warnf("批量解析 widget hostnames 失败: %v", err)
		}

		// 替换字段值
		for pi := range spec.Pages {
			pg := &spec.Pages[pi]
			for wi := range pg.Widgets {
				w := &pg.Widgets[wi]
				for fi := range w.Fields {
					f := &w.Fields[fi]
					if f.Name != "hostids" {
						continue
					}
					switch v := f.Value.(type) {
					case string:
						if v == "" {
							continue
						}
						// 如果是数字字符串直接转换
						if idNum, err := strconv.Atoi(v); err == nil {
							f.Value = idNum
							continue
						}
						if hid, ok := hostNameToID[v]; ok {
							if idNum, err := strconv.Atoi(hid); err == nil {
								f.Value = idNum
							} else {
								logger.L().Warnf("解析到的 hostid %s 不是数字，保留原值", hid)
							}
						} else {
							logger.L().Warnf("未解析到主机 %s 的 hostid，保留原值", v)
						}
					case []interface{}:
						outNums := make([]interface{}, 0, len(v))
						for _, e := range v {
							switch ev := e.(type) {
							case string:
								if ev == "" {
									continue
								}
								if idNum, err := strconv.Atoi(ev); err == nil {
									outNums = append(outNums, idNum)
									continue
								}
								if hid, ok := hostNameToID[ev]; ok {
									if idNum, err := strconv.Atoi(hid); err == nil {
										outNums = append(outNums, idNum)
									} else {
										logger.L().Warnf("解析到的 hostid %s 不是数字，跳过", hid)
									}
								} else {
									logger.L().Warnf("未解析到主机 %s 的 hostid，跳过", ev)
								}
							case float64:
								outNums = append(outNums, int(ev))
							case int:
								outNums = append(outNums, ev)
							}
						}
						if len(outNums) == 1 {
							f.Value = outNums[0]
						} else if len(outNums) > 1 {
							f.Value = outNums
						} else {
							logger.L().Warnf("widget %s 中 hostids 未解析出有效 id，保留原值", w.Name)
						}
					}
				}
			}
		}
	}

	dashboard, err := server.CreateDashboard(ctx, clientPool, spec, instanceName)
	if err != nil {
		logger.L().Errorf("调用 dashboard.create 失败: %v", err)
		return nil, fmt.Errorf("调用 dashboard.create 失败: %w", err)
	}
	return mcp.NewToolResultStructuredOnly(makeResult(dashboard)), nil
}

func GetDashboardHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	instance := ""
	dashboard_name := ""
	select_users := false
	select_usergroups := false
	select_pages := false
	if args, ok := req.Params.Arguments.(map[string]interface{}); ok {
		if v, ok2 := args["instance"].(string); ok2 {
			instance = v
		}
		if v, ok2 := args["dashboard_name"].(string); ok2 {
			dashboard_name = v
		}
		if v, ok2 := args["select_users"].(bool); ok2 {
			select_users = v
		}
		if v, ok2 := args["select_usergroups"].(bool); ok2 {
			select_usergroups = v
		}
		if v, ok2 := args["select_pages"].(bool); ok2 {
			select_pages = v
		}
	}
	if clientPool == nil {
		return mcp.NewToolResultStructuredOnly(makeResult([]map[string]interface{}{})), nil
	}
	spec := models.DashboardParams{Output: []string{"dashboardid", "name"}}
	if dashboard_name != "" {
		spec.Filter = map[string]interface{}{"name": dashboard_name}
	}
	if select_users {
		spec.SelectUsers = "extend"
		spec.Output = "extend"
	}
	if select_usergroups {
		spec.SelectUserGroups = "extend"
		spec.Output = "extend"
	}
	if select_pages {
		spec.SelectPages = "extend"
		spec.Output = "extend"
	}

	dashboard, err := server.GetDashboard(ctx, clientPool, spec, instance)
	if err != nil {
		return nil, fmt.Errorf("调用 dashboard.get 失败: %w", err)
	}
	return mcp.NewToolResultStructuredOnly(makeResult(dashboard)), nil
}

// DeleteDashboardHandler 删除仪表盘
func DeleteDashboardHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	instanceName := ""
	dashboardIDs := []string{}
	if args, ok := req.Params.Arguments.(map[string]interface{}); ok {
		if v, ok2 := args["instance"].(string); ok2 {
			instanceName = v
		}
		if arr, ok := args["dashboard_ids"].([]interface{}); ok {
			for _, v := range arr {
				if s, ok := v.(string); ok && s != "" {
					dashboardIDs = append(dashboardIDs, s)
				}
				if s, ok := v.(int64); ok && s != 0 {
					dashboardIDs = append(dashboardIDs, strconv.FormatInt(s, 10))
				}
			}
		}

		if v, ok := args["dashboard_ids"].(string); ok && v != "" {
			dashboardIDs = append(dashboardIDs, v)
		}
		if v, ok := args["dashboard_ids"].(int64); ok && v != 0 {
			dashboardIDs = append(dashboardIDs, strconv.FormatInt(v, 10))
		}
	}
	if clientPool == nil {
		return mcp.NewToolResultStructuredOnly(makeResult([]map[string]interface{}{})), nil
	}
	spec := models.DashboardParams{DashboardIDs: dashboardIDs}
	logger.L().Infof("调用 dashboard.delete 传入参数: %v", spec)
	logger.L().Infof("传入参数: %v", dashboardIDs)
	dashboard, err := server.DeleteDashboard(ctx, clientPool, spec, instanceName)
	if err != nil {
		logger.L().Errorf("调用 dashboard.delete 失败: %w", err)
		return nil, fmt.Errorf("调用 dashboard.delete 失败: %w", err)
	}
	return mcp.NewToolResultStructuredOnly(makeResult(dashboard)), nil
}
