/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2025-12-18 11:28:10
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2025-12-18 11:28:10
 * @FilePath: \zabbix-mcp-go\test\handler\common_test.go
 * @Description: 文件详情
 * @Copyright: Copyright (c) 2025 by fengzhilaoling@gmail.com, All Rights Reserved.
 */
package handler_test

import (
	"errors"
	"zabbixMcp/zabbix"
)

// mockHandler 用于测试，简化实现 zabbix.ZabbixClientHandler 接口
type mockHandler struct {
	callFn func(method string, params interface{}) (interface{}, error)
	info   []zabbix.ClientInfo
}

func (m *mockHandler) SetServerTimezone(tz string)            {}
func (m *mockHandler) GetCachedVersion() *zabbix.VersionInfo  { return nil }
func (m *mockHandler) SetCachedVersion(v *zabbix.VersionInfo) {}
func (m *mockHandler) ClearCachedVersion()                    {}
func (m *mockHandler) Info() []zabbix.ClientInfo              { return m.info }
func (m *mockHandler) Call(method string, params interface{}) (interface{}, error) {
	if m.callFn == nil {
		return nil, errors.New("no callFn set")
	}
	return m.callFn(method, params)
}
