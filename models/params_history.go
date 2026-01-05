/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2026-01-04 10:48:38
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2026-01-04 10:53:49
 * @FilePath: \zabbix-mcp-go\models\params_history.go
 * @Description: 历史数据参数结构体
 * @Copyright: Copyright (c) 2025 by fengzhilaoling@gmail.com, All Rights Reserved.
 */

package models

type ParamsHistory struct {
	// 0 - 数值型float;
	// 1 - 字符型;
	// 2 - 日志型;
	// 3 - (默认) 无符号数值型;
	// 4 - 文本型;
	// 5 - 二进制型.
	History                *int                   // history: 要返回的历史object类型 (0..5)。使用指针以允许显式设置 0
	HostIDs                []string               // hostids: 仅返回来自指定主机的历史数据
	ItemIDs                []string               // itemids: 仅返回来自指定监控项的历史数据
	TimeFrom               int                    // time_from: 仅返回在指定时间之后或等于该时间接收的值 (timestamp)
	TimeTill               int                    // time_till: 仅返回在指定时间之前或等于该时间接收的值 (timestamp)
	SortField              interface{}            // sortfield: 按指定属性排序结果，string或[]string，可用: itemid, clock, ns
	CountOutput            bool                   // countOutput: 返回计数而非详细结果
	Editable               bool                   // editable: 仅返回当前用户可编辑的历史数据
	ExcludeSearch          bool                   // excludeSearch: 对search条件执行排除匹配
	Filter                 map[string]interface{} // filter: 过滤条件
	Limit                  int                    // limit: 限制返回的记录数量
	Output                 interface{}            // output: 控制输出字段
	Search                 map[string]interface{} // search: 按LIKE模式模糊匹配属性
	SearchByAny            bool                   // searchByAny: search条件之间使用OR
	SearchWildcardsEnabled bool                   // searchWildcardsEnabled: 允许search中的通配符
	SortOrder              interface{}            // sortorder: 排序方向
	StartSearch            bool                   // startSearch: 将search作为前缀匹配
}

// BuildParams 将 ParamsHistory 转换为 Zabbix history.get 所需参数
func (p ParamsHistory) BuildParams() map[string]interface{} {
	params := map[string]interface{}{}

	if p.History != nil {
		params["history"] = *p.History
	}
	if len(p.HostIDs) > 0 {
		params["hostids"] = append([]string(nil), p.HostIDs...)
	}
	if len(p.ItemIDs) > 0 {
		params["itemids"] = append([]string(nil), p.ItemIDs...)
	}
	if p.TimeFrom != 0 {
		params["time_from"] = p.TimeFrom
	}
	if p.TimeTill != 0 {
		params["time_till"] = p.TimeTill
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

// 仅实现不调用
func (p ParamsHistory) BuildDeleteParams() []string {
	if len(p.HostIDs) > 0 {
		return append([]string(nil), p.HostIDs...)
	}
	return nil
}
