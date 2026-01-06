/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2026-01-06 18:59:21
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2026-01-06 19:02:54
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

// generateGraphWidgets 根据主机列表和图形ID自动生成 widgets 布局
// generateGraphWidgets 根据主机列表和图形ID自动生成 widgets 布局
// perRow: 每行 widget 个数（<=0 时默认为 2）
// width/height: 每个 widget 的宽度和高度（<=0 时使用默认 36/5）
// startX/startY: 起始偏移
func generateGraphWidgets(hosts []string, graphids []string, perRow, width, height, startX, startY int) ([]models.DashboardWidget, error) {
	widgets := make([]models.DashboardWidget, 0)
	if len(hosts) == 0 || len(graphids) == 0 {
		return widgets, nil
	}

	if perRow <= 0 {
		perRow = 2
	}
	if width <= 0 {
		width = 36
	}
	if height <= 0 {
		height = 5
	}

	// 约束：页面最大宽度/高度（以 Zabbix UI 限制为准）
	const pageWidth = 72
	const pageHeight = 64

	// 计算每行最大允许的 widget 数量，确保不会超出 x 范围
	maxPerRow := (pageWidth - startX) / width
	if maxPerRow <= 0 {
		maxPerRow = 1
	}
	if perRow > maxPerRow {
		perRow = maxPerRow
	}

	total := len(hosts) * len(graphids)
	neededRows := (total + perRow - 1) / perRow
	maxRows := (pageHeight - startY) / height
	if maxRows <= 0 {
		maxRows = 1
	}
	if neededRows > maxRows {
		return nil, fmt.Errorf("布局超出页高度: 需要行 %d，最大允许 %d；请减小 widgetHeight 或增加 rows", neededRows, maxRows)
	}

	widgetIndex := 0
	for _, host := range hosts {
		for _, graphid := range graphids {
			col := widgetIndex % perRow
			row := widgetIndex / perRow

			x := startX + col*width
			y := startY + row*height

			widget := models.DashboardWidget{
				Type:     "graph",
				Name:     fmt.Sprintf("%s - Graph %s", host, graphid),
				X:        x,
				Y:        y,
				Width:    width,
				Height:   height,
				ViewMode: 0,
				Fields: []models.DashboardWidgetField{
					{
						Type:  3, // 主机类型
						Name:  "hostids",
						Value: host,
					},
					{
						Type: 6, // 图形类型（需要数值）
						Name: "graphid",
						Value: func() interface{} {
							if n, err := strconv.Atoi(graphid); err == nil {
								return n
							}
							return graphid
						}(),
					},
				},
			}
			widgets = append(widgets, widget)
			widgetIndex++
		}
	}
	return widgets, nil
}

