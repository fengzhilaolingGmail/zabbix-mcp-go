/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2025-12-20 17:18:27
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2025-12-23 19:48:48
 * @FilePath: \zabbix-mcp-go\models\params_user_group.go
 * @Description: 用户组参数
 * @Copyright: Copyright (c) 2025 by fengzhilaoling@gmail.com, All Rights Reserved.
 */
package models

import (
	"fmt"
	"strconv"
)

type UserGroup struct {
	Name                string
	Groupids            []string
	GroupPer            map[int]int
	TemplatePer         map[int]int
	Users               []string
	Status              string // 0:启用 1:禁用
	Output              string // "extend" 等字符串形式
	Filter              map[string]interface{}
	SelectUsers         bool
	SelectRights        bool
	SelectTagFilters    bool
	HostgroupRights     map[string]int
	TemplategroupRights map[string]int
	TagFilters          map[string]string
}

func (p UserGroup) BuildParams() map[string]interface{} {
	params := map[string]interface{}{}
	if p.Name != "" {
		params["name"] = p.Name
	}
	if len(p.GroupPer) > 0 {
		var firstKey int
		for k := range p.GroupPer {
			firstKey = k
			break
		}
		params["rights"] = map[string]interface{}{
			"permission": p.GroupPer[firstKey],
			"id":         firstKey,
		}
		params["hostgroup_rights"] = map[string]interface{}{
			"permission": p.GroupPer[firstKey],
			"id":         firstKey,
		}
	}
	if len(p.Users) > 0 {
		users := make([]map[string]string, 0, len(p.Users))
		for _, v := range p.Users {
			if v == "" {
				continue
			}
			users = append(users, map[string]string{"userid": v})
		}
		if len(users) > 0 {
			fmt.Println(users)
			params["users"] = users
		}
	}
	if len(p.HostgroupRights) > 0 {
		hostgroupRights := make([]map[string]string, 0, len(p.HostgroupRights))
		for k, v := range p.HostgroupRights {
			hostgroupRights = append(hostgroupRights, map[string]string{"permission": strconv.Itoa(v), "id": k})
		}
		params["hostgroup_rights"] = hostgroupRights
	}
	if len(p.TemplategroupRights) > 0 {
		templategroupRights := make([]map[string]string, 0, len(p.TemplategroupRights))
		for k, v := range p.TemplategroupRights {
			templategroupRights = append(templategroupRights, map[string]string{"permission": strconv.Itoa(v), "id": k})
		}
		params["templategroup_rights"] = templategroupRights
	}
	if len(p.TagFilters) > 0 {
		tagFilters := make([]map[string]string, 0, len(p.TagFilters))
		for k, v := range p.TagFilters {
			tagFilters = append(tagFilters, map[string]string{"tag": k, "value": v})
		}
		params["tag_filters"] = tagFilters
	}
	if p.Status != "" {
		params["status"] = p.Status
	}
	if p.Output != "" {
		params["output"] = p.Output
	}
	if p.Filter != nil {
		params["filter"] = p.Filter
	}
	if p.SelectUsers {
		params["selectUsers"] = []string{"userid", "username", "alias", "name", "surname"}
	}
	if p.SelectRights {
		params["selectRights"] = []string{"permission", "id"}
	}
	if p.SelectTagFilters {
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
