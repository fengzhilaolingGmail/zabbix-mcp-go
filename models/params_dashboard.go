/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2026-01-06 18:50:15
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2026-01-06 18:50:30
 * @FilePath: \zabbix-mcp-go\models\params_dashboard.go
 * @Description: 仪表盘参数
 * @Copyright: Copyright (c) 2025 by fengzhilaoling@gmail.com, All Rights Reserved.
 */
package models

// DashboardWidgetField 仪表盘小部件字段
type DashboardWidgetField struct {
	Type  int         `json:"type"`  // 小部件字段类型 (0-13)
	Name  string      `json:"name"`  // 小部件字段名称
	Value interface{} `json:"value"` // 小部件字段值
}

// DashboardWidget 仪表盘小部件
type DashboardWidget struct {
	WidgetID string                 `json:"widgetid,omitempty"`  // 只读 - 小部件ID
	Type     string                 `json:"type"`                // 必填 - 小部件类型
	Name     string                 `json:"name,omitempty"`      // 自定义小部件名称
	X        int                    `json:"x,omitempty"`         // 水平位置 (0-71)
	Y        int                    `json:"y,omitempty"`         // 垂直位置 (0-63)
	Width    int                    `json:"width,omitempty"`     // 宽度 (1-72)
	Height   int                    `json:"height,omitempty"`    // 高度 (1-64)
	ViewMode int                    `json:"view_mode,omitempty"` // 视图模式 (0-1)
	Fields   []DashboardWidgetField `json:"fields,omitempty"`    // 小部件字段
}

// DashboardPage 仪表盘页面
type DashboardPage struct {
	DashboardPageID string            `json:"dashboard_pageid,omitempty"` // 只读 - 页面ID
	Name            string            `json:"name,omitempty"`             // 页面名称
	DisplayPeriod   int               `json:"display_period,omitempty"`   // 显示周期 (秒)
	Widgets         []DashboardWidget `json:"widgets,omitempty"`          // 小部件数组
}

// DashboardUser 仪表盘用户共享
type DashboardUser struct {
	UserID     string `json:"userid"`     // 用户ID
	Permission int    `json:"permission"` // 权限级别
}

// DashboardUserGroup 仪表盘用户组共享
type DashboardUserGroup struct {
	Usrgrpid   string `json:"usrgrpid"`   // 用户组ID
	Permission int    `json:"permission"` // 权限级别
}

// DashboardParams 仪表盘参数
type DashboardParams struct {
	// 基本字段
	DashboardID string `json:"dashboardid,omitempty"` // 仪表盘ID (用于更新)
	Name        string `json:"name"`                  // 仪表盘名称

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

	// 必填参数 - pages
	if len(p.Pages) > 0 {
		pages := make([]map[string]interface{}, 0, len(p.Pages))
		for _, page := range p.Pages {
			pageMap := map[string]interface{}{
				"name":           page.Name,
				"display_period": page.DisplayPeriod,
			}

			// 处理 widgets
			if len(page.Widgets) > 0 {
				widgets := make([]map[string]interface{}, 0, len(page.Widgets))
				for _, widget := range page.Widgets {
					widgetMap := map[string]interface{}{
						"type":      widget.Type,
						"name":      widget.Name,
						"x":         widget.X,
						"y":         widget.Y,
						"width":     widget.Width,
						"height":    widget.Height,
						"view_mode": widget.ViewMode,
					}

					// 处理 fields
					if len(widget.Fields) > 0 {
						fields := make([]map[string]interface{}, 0, len(widget.Fields))
						for _, field := range widget.Fields {
							fieldMap := map[string]interface{}{
								"type":  field.Type,
								"name":  field.Name,
								"value": field.Value,
							}
							fields = append(fields, fieldMap)
						}
						widgetMap["fields"] = fields
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