// CreateDashboardHandler 通过注入的 ClientProvider 调用 dashboard.create 并返回结果
func CreateDashboardHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	instanceName := ""
	name := ""
	pages := make([]models.DashboardPage, 0)
	users := make([]models.DashboardUser, 0)
	userGroups := make([]models.DashboardUserGroup, 0)
	hosts := make([]string, 0)
	graphids := make([]string, 0)

	if args, ok := req.Params.Arguments.(map[string]interface{}); ok {
		if v, ok2 := args["instance"].(string); ok2 {
			instanceName = v
		}
		if v, ok2 := args["name"].(string); ok2 {
			name = v
		}

		// 解析 hosts 和 graphids（用于自动生成 graph widgets）
		if hostsArr, ok2 := args["hosts"].([]interface{}); ok2 {
			for _, hostItf := range hostsArr {
				if hostStr, ok3 := hostItf.(string); ok3 && hostStr != "" {
					hosts = append(hosts, hostStr)
				}
			}
		}
		if graphidsArr, ok2 := args["graphids"].([]interface{}); ok2 {
			for _, graphidItf := range graphidsArr {
				if graphidStr, ok3 := graphidItf.(string); ok3 && graphidStr != "" {
					graphids = append(graphids, graphidStr)
				}
			}
		}

		// 解析 pages
		if pagesArr, ok2 := args["pages"].([]interface{}); ok2 {
			for _, pageItf := range pagesArr {
				if pageMap, ok3 := pageItf.(map[string]interface{}); ok3 {
					page := models.DashboardPage{
						Name:          "",
						DisplayPeriod: 0,
						Widgets:       make([]models.DashboardWidget, 0),
					}

					if pageName, ok4 := pageMap["name"].(string); ok4 {
						page.Name = pageName
					}
					if displayPeriod, ok4 := pageMap["display_period"].(float64); ok4 {
						page.DisplayPeriod = int(displayPeriod)
					}

					// 解析 widgets
					if widgetsArr, ok4 := pageMap["widgets"].([]interface{}); ok4 {
						for _, widgetItf := range widgetsArr {
							if widgetMap, ok5 := widgetItf.(map[string]interface{}); ok5 {
								widget := models.DashboardWidget{
									Type:     "",
									Name:     "",
									X:        0,
									Y:        0,
									Width:    1,
									Height:   1,
									ViewMode: 0,
									Fields:   make([]models.DashboardWidgetField, 0),
								}

								if widgetType, ok6 := widgetMap["type"].(string); ok6 {
									widget.Type = widgetType
								}
								if widgetName, ok6 := widgetMap["name"].(string); ok6 {
									widget.Name = widgetName
								}
								if x, ok6 := widgetMap["x"].(float64); ok6 {
									widget.X = int(x)
								}
								if y, ok6 := widgetMap["y"].(float64); ok6 {
									widget.Y = int(y)
								}
								if width, ok6 := widgetMap["width"].(float64); ok6 {
									widget.Width = int(width)
								}
								if height, ok6 := widgetMap["height"].(float64); ok6 {
									widget.Height = int(height)
								}
								if viewMode, ok6 := widgetMap["view_mode"].(float64); ok6 {
									widget.ViewMode = int(viewMode)
								}

								// 解析 fields
								if fieldsArr, ok6 := widgetMap["fields"].([]interface{}); ok6 {
									for _, fieldItf := range fieldsArr {
										if fieldMap, ok7 := fieldItf.(map[string]interface{}); ok7 {
											field := models.DashboardWidgetField{
												Type:  0,
												Name:  "",
												Value: nil,
											}

											if fieldType, ok8 := fieldMap["type"].(float64); ok8 {
												field.Type = int(fieldType)
											}
											if fieldName, ok8 := fieldMap["name"].(string); ok8 {
												field.Name = fieldName
											}
											if fieldValue, ok8 := fieldMap["value"]; ok8 {
												field.Value = fieldValue
												// 如果是 graphid (type 6)，尝试将字符串数值转换为数字类型，避免 Zabbix 返回 "a number is expected"
												if field.Type == 6 {
													switch v := fieldValue.(type) {
													case string:
														if n, err := strconv.Atoi(v); err == nil {
															field.Value = n
														}
													case float64:
														// 保持为整数类型以符合 API 期望
														field.Value = int(v)
													}
												}
												// 如果是 hostids (type 3)，也尝试转换为数字类型
												if field.Type == 3 {
													switch v := fieldValue.(type) {
													case string:
														if n, err := strconv.Atoi(v); err == nil {
															field.Value = n
														}
													case float64:
														field.Value = int(v)
													case []interface{}:
														// 将数组内的数字尝试转换为 int
														arr := make([]interface{}, 0, len(v))
														for _, it := range v {
															switch vv := it.(type) {
															case string:
																if n, err := strconv.Atoi(vv); err == nil {
																	arr = append(arr, n)
																} else {
																	arr = append(arr, vv)
																}
															case float64:
																arr = append(arr, int(vv))
															default:
																arr = append(arr, vv)
															}
														}
														field.Value = arr
													}
												}
											}

											widget.Fields = append(widget.Fields, field)
										}
									}
								}

								page.Widgets = append(page.Widgets, widget)
							}
						}
					}

					pages = append(pages, page)
				}
			}
		}

		// 解析 users
		if usersArr, ok2 := args["users"].([]interface{}); ok2 {
			for _, userItf := range usersArr {
				if userMap, ok3 := userItf.(map[string]interface{}); ok3 {
					user := models.DashboardUser{
						UserID:     "",
						Permission: 0,
					}

					if userID, ok4 := userMap["userid"].(string); ok4 {
						user.UserID = userID
					}
					if permission, ok4 := userMap["permission"].(float64); ok4 {
						user.Permission = int(permission)
					}

					users = append(users, user)
				}
			}
		}

		// 解析 userGroups
		if userGroupsArr, ok2 := args["userGroups"].([]interface{}); ok2 {
			for _, ugItf := range userGroupsArr {
				if ugMap, ok3 := ugItf.(map[string]interface{}); ok3 {
					ug := models.DashboardUserGroup{
						Usrgrpid:   "",
						Permission: 0,
					}

					if usrgrpid, ok4 := ugMap["usrgrpid"].(string); ok4 {
						ug.Usrgrpid = usrgrpid
					}
					if permission, ok4 := ugMap["permission"].(float64); ok4 {
						ug.Permission = int(permission)
					}

					userGroups = append(userGroups, ug)
				}
			}
		}
	}

	if clientPool == nil {
		return mcp.NewToolResultStructuredOnly(makeResult([]map[string]interface{}{})), nil
	}

	// 如果传入了简化的创建参数（通过 create_graph_dashboard 工具），handler 会被 CreateGraphDashboardHandler 使用，这里仅保留兼容逻辑

	// 如果提供了 hosts 和 graphids，自动生成 graph widgets（使用默认布局）
	if len(hosts) > 0 && len(graphids) > 0 {
		autoWidgets, err := generateGraphWidgets(hosts, graphids, 2, 36, 5, 0, 0)
		if err != nil {
			return nil, fmt.Errorf("自动生成 widgets 失败: %w", err)
		}
		if len(pages) == 0 {
			// 如果没有提供 pages，创建一个默认页面
			pages = []models.DashboardPage{
				{
					Name:          "Default Page",
					DisplayPeriod: 0,
					Widgets:       autoWidgets,
				},
			}
		} else {
			// 如果提供了 pages，将自动生成的 widgets 添加到第一个页面
			pages[0].Widgets = append(pages[0].Widgets, autoWidgets...)
		}
	}

	// 验证：如果没有 pages，返回错误
	if len(pages) == 0 {
		return nil, fmt.Errorf("必须提供 pages 参数，或者提供 hosts 和 graphids 参数来自动生成")
	}

	spec := models.DashboardParams{
		Name:       name,
		Pages:      pages,
		Users:      users,
		UserGroups: userGroups,
	}

	result, err := server.CreateDashboard(ctx, clientPool, spec, instanceName)
	if err != nil {
		logger.L().Errorf("调用 dashboard.create 失败: %v", err)
		return nil, fmt.Errorf("调用 dashboard.create 失败: %w", err)
	}
	return mcp.NewToolResultStructuredOnly(makeResult(result)), nil
}
