/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2026-01-06 10:00:00
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2026-01-06 10:00:00
 * @FilePath: \zabbix-mcp-go\models\params_graph.go
 * @Description: 图表参数
 * Copyright (c) 2026 by fengzhilaoling@gmail.com, All Rights Reserved.
 */
package models

// GraphGItem 表示图表监控项 (gitem)，参考 Zabbix API: graphitem object
type GraphGItem struct {
	GItemID   string `json:"gitemid,omitempty"`   // 只读 ID
	Color     string `json:"color,omitempty"`     // 颜色 (十六进制)
	ItemID    string `json:"itemid,omitempty"`    // 监控项 ID
	CalcFnc   int    `json:"calc_fnc,omitempty"`  // 计算函数
	DrawType  int    `json:"drawtype,omitempty"`  // 绘制样式
	GraphID   string `json:"graphid,omitempty"`   // 所属图表 ID
	SortOrder int    `json:"sortorder,omitempty"` // 排序序号
	Type      int    `json:"type,omitempty"`      // 类型 (0 简单, 2 总和)
	YaxisSide int    `json:"yaxisside,omitempty"` // Y 轴侧边 (0 左,1 右)
}

// GraphParams 表示用于 graph.get/create/update 的参数（兼容 Zabbix 7.x）
type GraphParams struct {
	// 标识字段
	GraphIDs []string `json:"graphids,omitempty"` // graphids：仅返回指定图表ID的图表
	GraphID  string   `json:"graphid,omitempty"`  // graphid：单个图表ID，用于 graph.update

	// 基本属性
	Name         string  `json:"name,omitempty"`
	Height       int     `json:"height,omitempty"`
	Width        int     `json:"width,omitempty"`
	Flags        int     `json:"flags,omitempty"`
	Graphtype    int     `json:"graphtype,omitempty"`
	PercentLeft  float64 `json:"percent_left,omitempty"`
	PercentRight float64 `json:"percent_right,omitempty"`

	// 显示相关：使用 int (0/1) 与官方文档一致
	Show3d         int `json:"show_3d,omitempty"`
	ShowLegend     int `json:"show_legend,omitempty"`
	ShowWorkPeriod int `json:"show_work_period,omitempty"`
	ShowTriggers   int `json:"show_triggers,omitempty"`

	// 其它属性
	TemplateID string  `json:"templateid,omitempty"`
	YaxisMax   float64 `json:"yaxismax,omitempty"`
	YaxisMin   float64 `json:"yaxismin,omitempty"`
	YmaxItemID string  `json:"ymax_itemid,omitempty"`
	YmaxType   int     `json:"ymax_type,omitempty"`
	YminItemID string  `json:"ymin_itemid,omitempty"`
	YminType   int     `json:"ymin_type,omitempty"`
	UUID       string  `json:"uuid,omitempty"`

	// Create/Update 子对象
	GItems []GraphGItem `json:"gitems,omitempty"` // graph.create/graph.update 时可用

	// 常见通用字段
	SelectItems            interface{}            `json:"selectItems,omitempty"`
	SelectHosts            interface{}            `json:"selectHosts,omitempty"`
	SelectTemplates        interface{}            `json:"selectTemplates,omitempty"`
	Filter                 map[string]interface{} `json:"filter,omitempty"`
	LimitSelects           int                    `json:"limitSelects,omitempty"`
	SortField              interface{}            `json:"sortfield,omitempty"`
	CountOutput            bool                   `json:"countOutput,omitempty"`
	Editable               bool                   `json:"editable,omitempty"`
	ExcludeSearch          bool                   `json:"excludeSearch,omitempty"`
	Limit                  int                    `json:"limit,omitempty"`
	Output                 interface{}            `json:"output,omitempty"`
	PreserveKeys           bool                   `json:"preservekeys,omitempty"`
	Search                 map[string]interface{} `json:"search,omitempty"`
	SearchByAny            bool                   `json:"searchByAny,omitempty"`
	SearchWildcardsEnabled bool                   `json:"searchWildcardsEnabled,omitempty"`
	SortOrder              interface{}            `json:"sortorder,omitempty"`
	StartSearch            bool                   `json:"startSearch,omitempty"`
}

