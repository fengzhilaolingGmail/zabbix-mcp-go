/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2026-01-02 15:51:03
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2026-01-02 16:26:05
 * @FilePath: \zabbix-mcp-go\models\params_item.go
 * @Description: 监控项参数
 * Copyright (c) 2026 by fengzhilaoling@gmail.com, All Rights Reserved.
 */
package models

type ParamsItem struct {
	ItemIDs                      []string                 // itemids：仅返回具有指定ID的监控项
	GroupIDs                     []string                 // groupids：仅返回属于给定组中主机的监控项
	TemplateIDs                  []string                 // templateids：仅返回属于给定模板的监控项
	HostIDs                      []string                 // hostids：仅返回属于指定主机的监控项
	ProxyIDs                     []string                 // proxyids：仅返回由给定proxies监控的监控项
	InterfaceIDs                 []string                 // interfaceids：仅返回使用指定主机接口的监控项
	GraphIDs                     []string                 // graphids：仅返回给定图表中使用的监控项
	TriggerIDs                   []string                 // triggerids：仅返回给定触发器中使用的监控项
	WebItems                     bool                     // webitems：在结果中包含web监控项
	Inherited                    bool                     // inherited：若设置为true则仅返回从模板继承的监控项
	Templated                    bool                     // templated：如果设置为true则仅返回属于模板的监控项
	Monitored                    bool                     // monitored：如果设置为true则仅返回属于受监控主机的已启用监控项
	Group                        string                   // group：仅返回属于指定名称组的监控项
	Host                         string                   // host：仅返回属于指定名称的一个主机的监控项
	EvalType                     int                      // evaltype：标签搜索规则(0-与/或, 2-或)
	Tags                         []map[string]interface{} // tags：仅返回带有指定标签的监控项
	WithTriggers                 bool                     // with_triggers：如果设置为true则仅返回触发器中使用到的监控项
	SelectHosts                  interface{}              // selectHosts：返回hosts属性，包含主机信息
	SelectInterfaces             interface{}              // selectInterfaces：返回interfaces属性，包含主机接口信息
	SelectTriggers               interface{}              // selectTriggers：返回triggers属性，包含触发器信息
	SelectGraphs                 interface{}              // selectGraphs：返回graphs属性，包含图表信息
	SelectDiscoveryData          interface{}              // selectDiscoveryData：返回discoveryData属性，包含监控项发现对象数据
	SelectDiscoveryRule          interface{}              // selectDiscoveryRule：返回discoveryRule属性，包含LLD规则
	SelectDiscoveryRulePrototype interface{}              // selectDiscoveryRulePrototype：返回discoveryRulePrototype属性
	SelectItemDiscovery          interface{}              // selectItemDiscovery：返回itemDiscovery属性(已弃用)
	SelectPreprocessing          interface{}              // selectPreprocessing：返回preprocessing属性，包含预处理选项
	SelectTags                   interface{}              // selectTags：返回tags属性，包含监控项标签
	SelectValueMap               interface{}              // selectValueMap：返回valuemap属性，包含值映射
	Filter                       map[string]interface{}   // filter：仅返回与给定筛选条件完全匹配的结果
	LimitSelects                 int                      // limitSelects：限制子查询返回的记录数量
	SortField                    interface{}              // sortfield：按给定属性对结果进行排序
	CountOutput                  bool                     // countOutput：返回计数而非详细结果
	Editable                     bool                     // editable：仅返回当前用户可编辑的监控项
	ExcludeSearch                bool                     // excludeSearch：对search条件执行排除匹配
	Limit                        int                      // limit：限制返回的监控项数量
	Output                       interface{}              // output：控制输出字段
	PreserveKeys                 bool                     // preservekeys：保持返回结果使用监控项ID作为key
	Search                       map[string]interface{}   // search：按LIKE模式模糊匹配属性
	SearchByAny                  bool                     // searchByAny：search条件之间使用OR
	SearchWildcardsEnabled       bool                     // searchWildcardsEnabled：允许search中的通配符
	SortOrder                    interface{}              // sortorder：排序方向
	StartSearch                  bool                     // startSearch：将search作为前缀匹配
}

