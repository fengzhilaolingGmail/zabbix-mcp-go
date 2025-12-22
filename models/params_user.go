/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2025-12-19 09:10:52
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2025-12-22 14:30:55
 * @FilePath: \zabbix-mcp-go\models\params_user.go
 * @Description: 用户参数
 * @Copyright: Copyright (c) 2025 by fengzhilaoling@gmail.com, All Rights Reserved.
 */

package models

type UserParams struct {
	UserIDs       []string // userids
	MediaTypeIDs  []string // mediatypeids
	SelectUsrgrps []string
	UserName      string
	Name          string
	Passwd        string
	Roleid        string
	UserGroup     string
	Userid        string
	Surname       string
	CurrentPasswd string
	Usrgrps       []string
	Alias         string                 // 兼容旧用法：filter.alias = alias
	Filter        map[string]interface{} // 任意 filter 条件，例如 {"username": []string{"admin"}}
	Search        map[string]interface{} // search 条件，例如 {"username": "ops*"}

	Output       string   // "extend" 等字符串形式
	OutputFields []string // 明确字段列表

	GetAccess bool

	SelectMediasAll     bool
	SelectMediasFields  []string
	SelectRole          bool
	SelectUsrgrpsFields []string

	SortField string
	SortOrder string
	Limit     int
}

func (p UserParams) BuildParams() map[string]interface{} {
	params := map[string]interface{}{}

	if len(p.UserIDs) > 0 {
		params["userids"] = append([]string(nil), p.UserIDs...)
	}
	if len(p.MediaTypeIDs) > 0 {
		params["mediatypeids"] = append([]string(nil), p.MediaTypeIDs...)
	}
	if len(p.SelectUsrgrps) > 0 {
		params["selectUsrgrps"] = p.SelectUsrgrps
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

	if p.SortField != "" {
		params["sortfield"] = p.SortField
	}
	if p.SortOrder != "" {
		params["sortorder"] = p.SortOrder
	}
	if p.Limit > 0 {
		params["limit"] = p.Limit
	}
	if p.UserName != "" {
		params["alias"] = p.UserName
		params["username"] = p.UserName
	}
	if p.Name != "" {
		params["name"] = p.Name
	}
	if p.Surname != "" {
		params["surname"] = p.Surname
	}
	if p.Passwd != "" {
		params["passwd"] = p.Passwd
	}
	if p.Roleid != "" {
		params["roleid"] = p.Roleid
	}
	if p.UserGroup != "" {
		params["usrgrps"] = []map[string]interface{}{{"usrgrpid": p.UserGroup}}
	}
	if p.Userid != "" {
		params["userid"] = p.Userid
	}
	if p.CurrentPasswd != "" {
		params["currentpasswd"] = p.CurrentPasswd
	}

	if len(p.Usrgrps) > 0 {
		var groups []map[string]interface{}
		if existing, ok := params["usrgrps"]; ok {
			if typed, ok := existing.([]map[string]interface{}); ok {
				groups = append(groups, typed...)
			}
		}
		groups = append(groups, make([]map[string]interface{}, 0, len(p.Usrgrps))...)
		for _, v := range p.Usrgrps {
			groups = append(groups, map[string]interface{}{"usrgrpid": v})
		}
		params["usrgrps"] = groups
	}
	return params
}

func (p UserParams) BuildDeleteParams() []string {
	switch {
	case len(p.UserIDs) > 0:
		return append([]string(nil), p.UserIDs...)
	case p.Userid != "":
		return []string{p.Userid}
	default:
		return nil
	}
}