// BuildParams 将 GraphParams 转换为 Zabbix graph.get 所需参数
func (p GraphParams) BuildParams() map[string]interface{} {
	params := map[string]interface{}{}

	if len(p.GraphIDs) > 0 {
		params["graphids"] = append([]string(nil), p.GraphIDs...)
	}
	// update uses single graphid
	if p.GraphID != "" {
		params["graphid"] = p.GraphID
	}
	if p.Name != "" {
		params["name"] = p.Name
	}
	if p.Height > 0 {
		params["height"] = p.Height
	}
	if p.Width > 0 {
		params["width"] = p.Width
	}
	if p.Flags != 0 {
		params["flags"] = p.Flags
	}
	if p.Graphtype != 0 {
		params["graphtype"] = p.Graphtype
	}
	// percents: include even if zero only if explicitly set (we treat zero as omission);
	if p.PercentLeft != 0 {
		params["percent_left"] = p.PercentLeft
	}
	if p.PercentRight != 0 {
		params["percent_right"] = p.PercentRight
	}

	// show_* fields are ints (0/1) in Zabbix API
	if p.Show3d != 0 {
		params["show_3d"] = p.Show3d
	}
	if p.ShowLegend != 0 {
		params["show_legend"] = p.ShowLegend
	}
	if p.ShowWorkPeriod != 0 {
		params["show_work_period"] = p.ShowWorkPeriod
	}
	if p.ShowTriggers != 0 {
		params["show_triggers"] = p.ShowTriggers
	}
	if p.TemplateID != "" {
		params["templateid"] = p.TemplateID
	}
	if p.YaxisMax != 0 {
		params["yaxismax"] = p.YaxisMax
	}
	if p.YaxisMin != 0 {
		params["yaxismin"] = p.YaxisMin
	}
	if p.YmaxItemID != "" {
		params["ymax_itemid"] = p.YmaxItemID
	}
	if p.YmaxType != 0 {
		params["ymax_type"] = p.YmaxType
	}
	if p.YminItemID != "" {
		params["ymin_itemid"] = p.YminItemID
	}
	if p.YminType != 0 {
		params["ymin_type"] = p.YminType
	}
	if p.UUID != "" {
		params["uuid"] = p.UUID
	}

	// gitems 支持在 create/update 操作中传入
	if len(p.GItems) > 0 {
		gitems := make([]map[string]interface{}, 0, len(p.GItems))
		for _, gi := range p.GItems {
			giMap := map[string]interface{}{}
			if gi.GItemID != "" {
				giMap["gitemid"] = gi.GItemID
			}
			if gi.Color != "" {
				giMap["color"] = gi.Color
			}
			if gi.ItemID != "" {
				giMap["itemid"] = gi.ItemID
			}
			if gi.CalcFnc != 0 {
				giMap["calc_fnc"] = gi.CalcFnc
			}
			if gi.DrawType != 0 {
				giMap["drawtype"] = gi.DrawType
			}
			if gi.GraphID != "" {
				giMap["graphid"] = gi.GraphID
			}
			if gi.SortOrder != 0 {
				giMap["sortorder"] = gi.SortOrder
			}
			if gi.Type != 0 {
				giMap["type"] = gi.Type
			}
			if gi.YaxisSide != 0 {
				giMap["yaxisside"] = gi.YaxisSide
			}
			gitems = append(gitems, giMap)
		}
		params["gitems"] = gitems
	}

	if p.SelectItems != nil {
		params["selectItems"] = p.SelectItems
	}
	if p.SelectHosts != nil {
		params["selectHosts"] = p.SelectHosts
	}
	if p.SelectTemplates != nil {
		params["selectTemplates"] = p.SelectTemplates
	}
	if len(p.Filter) > 0 {
		filter := make(map[string]interface{}, len(p.Filter))
		for k, v := range p.Filter {
			filter[k] = v
		}
		params["filter"] = filter
	}
	if p.LimitSelects > 0 {
		params["limitSelects"] = p.LimitSelects
	}
	if p.SortField != nil {
		params["sortfield"] = p.SortField
	}
	if p.CountOutput {
		params["countOutput"] = true
	}
	if p.Editable {
		params["editable"] = true
	}
	if p.ExcludeSearch {
		params["excludeSearch"] = true
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
	if p.SearchByAny {
		params["searchByAny"] = true
	}
	if p.SearchWildcardsEnabled {
		params["searchWildcardsEnabled"] = true
	}
	if p.SortOrder != nil {
		params["sortorder"] = p.SortOrder
	}
	if p.StartSearch {
		params["startSearch"] = true
	}

	return params
}

// BuildDeleteParams 返回 graph.delete 所需的 graphids 列表
func (p GraphParams) BuildDeleteParams() []string {
	if len(p.GraphIDs) > 0 {
		return append([]string(nil), p.GraphIDs...)
	}
	return nil
}
