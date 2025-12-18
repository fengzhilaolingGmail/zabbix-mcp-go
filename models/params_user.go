package models

// UserGetParams 提供更类型化的 user.get 参数描述
type UserGetParams struct {
	UserIDs []string
	Alias   string
	Output  string
}

// BuildParams 将 UserGetParams 转成 Zabbix API 所需的 map
func (p UserGetParams) BuildParams() map[string]interface{} {
	params := map[string]interface{}{}
	if len(p.UserIDs) > 0 {
		params["userids"] = p.UserIDs
	}
	if p.Alias != "" {
		params["filter"] = map[string]interface{}{"alias": p.Alias}
	}
	if p.Output != "" {
		params["output"] = p.Output
	}
	return params
}
