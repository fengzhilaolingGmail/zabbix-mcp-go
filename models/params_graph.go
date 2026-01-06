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

// GraphParams 表示用于 graph.get/create/update 的参数
type GraphParams struct {
	GraphIDs       []string // graphids：仅返回指定图表ID的图表
	GraphID        string   // graphid：单个图表ID，用于 graph.update
	Name           string   // name：图表名称
	Height         int      // height：图表高度
	Width          int      // width：图表宽度
	Flags          int      // flags：图表来源（0 普通，4 发现），只读
	Graphtype      int      // graphtype：图表类型（0 普通，1 堆叠，2 饼图，3 爆炸）
	PercentLeft    float64  // percent_left：左侧百分比
	PercentRight   float64  // percent_right：右侧百分比
	Show3d         bool     // show_3d：是否以3D显示饼图/爆炸图
	ShowLegend     bool     // show_legend：是否显示图例
	ShowWorkPeriod bool     // show_work_period：是否显示工作时间
	ShowTriggers   bool     // show_triggers：是否显示触发器线
	TemplateID     string   // templateid：父模板图表ID，只读
	YaxisMax       float64  // yaxismax：Y轴最大值（固定）
	YaxisMin       float64  // yaxismin：Y轴最小值（固定）
	YmaxItemID     string   // ymax_itemid：作为Y轴最大值的监控项ID
	YmaxType       int      // ymax_type：Y轴最大值计算方法（0 计算，1 固定，2 监控项）
	YminItemID     string   // ymin_itemid：作为Y轴最小值的监控项ID
	YminType       int      // ymin_type：Y轴最小值计算方法（0 计算，1 固定，2 监控项）
	UUID           string   // uuid：通用唯一标识符

	// 常见通用字段
	SelectItems            interface{}            // selectItems：返回 items 属性
	SelectHosts            interface{}            // selectHosts：返回 hosts 属性
	SelectTemplates        interface{}            // selectTemplates：返回 templates 属性
	Filter                 map[string]interface{} // filter：仅返回完全匹配给定筛选条件的图表
	LimitSelects           int                    // limitSelects：限制子查询返回的记录数量
	SortField              interface{}            // sortfield：按给定属性对结果进行排序
	CountOutput            bool                   // countOutput：返回计数而非详细结果
	Editable               bool                   // editable：仅返回当前用户可编辑的图表
	ExcludeSearch          bool                   // excludeSearch：对search条件执行排除匹配
	Limit                  int                    // limit：限制返回的图表数量
	Output                 interface{}            // output：控制输出字段
	PreserveKeys           bool                   // preservekeys：保持返回结果使用图表ID作为key
	Search                 map[string]interface{} // search：按LIKE模式模糊匹配属性
	SearchByAny            bool                   // searchByAny：search条件之间使用OR
	SearchWildcardsEnabled bool                   // searchWildcardsEnabled：允许search中的通配符
	SortOrder              interface{}            // sortorder：排序方向
	StartSearch            bool                   // startSearch：将search作为前缀匹配
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
	if p.Show3d {
		params["show_3d"] = true
	}
	if p.ShowLegend {
		params["show_legend"] = true
	}
	if p.ShowWorkPeriod {
		params["show_work_period"] = true
	}
	if p.ShowTriggers {
		params["show_triggers"] = true
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
