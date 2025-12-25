package models

type HostParams struct {
	GroupIDs                      []string                 // groupids：仅返回属于指定主机组的主机
	DServiceIDs                   []string                 // dserviceids：仅返回与指定发现服务相关的主机
	GraphIDs                      []string                 // graphids：仅返回包含指定图形的主机
	HostIDs                       []string                 // hostids：仅返回具有指定主机 ID 的主机
	HTTPTestIDs                   []string                 // httptestids：仅返回包含指定 Web 检查的主机
	InterfaceIDs                  []string                 // interfaceids：仅返回使用指定接口的主机
	ItemIDs                       []string                 // itemids：仅返回包含指定监控项的主机
	MaintenanceIDs                []string                 // maintenanceids：仅返回受指定维护影响的主机
	MonitoredHosts                bool                     // monitored_hosts：仅返回受监控的主机
	ProxyIDs                      []string                 // proxyids：仅返回由指定代理监控的主机
	ProxyGroupIDs                 []string                 // proxy_groupids：仅返回由指定代理组监控的主机
	TemplatedHosts                bool                     // templated_hosts：同时返回主机和模板
	TemplateIDs                   []string                 // templateids：仅返回链接到指定模板的主机
	TriggerIDs                    []string                 // triggerids：仅返回包含指定触发器的主机
	WithItems                     bool                     // with_items：仅返回包含监控项的主机（覆盖 with_monitored_items 和 with_simple_graph_items）
	WithItemPrototypes            bool                     // with_item_prototypes：仅返回包含监控项原型的主机（覆盖 with_simple_graph_item_prototypes）
	WithSimpleGraphItemPrototypes bool                     // with_simple_graph_item_prototypes：仅返回包含数值型监控项原型的主机
	WithGraphs                    bool                     // with_graphs：仅返回包含图形的主机
	WithGraphPrototypes           bool                     // with_graph_prototypes：仅返回包含图形原型的主机
	WithHTTPTests                 bool                     // with_httptests：仅返回包含 Web 检查的主机（覆盖 with_monitored_httptests）
	WithMonitoredHTTPTests        bool                     // with_monitored_httptests：仅返回包含已启用 Web 检查的主机
	WithMonitoredItems            bool                     // with_monitored_items：仅返回包含已启用监控项的主机（覆盖 with_simple_graph_items）
	WithMonitoredTriggers         bool                     // with_monitored_triggers：仅返回包含已启用触发器的主机
	WithSimpleGraphItems          bool                     // with_simple_graph_items：仅返回包含数值型监控项的主机
	WithTriggers                  bool                     // with_triggers：仅返回包含触发器的主机（覆盖 with_monitored_triggers）
	WithProblemsSuppressed        *bool                    // withProblemsSuppressed：控制是否返回具有被抑制问题的主机
	EvalType                      int                      // evaltype：标签搜索规则（0 与/或，2 或）
	Severities                    interface{}              // severities：仅返回包含指定严重性问题的主机，支持单个整数或整数数组
	Tags                          []map[string]interface{} // tags：仅返回具有指定标签的主机
	InheritedTags                 *bool                    // inheritedTags：是否要求链接模板也包含指定标签
	SelectDiscoveries             interface{}              // selectDiscoveries：返回 discoverie 属性，可为 true、"extend"、字段数组或 count
	SelectDiscoveryData           interface{}              // selectDiscoveryData：返回 discoveryData 属性
	SelectDiscoveryRule           interface{}              // selectDiscoveryRule：返回 discoveryRule 属性
	SelectDiscoveryRulePrototype  interface{}              // selectDiscoveryRulePrototype：返回 discoveryRulePrototype 属性
	SelectGraphs                  interface{}              // selectGraphs：返回 graphs 属性
	SelectHostDiscovery           interface{}              // selectHostDiscovery：返回 hostDiscovery 属性（已弃用）
	SelectHostGroups              interface{}              // selectHostGroups：返回 hostgroups 属性
	SelectHTTPTests               interface{}              // selectHttpTests：返回 httpTests 属性
	SelectInterfaces              interface{}              // selectInterfaces：返回 interfaces 属性
	SelectInventory               interface{}              // selectInventory：返回 inventory 属性
	SelectItems                   interface{}              // selectItems：返回 items 属性
	SelectMacros                  interface{}              // selectMacros：返回 macros 属性
	SelectParentTemplates         interface{}              // selectParentTemplates：返回 parentTemplates 属性
	SelectDashboards              interface{}              // selectDashboards：返回 dashboards 属性
	SelectTags                    interface{}              // selectTags：返回 tags 属性
	SelectInheritedTags           interface{}              // selectInheritedTags：返回 inheritedTags 属性
	SelectTriggers                interface{}              // selectTriggers：返回 triggers 属性
	SelectValueMaps               interface{}              // selectValueMaps：返回 valuemaps 属性
	Filter                        map[string]interface{}   // filter：仅返回完全匹配给定筛选条件的主机
	LimitSelects                  int                      // limitSelects：限制子查询返回的记录数
	Search                        map[string]interface{}   // search：按 LIKE 模式模糊匹配属性
	SearchInventory               map[string]interface{}   // searchInventory：按 LIKE 模式匹配清单数据
	SortField                     interface{}              // sortfield：按 hostid/host/name/status 排序，支持字符串或数组
	CountOutput                   bool                     // countOutput：返回计数而非详细结果
	Editable                      bool                     // editable：仅返回当前用户可编辑的主机
	ExcludeSearch                 bool                     // excludeSearch：对 search 条件执行排除匹配
	Limit                         int                      // limit：限制返回的主机数量
	Output                        interface{}              // output：控制输出字段，可为 "extend"、字段数组等
	PreserveKeys                  bool                     // preservekeys：保持返回结果使用主机 ID 作为 key
	SearchByAny                   bool                     // searchByAny：search 条件之间使用 OR
	SearchWildcardsEnabled        bool                     // searchWildcardsEnabled：允许 search 中的通配符
	SortOrder                     interface{}              // sortorder：排序方向，支持字符串或数组
	StartSearch                   bool                     // startSearch：将 search 作为前缀匹配
}