// BuildParams 将 ParamsItem 转换为 Zabbix item.get 所需参数
func (p ParamsItem) BuildParams() map[string]interface{} {
	params := map[string]interface{}{}

	if len(p.ItemIDs) > 0 {
		params["itemids"] = append([]string(nil), p.ItemIDs...)
	}
	if len(p.GroupIDs) > 0 {
		params["groupids"] = append([]string(nil), p.GroupIDs...)
	}
	if len(p.TemplateIDs) > 0 {
		params["templateids"] = append([]string(nil), p.TemplateIDs...)
	}
	if len(p.HostIDs) > 0 {
		params["hostids"] = append([]string(nil), p.HostIDs...)
	}
	if len(p.ProxyIDs) > 0 {
		params["proxyids"] = append([]string(nil), p.ProxyIDs...)
	}
	if len(p.InterfaceIDs) > 0 {
		params["interfaceids"] = append([]string(nil), p.InterfaceIDs...)
	}
	if len(p.GraphIDs) > 0 {
		params["graphids"] = append([]string(nil), p.GraphIDs...)
	}
	if len(p.TriggerIDs) > 0 {
		params["triggerids"] = append([]string(nil), p.TriggerIDs...)
	}
	if p.WebItems {
		params["webitems"] = true
	}
	if p.Inherited {
		params["inherited"] = true
	}
	if p.Templated {
		params["templated"] = true
	}
	if p.Monitored {
		params["monitored"] = true
	}
	if p.Group != "" {
		params["group"] = p.Group
	}
	if p.Host != "" {
		params["host"] = p.Host
	}
	if p.EvalType != 0 {
		params["evaltype"] = p.EvalType
	}
	if len(p.Tags) > 0 {
		tags := make([]map[string]interface{}, 0, len(p.Tags))
		for _, tag := range p.Tags {
			if tag == nil {
				continue
			}
			copied := make(map[string]interface{}, len(tag))
			for k, v := range tag {
				copied[k] = v
			}
			tags = append(tags, copied)
		}
		if len(tags) > 0 {
			params["tags"] = tags
		}
	}
	if p.WithTriggers {
		params["with_triggers"] = true
	}
	if p.SelectHosts != nil {
		params["selectHosts"] = p.SelectHosts
	}
	if p.SelectInterfaces != nil {
		params["selectInterfaces"] = p.SelectInterfaces
	}
	if p.SelectTriggers != nil {
		params["selectTriggers"] = p.SelectTriggers
	}
	if p.SelectGraphs != nil {
		params["selectGraphs"] = p.SelectGraphs
	}
	if p.SelectDiscoveryData != nil {
		params["selectDiscoveryData"] = p.SelectDiscoveryData
	}
	if p.SelectDiscoveryRule != nil {
		params["selectDiscoveryRule"] = p.SelectDiscoveryRule
	}
	if p.SelectDiscoveryRulePrototype != nil {
		params["selectDiscoveryRulePrototype"] = p.SelectDiscoveryRulePrototype
	}
	if p.SelectItemDiscovery != nil {
		params["selectItemDiscovery"] = p.SelectItemDiscovery
	}
	if p.SelectPreprocessing != nil {
		params["selectPreprocessing"] = p.SelectPreprocessing
	}
	if p.SelectTags != nil {
		params["selectTags"] = p.SelectTags
	}
	if p.SelectValueMap != nil {
		params["selectValueMap"] = p.SelectValueMap
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

func (p ParamsItem) BuildDeleteParams() []string {
	var itemIDs []string
	for _, id := range p.ItemIDs {
		itemIDs = append(itemIDs, id)
	}
	return itemIDs
}
