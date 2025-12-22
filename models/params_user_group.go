/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2025-12-20 17:18:27
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2025-12-22 14:48:46
 * @FilePath: \zabbix-mcp-go\models\params_user_group.go
 * @Description: 用户组参数
 * @Copyright: Copyright (c) 2025 by fengzhilaoling@gmail.com, All Rights Reserved.
 */
package models

type UserGroup struct {
	Name             string
	Groupids         []string
	GroupPer         map[int]int
	templatePer      map[int]int
	TagFilters       []string
	Users            []string
	Status           int    // 0:启用 1:禁用
	Output           string // "extend" 等字符串形式
	Filter           map[string]interface{}
	SelectUsers      bool
	SelectRights     bool
	SelectTagFilters bool
}

func (P UserGroup) BuildParams() map[string]interface{} {
	params := map[string]interface{}{}
	if P.Name != "" {
		params["name"] = P.Name
	}
	if len(P.GroupPer) > 0 {
		var firstKey int
		for k := range P.GroupPer {
			firstKey = k
			break
		}
		params["rights"] = map[string]interface{}{
			"permission": P.GroupPer[firstKey],
			"id":         firstKey,
		}
		params["hostgroup_rights"] = map[string]interface{}{
			"permission": P.GroupPer[firstKey],
			"id":         firstKey,
		}
	}
	if P.Users != nil {
		params["users"] = P.Users
	}
	params["status"] = P.Status
	if P.Output != "" {
		params["output"] = P.Output
	}
	if P.Filter != nil {
		params["filter"] = P.Filter
	}
	if P.SelectUsers {
		params["selectUsers"] = []string{"userid", "username", "alias", "name", "surname"}
	}
	if P.SelectRights {
		params["selectRights"] = []string{"permission", "id"}
	}
	if P.SelectTagFilters {
		params["selectTagFilters"] = []string{"tag", "value"}
	}
	return params
}

func (p UserGroup) BuildDeleteParams() []string {
	if len(p.Groupids) > 0 {
		return append([]string(nil), p.Groupids...)

	}
	return nil
}
