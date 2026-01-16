/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2025-12-18 22:03:59
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2026-01-16 10:29:12
 * @FilePath: \zabbix-mcp-go\models\params_base.go
 * @Description: 文件解释
 * Copyright (c) 2025 by fengzhilaoling@gmail.com, All Rights Reserved.
 */
package models

// ParamSpec 用于描述一个方法的业务参数，由具体实现转换成 Zabbix API 所需的 map
// 每个具体的 API spec 应当实现 BuildParams 以便统一适配
// MapParams 为旧的 map[string]interface{} 调用方式提供兼容

type ParamSpec interface {
	BuildParams() map[string]interface{}
	BuildDeleteParams() []string
}

// MapParams 允许沿用 map[string]interface{} 的方式，同时实现 ParamSpec 接口
type MapParams map[string]interface{}

// BuildParams 返回 map 的浅拷贝，避免调用方修改底层存储
func (m MapParams) BuildParams() map[string]interface{} {
	cloned := make(map[string]interface{}, len(m))
	for k, v := range m {
		cloned[k] = v
	}
	return cloned
}

func (m MapParams) BuildDeleteParams() []string {
	return []string{}
}

// tag类型
type Tag struct {
	Tag   string `json:"tag"`   // 标签字符串
	Value string `json:"value"` // 标签值字符串
}

type Macro struct {
	GlobalMacroID string `json:"globalmacroid"` // 全局宏的ID
	Macro         string `json:"macro"`         // 宏字符串
	Value         string `json:"value"`         // 宏的值
	Type          int    `json:"type"`          // 宏的类型 0 - (默认) 文本宏, 1 - 密文宏, 2 - 密钥宏
	Description   string `json:"description"`   // 宏描述信息
}

type ZabbixInterface struct {
	InterfaceID  string   `json:"interfaceid"`             // 接口ID.
	Available    int      `json:"available,omitempty"`     // 主机 接口的可用性. 0 - (默认) 未知; 1 - 可用; 2 - 不可用.
	HostID       string   `json:"hostid"`                  // 接口所属的 主机 ID.
	Type         int      `json:"type,omitempty"`          // 接口类型. 1 - Agent; 2 - SNMP; 3 - IPMI; 4 - JMX.
	IP           string   `json:"ip,omitempty"`            // 接口使用的IP地址. 如果通过DNS连接可以为空.
	DNS          string   `json:"dns,omitempty"`           // 接口使用的DNS名称. 如果通过IP连接可以为空.
	Port         string   `json:"port,omitempty"`          // 接口使用的端口号. 可包含用户宏.
	UseIP        int      `json:"useip,omitempty"`         // 是否应通过IP连接. 0 - 使用 主机DNS名称 连接; 1 - 使用 主机IP地址 连接.(默认)
	Main         int      `json:"main,omitempty"`          // 接口是否作为 主机 的默认接口. 每种类型只能有一个接口在 一个主机 上设置为默认. 0 - 非默认; 1 - 默认. (默认)
	Details      []string `json:"details,omitempty"`       // 接口的额外详情 object.
	DisableUntil int      `json:"disable_until,omitempty"` // 不可用 主机 接口的下次轮询时间.
	Error        string   `json:"error,omitempty"`         // 当 主机 接口不可用时的错误文本.
	ErrorsFrom   int      `json:"errors_from,omitempty"`   // 主机 接口变为不可用的时间.
}
