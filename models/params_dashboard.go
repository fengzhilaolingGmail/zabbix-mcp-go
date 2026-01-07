/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2026-01-06 18:50:15
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2026-01-07 12:20:00
 * @FilePath: \zabbix-mcp-go\models\params_dashboard.go
 * @Description: 仪表盘参数（参考 Zabbix 7.4 dashboard object）
 * Copyright (c) 2025 by fengzhilaoling@gmail.com, All Rights Reserved.
 */
package models

import "strconv"

// DashboardWidgetField 仪表盘小部件字段
type DashboardWidgetField struct {
	Type  int         `json:"type"`  // 小部件字段类型 (0-13)
	Name  string      `json:"name"`  // 小部件字段名称
	Value interface{} `json:"value"` // 小部件字段值（mixed）
}

// DashboardWidget 仪表盘小部件
type DashboardWidget struct {
	WidgetID string `json:"widgetid,omitempty"` // 只读 - 小部件ID
	// Type defines widget type. Use WidgetType constants when possible.
	Type     WidgetType             `json:"type"`                // 必填 - 小部件类型
	Name     string                 `json:"name,omitempty"`      // 自定义小部件名称
	X        *int                   `json:"x,omitempty"`         // 水平位置 (0-71)
	Y        *int                   `json:"y,omitempty"`         // 垂直位置 (0-63)
	Width    *int                   `json:"width,omitempty"`     // 宽度 (1-72)
	Height   *int                   `json:"height,omitempty"`    // 高度 (1-64)
	ViewMode *int                   `json:"view_mode,omitempty"` // 视图模式 (0-1)
	Fields   []DashboardWidgetField `json:"fields,omitempty"`    // 小部件字段
}

// WidgetType 仪表盘小部件类型（参考 Zabbix 7.4 文档）
type WidgetType string

const (
	WidgetTypeActionLog          WidgetType = "actionlog"
	WidgetTypeClock              WidgetType = "clock"
	WidgetTypeDiscovery          WidgetType = "discovery"
	WidgetTypeFavGraphs          WidgetType = "favgraphs"
	WidgetTypeFavMaps            WidgetType = "favmaps"
	WidgetTypeGauge              WidgetType = "gauge"
	WidgetTypeGeoMap             WidgetType = "geomap"
	WidgetTypeGraph              WidgetType = "graph"
	WidgetTypeGraphProto         WidgetType = "graphprototype"
	WidgetTypeHoneycomb          WidgetType = "honeycomb"
	WidgetTypeHostAvail          WidgetType = "hostavail"
	WidgetTypeHostCard           WidgetType = "hostcard"
	WidgetTypeHostNav            WidgetType = "hostnavigator"
	WidgetTypeItemCard           WidgetType = "itemcard"
	WidgetTypeItemHistory        WidgetType = "itemhistory"
	WidgetTypeItemNav            WidgetType = "itemnavigator"
	WidgetTypeItemValue          WidgetType = "item" // 监控项
	WidgetTypeMap                WidgetType = "map"
	WidgetTypeNavTree            WidgetType = "navtree"
	WidgetTypePieChart           WidgetType = "piechart"
	WidgetTypeProblemHosts       WidgetType = "problemhosts"
	WidgetTypeProblems           WidgetType = "problems"
	WidgetTypeProblemsBySeverity WidgetType = "problemsbysv"
	WidgetTypeSLAReport          WidgetType = "slareport"
	WidgetTypeSVGGraph           WidgetType = "svggraph"
	WidgetTypeSystemInfo         WidgetType = "systeminfo"
	WidgetTypeTopHosts           WidgetType = "tophosts"
	WidgetTypeTopItems           WidgetType = "topitems"
	WidgetTypeTopTriggers        WidgetType = "toptriggers"
	WidgetTypeTrigOver           WidgetType = "trigover"
	WidgetTypeURL                WidgetType = "url"
	WidgetTypeWeb                WidgetType = "web"
)

// IsValid 返回 widget type 是否为已知类型
func (wt WidgetType) IsValid() bool {
	switch wt {
	case WidgetTypeActionLog, WidgetTypeClock, WidgetTypeDiscovery, WidgetTypeFavGraphs,
		WidgetTypeFavMaps, WidgetTypeGauge, WidgetTypeGeoMap, WidgetTypeGraph,
		WidgetTypeGraphProto, WidgetTypeHoneycomb, WidgetTypeHostAvail, WidgetTypeHostCard,
		WidgetTypeHostNav, WidgetTypeItemCard, WidgetTypeItemHistory, WidgetTypeItemNav,
		WidgetTypeItemValue, WidgetTypeMap, WidgetTypeNavTree, WidgetTypePieChart,
		WidgetTypeProblemHosts, WidgetTypeProblems, WidgetTypeProblemsBySeverity, WidgetTypeSLAReport,
		WidgetTypeSVGGraph, WidgetTypeSystemInfo, WidgetTypeTopHosts, WidgetTypeTopItems,
		WidgetTypeTopTriggers, WidgetTypeTrigOver, WidgetTypeURL, WidgetTypeWeb:
		return true
	}
	return false
}

