/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2025-12-16 20:43:12
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2025-12-16 20:43:25
 * @FilePath: \zabbix-mcp-go\zabbix\auth.go
 * @Description: 登录和认证相关功能
 * Copyright (c) 2025 by fengzhilaoling@gmail.com, All Rights Reserved.
 */
package zabbix

import (
	"context"
	"encoding/json"
	"fmt"
)

// Login 登录Zabbix API
func (c *ZabbixClient) Login(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 如果已经设置了token认证，直接验证token有效性
	if c.AuthType == "token" && c.AuthToken != "" {
		// 保存当前的AuthToken，临时清空它以调用不需要认证的API
		savedToken := c.AuthToken
		c.AuthToken = ""

		// 尝试调用apiinfo.version来验证连接是否正常
		// apiinfo.version不需要认证，所以传入空auth参数
		_, err := c.call(ctx, "apiinfo.version", map[string]interface{}{}, "")

		// 恢复AuthToken
		c.AuthToken = savedToken

		if err != nil {
			return fmt.Errorf("token认证失败: %w", err)
		}
		return nil
	}

	// 密码认证
	params := map[string]string{
		"user":     c.User,
		"password": c.Pass,
	}

	// 使用内部调用，传入空auth进行登录
	response, err := c.callWithAuth(ctx, "user.login", params, "")
	if err != nil {
		return fmt.Errorf("登录失败: %w", err)
	}

	var authToken string
	if err := json.Unmarshal(response, &authToken); err != nil {
		return fmt.Errorf("解析登录响应失败: %w", err)
	}

	c.AuthToken = authToken
	return nil
}

// Logout 登出Zabbix API
func (c *ZabbixClient) Logout(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.AuthToken == "" {
		return nil
	}

	_, err := c.call(ctx, "user.logout", nil, c.AuthToken)
	c.AuthToken = ""
	return err
}

// SetAuthToken 设置认证token
func (c *ZabbixClient) SetAuthToken(token string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.AuthToken = token
	c.AuthType = "token"
}

// SetAuthType 设置认证方式
func (c *ZabbixClient) SetAuthType(authType string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.AuthType = authType
}

// GetAuthToken 获取当前认证token
func (c *ZabbixClient) GetAuthToken() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.AuthToken
}

// GetAuthType 获取当前认证方式
func (c *ZabbixClient) GetAuthType() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.AuthType
}
