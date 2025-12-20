package models

// UserGetParams 提供更类型化的 user.get 参数描述
//
// 设计目标：在保持常见字段（userids/alias/output）简单易用的同时，兼容
// Zabbix user.get 的更多能力，例如访问权限、分组/媒介筛选、分页排序等。
// 结构中大部分字段都是可选的，对应 API 中的等名参数。
type UserGetParams struct {
	UserIDs      []string               // userids
	UserGroupIDs []string               // usrgrpids
	MediaTypeIDs []string               // mediatypeids
	Alias        string                 // 兼容旧用法：filter.alias = alias
	Filter       map[string]interface{} // 任意 filter 条件，例如 {"username": []string{"admin"}}
	Search       map[string]interface{} // search 条件，例如 {"username": "ops*"}

	Output       string   // "extend" 等字符串形式
	OutputFields []string // 明确字段列表

	GetAccess bool

	SelectMediasAll     bool
	SelectMediasFields  []string
	SelectRole          bool
	SelectUsrgrpsAll    bool
	SelectUsrgrpsFields []string

	SortField string
	SortOrder string
	Limit     int
}

// BuildParams 将 UserGetParams 转成 Zabbix API 所需的 map
func (p UserGetParams) BuildParams() map[string]interface{} {
	params := map[string]interface{}{}

	if len(p.UserIDs) > 0 {
		params["userids"] = append([]string(nil), p.UserIDs...)
	}
	if len(p.UserGroupIDs) > 0 {
		params["usrgrpids"] = append([]string(nil), p.UserGroupIDs...)
	}
	if len(p.MediaTypeIDs) > 0 {
		params["mediatypeids"] = append([]string(nil), p.MediaTypeIDs...)
	}

	filter := map[string]interface{}{}
	if len(p.Filter) > 0 {
		for k, v := range p.Filter {
			filter[k] = v
		}
	}
	if p.Alias != "" {
		filter["alias"] = p.Alias
	}
	if len(filter) > 0 {
		params["filter"] = filter
	}

	if len(p.Search) > 0 {
		search := make(map[string]interface{}, len(p.Search))
		for k, v := range p.Search {
			search[k] = v
		}
		params["search"] = search
	}

	if len(p.OutputFields) > 0 {
		params["output"] = append([]string(nil), p.OutputFields...)
	} else if p.Output != "" {
		params["output"] = p.Output
	}

	if p.GetAccess {
		params["getAccess"] = true
	}

	if len(p.SelectMediasFields) > 0 {
		params["selectMedias"] = append([]string(nil), p.SelectMediasFields...)
	} else if p.SelectMediasAll {
		params["selectMedias"] = true
	}

	if p.SelectRole {
		params["selectRole"] = true
	}

	if len(p.SelectUsrgrpsFields) > 0 {
		params["selectUsrgrps"] = append([]string(nil), p.SelectUsrgrpsFields...)
	} else if p.SelectUsrgrpsAll {
		params["selectUsrgrps"] = true
	}

	if p.SortField != "" {
		params["sortfield"] = p.SortField
	}
	if p.SortOrder != "" {
		params["sortorder"] = p.SortOrder
	}
	if p.Limit > 0 {
		params["limit"] = p.Limit
	}

	return params
}

type UserCreateParams struct {
	UserName  string
	Name      string
	Passwd    string
	Roleid    string
	UserGroup string
}

func (P UserCreateParams) BuildParams() map[string]interface{} {
	params := map[string]interface{}{}
	if P.UserName != "" {
		params["alias"] = P.UserName
		params["username"] = P.UserName
	}
	if P.Name != "" {
		params["name"] = P.Name
	}
	if P.Passwd != "" {
		params["passwd"] = P.Passwd
	}
	if P.Roleid != "" {
		params["roleid"] = P.Roleid
	}
	if P.UserGroup != "" {
		params["usrgrps"] = []map[string]interface{}{{"usrgrpid": P.UserGroup}}
	}
	return params
}