// isValidWidgetFieldType 检查 Dashboard widget field 的 type 是否在 0-13 范围内
// IsValidWidgetFieldType checks whether dashboard widget field type is within 0-13
func IsValidWidgetFieldType(t int) bool {
	return t >= 0 && t <= 13
}

// DashboardPage 仪表盘页面
type DashboardPage struct {
	DashboardPageID string            `json:"dashboard_pageid,omitempty"` // 只读 - 页面ID
	Name            string            `json:"name,omitempty"`             // 页面名称
	DisplayPeriod   *int              `json:"display_period,omitempty"`   // 显示周期 (秒)
	Widgets         []DashboardWidget `json:"widgets,omitempty"`          // 小部件数组
}

// DashboardUser 仪表盘用户共享
type DashboardUser struct {
	UserID     string `json:"userid"`     // 用户ID
	Permission int    `json:"permission"` // 权限级别 (2 只读, 3 读写)
}

// DashboardUserGroup 仪表盘用户组共享
type DashboardUserGroup struct {
	Usrgrpid   string `json:"usrgrpid"`   // 用户组ID
	Permission int    `json:"permission"` // 权限级别 (2 只读, 3 读写)
}

// DashboardParams 仪表盘参数
type DashboardParams struct {
	// 基本字段
	DashboardID string `json:"dashboardid,omitempty"` // 仪表盘ID (用于更新)
	Name        string `json:"name"`                  // 仪表盘名称 (create 必填)

	// 顶层显示/权限相关（指针以便可以显式发送 0 值）
	UserID        string `json:"userid,omitempty"`         // 仪表板所有者用户ID
	Private       *int   `json:"private,omitempty"`        // 仪表板共享类型 (0 公共, 1 私有)
	DisplayPeriod *int   `json:"display_period,omitempty"` // 默认页面显示周期(秒)
	AutoStart     *int   `json:"auto_start,omitempty"`     // 自动启动幻灯片播放 (0/1)

	// 必填参数
	Pages      []DashboardPage      `json:"pages"`                // 仪表盘页面数组 (必填)
	Users      []DashboardUser      `json:"users,omitempty"`      // 用户共享
	UserGroups []DashboardUserGroup `json:"userGroups,omitempty"` // 用户组共享

	// 查询参数
	DashboardIDs     []string    `json:"dashboardids,omitempty"`     // 仪表盘ID数组
	SelectPages      interface{} `json:"selectPages,omitempty"`      // 选择页面
	SelectUsers      interface{} `json:"selectUsers,omitempty"`      // 选择用户
	SelectUserGroups interface{} `json:"selectUserGroups,omitempty"` // 选择用户组

	// 通用查询参数
	CountOutput            bool                   `json:"countOutput,omitempty"`
	Editable               bool                   `json:"editable,omitempty"`
	Filter                 map[string]interface{} `json:"filter,omitempty"`
	Limit                  int                    `json:"limit,omitempty"`
	Output                 interface{}            `json:"output,omitempty"`
	PreserveKeys           bool                   `json:"preservekeys,omitempty"`
	Search                 map[string]interface{} `json:"search,omitempty"`
	SearchWildcardsEnabled bool                   `json:"searchWildcardsEnabled,omitempty"`
	SortField              interface{}            `json:"sortfield,omitempty"`
	SortOrder              interface{}            `json:"sortorder,omitempty"`
}