// BuildParams 将 HostParams 转换为 Zabbix host.get 所需参数
func (p HostParams) BuildParams() map[string]interface{} {
	params := map[string]interface{}{}

	if len(p.GroupIDs) > 0 {
		params["groupids"] = append([]string(nil), p.GroupIDs...)
	}
	if len(p.DServiceIDs) > 0 {
		params["dserviceids"] = append([]string(nil), p.DServiceIDs...)
	}
	if len(p.GraphIDs) > 0 {
		params["graphids"] = append([]string(nil), p.GraphIDs...)
	}
	if len(p.HostIDs) > 0 {
		params["hostids"] = append([]string(nil), p.HostIDs...)
	}
	if len(p.HTTPTestIDs) > 0 {
		params["httptestids"] = append([]string(nil), p.HTTPTestIDs...)
	}
	if len(p.InterfaceIDs) > 0 {
		params["interfaceids"] = append([]string(nil), p.InterfaceIDs...)
	}
	if len(p.ItemIDs) > 0 {
		params["itemids"] = append([]string(nil), p.ItemIDs...)
	}
	if len(p.MaintenanceIDs) > 0 {
		params["maintenanceids"] = append([]string(nil), p.MaintenanceIDs...)
	}
	if p.MonitoredHosts {
		params["monitored_hosts"] = true
	}
	if len(p.ProxyIDs) > 0 {
		params["proxyids"] = append([]string(nil), p.ProxyIDs...)
	}
	if len(p.ProxyGroupIDs) > 0 {
		params["proxy_groupids"] = append([]string(nil), p.ProxyGroupIDs...)
	}
	if p.TemplatedHosts {
		params["templated_hosts"] = true
	}
	if len(p.TemplateIDs) > 0 {
		params["templateids"] = append([]string(nil), p.TemplateIDs...)
	}
	if len(p.TriggerIDs) > 0 {
		params["triggerids"] = append([]string(nil), p.TriggerIDs...)
	}
	if p.WithItems {
		params["with_items"] = true
	}
	if p.WithItemPrototypes {
		params["with_item_prototypes"] = true
	}
	if p.WithSimpleGraphItemPrototypes {
		params["with_simple_graph_item_prototypes"] = true
	}
	if p.WithGraphs {
		params["with_graphs"] = true
	}
	if p.WithGraphPrototypes {
		params["with_graph_prototypes"] = true
	}
	if p.WithHTTPTests {
		params["with_httptests"] = true
	}
	if p.WithMonitoredHTTPTests {
		params["with_monitored_httptests"] = true
	}
	if p.WithMonitoredItems {
		params["with_monitored_items"] = true
	}
	if p.WithMonitoredTriggers {
		params["with_monitored_triggers"] = true
	}
	if p.WithSimpleGraphItems {
		params["with_simple_graph_items"] = true
	}
	if p.WithTriggers {
		params["with_triggers"] = true
	}
	if p.WithProblemsSuppressed != nil {
		params["withProblemsSuppressed"] = *p.WithProblemsSuppressed
	}
	if p.EvalType != 0 {
		params["evaltype"] = p.EvalType
	}
	if p.Severities != nil {
		params["severities"] = p.Severities
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
	if p.InheritedTags != nil {
		params["inheritedTags"] = *p.InheritedTags
	}
	if p.SelectDiscoveries != nil {
		params["selectDiscoveries"] = p.SelectDiscoveries
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
	if p.SelectGraphs != nil {
		params["selectGraphs"] = p.SelectGraphs
	}
	if p.SelectHostDiscovery != nil {
		params["selectHostDiscovery"] = p.SelectHostDiscovery
	}
	if p.SelectHostGroups != nil {
		params["selectHostGroups"] = p.SelectHostGroups
	}
	if p.SelectHTTPTests != nil {
		params["selectHttpTests"] = p.SelectHTTPTests
	}
	if p.SelectInterfaces != nil {
		params["selectInterfaces"] = p.SelectInterfaces
	}
	if p.SelectInventory != nil {
		params["selectInventory"] = p.SelectInventory
	}
	if p.SelectItems != nil {
		params["selectItems"] = p.SelectItems
	}
	if p.SelectMacros != nil {
		params["selectMacros"] = p.SelectMacros
	}
	if p.SelectParentTemplates != nil {
		params["selectParentTemplates"] = p.SelectParentTemplates
	}
	if p.SelectDashboards != nil {
		params["selectDashboards"] = p.SelectDashboards
	}
	if p.SelectTags != nil {
		params["selectTags"] = p.SelectTags
	}
	if p.SelectInheritedTags != nil {
		params["selectInheritedTags"] = p.SelectInheritedTags
	}
	if p.SelectTriggers != nil {
		params["selectTriggers"] = p.SelectTriggers
	}
	if p.SelectValueMaps != nil {
		params["selectValueMaps"] = p.SelectValueMaps
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
	if len(p.Search) > 0 {
		search := make(map[string]interface{}, len(p.Search))
		for k, v := range p.Search {
			search[k] = v
		}
		params["search"] = search
	}
	if len(p.SearchInventory) > 0 {
		searchInventory := make(map[string]interface{}, len(p.SearchInventory))
		for k, v := range p.SearchInventory {
			searchInventory[k] = v
		}
		params["searchInventory"] = searchInventory
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

// BuildDeleteParams 返回 host.delete 所需的 hostids 列表
func (p HostParams) BuildDeleteParams() []string {
	if len(p.HostIDs) > 0 {
		return append([]string(nil), p.HostIDs...)
	}
	return nil
}
