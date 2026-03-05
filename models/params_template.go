package models

type ParamsTemplate struct {
	TemplateIDs       []string `json:"templateids"`       // 模板ID列表
	GroupIDs          []string `json:"groupids"`          // 组ID列表
	ParentTemplateIDs []string `json:"parenttemplateids"` // 父模板ID列表
	HostIDs           []string `json:"hostids"`           // 主机ID列表
	GraphIDs          []string `json:"graphids"`          // 图表ID列表
	ItemIDs           []string `json:"itemids"`           // 监控项ID列表
	TriggerIDs        []string `json:"triggerids"`        // 触发器ID列表
}

func (p ParamsTemplate) BuildParams() map[string]interface{} {
	return nil
}

func (p ParamsTemplate) BuildQueryParams() map[string]interface{} {
	return nil
}

func (p ParamsTemplate) BuildDeleteParams() []string {
	var templateids []string
	for _, id := range p.TemplateIDs {
		templateids = append(templateids, id)
	}
	return templateids
}

// https://www.zabbix.com/documentation/7.4/zh/manual/api/reference/template/get