// BuildParams 将 DashboardParams 转换为 Zabbix dashboard 所需参数
func (p DashboardParams) BuildParams() map[string]interface{} {
	params := map[string]interface{}{}

	// 基本字段
	if p.DashboardID != "" {
		params["dashboardid"] = p.DashboardID
	}
	if p.Name != "" {
		params["name"] = p.Name
	}

	// 顶层可选字段
	if p.UserID != "" {
		params["userid"] = p.UserID
	}
	if p.Private != nil {
		params["private"] = *p.Private
	}
	if p.DisplayPeriod != nil {
		params["display_period"] = *p.DisplayPeriod
	}
	if p.AutoStart != nil {
		params["auto_start"] = *p.AutoStart
	}

	// 必填参数 - pages
	if len(p.Pages) > 0 {
		pages := make([]map[string]interface{}, 0, len(p.Pages))
		for _, page := range p.Pages {
			pageMap := map[string]interface{}{
				"name": page.Name,
			}

			if page.DisplayPeriod != nil {
				pageMap["display_period"] = *page.DisplayPeriod
			}

			// 处理 widgets
			if len(page.Widgets) > 0 {
				widgets := make([]map[string]interface{}, 0, len(page.Widgets))
				for _, widget := range page.Widgets {
					// skip widgets with unknown/unsupported types
					if widget.Type == "" || !widget.Type.IsValid() {
						continue
					}

					widgetMap := map[string]interface{}{
						"type": widget.Type,
						"name": widget.Name,
					}

					if widget.X != nil {
						widgetMap["x"] = *widget.X
					}
					if widget.Y != nil {
						widgetMap["y"] = *widget.Y
					}
					if widget.Width != nil {
						widgetMap["width"] = *widget.Width
					}
					if widget.Height != nil {
						widgetMap["height"] = *widget.Height
					}
					if widget.ViewMode != nil {
						widgetMap["view_mode"] = *widget.ViewMode
					}

					// 处理 fields -> parameters；过滤无效的 field type
					if len(widget.Fields) > 0 {
						// We still produce "fields" array expected by API, normalizing values
						fields := make([]map[string]interface{}, 0, len(widget.Fields))
						for _, field := range widget.Fields {
							// skip invalid field types
							if !IsValidWidgetFieldType(field.Type) {
								continue
							}
							var normalized interface{}
							switch field.Name {
							case "graphid":
								switch v := field.Value.(type) {
								case string:
									if n, err := strconv.Atoi(v); err == nil {
										normalized = n
									} else {
										normalized = v
									}
								case float64:
									normalized = int(v)
								default:
									normalized = v
								}
							case "hostids":
								switch v := field.Value.(type) {
								case string:
									if n, err := strconv.Atoi(v); err == nil {
										normalized = []int{n}
									} else {
										normalized = []string{v}
									}
								case float64:
									normalized = []int{int(v)}
								case []interface{}:
									arr := make([]interface{}, 0, len(v))
									for _, it := range v {
										switch itv := it.(type) {
										case string:
											if n, err := strconv.Atoi(itv); err == nil {
												arr = append(arr, n)
											} else {
												arr = append(arr, itv)
											}
										case float64:
											arr = append(arr, int(itv))
										default:
											arr = append(arr, itv)
										}
									}
									normalized = arr
								default:
									normalized = field.Value
								}
							default:
								normalized = field.Value
							}
							fieldMap := map[string]interface{}{
								"type":  field.Type,
								"name":  field.Name,
								"value": normalized,
							}
							fields = append(fields, fieldMap)
						}
						if len(fields) > 0 {
							widgetMap["fields"] = fields
						}
					}

					widgets = append(widgets, widgetMap)
				}
				pageMap["widgets"] = widgets
			}

			pages = append(pages, pageMap)
		}
		params["pages"] = pages
	}

	// 用户共享
	if len(p.Users) > 0 {
		users := make([]map[string]interface{}, 0, len(p.Users))
		for _, user := range p.Users {
			userMap := map[string]interface{}{
				"userid":     user.UserID,
				"permission": user.Permission,
			}
			users = append(users, userMap)
		}
		params["users"] = users
	}

	// 用户组共享
	if len(p.UserGroups) > 0 {
		userGroups := make([]map[string]interface{}, 0, len(p.UserGroups))
		for _, ug := range p.UserGroups {
			ugMap := map[string]interface{}{
				"usrgrpid":   ug.Usrgrpid,
				"permission": ug.Permission,
			}
			userGroups = append(userGroups, ugMap)
		}
		params["userGroups"] = userGroups
	}

	// 查询参数
	if len(p.DashboardIDs) > 0 {
		params["dashboardids"] = append([]string(nil), p.DashboardIDs...)
	}
	if p.SelectPages != nil {
		params["selectPages"] = p.SelectPages
	}
	if p.SelectUsers != nil {
		params["selectUsers"] = p.SelectUsers
	}
	if p.SelectUserGroups != nil {
		params["selectUserGroups"] = p.SelectUserGroups
	}

	// 通用查询参数
	if p.CountOutput {
		params["countOutput"] = true
	}
	if p.Editable {
		params["editable"] = true
	}
	if len(p.Filter) > 0 {
		filter := make(map[string]interface{}, len(p.Filter))
		for k, v := range p.Filter {
			filter[k] = v
		}
		params["filter"] = filter
	}
	if p.Limit > 0 {
		params["limit"] = p.Limit
	}
	if p.Output != nil {
		params["output"] = p.Output
	}
	if p.PreserveKeys {
		params["preservekeys"] = true
	}
	if len(p.Search) > 0 {
		search := make(map[string]interface{}, len(p.Search))
		for k, v := range p.Search {
			search[k] = v
		}
		params["search"] = search
	}
	if p.SearchWildcardsEnabled {
		params["searchWildcardsEnabled"] = true
	}
	if p.SortField != nil {
		params["sortfield"] = p.SortField
	}
	if p.SortOrder != nil {
		params["sortorder"] = p.SortOrder
	}

	return params
}

// BuildDeleteParams 返回 dashboard.delete 所需的 dashboardids 列表
func (p DashboardParams) BuildDeleteParams() []string {
	if len(p.DashboardIDs) > 0 {
		return append([]string(nil), p.DashboardIDs...)
	}
	return nil
}
