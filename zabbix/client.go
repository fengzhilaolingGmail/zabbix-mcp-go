/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2025-12-16 20:19:03
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2025-12-16 21:11:20
 * @FilePath: \zabbix-mcp-go\zabbix\client.go
 * @Description: Zabbix客户端相关功能
 * Copyright (c) 2025 by fengzhilaoling@gmail.com, All Rights Reserved.
 */
package zabbix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

type ZabbixClient struct {
	URL        string
	User       string
	Pass       string
	AuthToken  string
	AuthType   string // "password" 或 "token"
	ServerTZ   string
	HTTPClient *http.Client
	mu         sync.Mutex
	// 缓存检测到的版本（防止频繁请求）
	cachedVersion *VersionInfo
	cacheLock     sync.RWMutex
}

// NewZabbixClient 创建新的Zabbix客户端
func NewZabbixClient(url, user, pass string, timeout int) *ZabbixClient {
	var HTTPTimeout time.Duration = 120 * time.Second // 默认超时时间为120秒
	if timeout > 0 {
		HTTPTimeout = time.Duration(timeout) * time.Second
	}
	return &ZabbixClient{
		URL:      url,
		User:     user,
		Pass:     pass,
		ServerTZ: "",
		AuthType: "password", // 默认为密码认证
		HTTPClient: &http.Client{
			Timeout: HTTPTimeout,
		},
	}
}

// SetServerTimezone 设置服务器时区
func (c *ZabbixClient) SetServerTimezone(tz string) {
	if tz == "" {
		tz = time.Local.String()
	}
	c.ServerTZ = tz
}

// 获取客户端缓存的版本（返回拷贝以防外部修改）
func (c *ZabbixClient) GetCachedVersion() *VersionInfo {
	c.cacheLock.RLock()
	defer c.cacheLock.RUnlock()
	if c.cachedVersion == nil {
		return nil
	}
	v := *c.cachedVersion
	return &v
}

// 设置客户端缓存版本（存储拷贝）
func (c *ZabbixClient) SetCachedVersion(v *VersionInfo) {
	c.cacheLock.Lock()
	defer c.cacheLock.Unlock()
	if v == nil {
		c.cachedVersion = nil
		return
	}
	vv := *v
	c.cachedVersion = &vv
}

// 清除缓存的版本信息
func (c *ZabbixClient) ClearCachedVersion() {
	c.cacheLock.Lock()
	defer c.cacheLock.Unlock()
	c.cachedVersion = nil
}

func (c *ZabbixClient) call(method string, params interface{}, auth string) (interface{}, error) {
	// 检测Zabbix版本以确定认证方式
	version, err := NewVersionDetector(c).DetectVersion()
	if err != nil {
		// 如果版本检测失败，使用传统方式
		return c.callWithAuth(method, params, auth)
	}

	// Zabbix 7.0+ 不再使用auth参数，改用HTTP头部认证
	if version.Major >= 7 {
		return c.callWithHeaderAuth(method, params, auth)
	}

	// 旧版本使用传统的auth参数
	return c.callWithAuth(method, params, auth)
}

// 旧版本使用auth参数进行认证
func (c *ZabbixClient) callWithAuth(method string, params interface{}, auth string) (interface{}, error) {
	request := JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
		Auth:    auth,
	}

	requestData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	// 构建完整的API URL，如果URL中没有包含api_jsonrpc.php则自动添加
	if !strings.Contains(c.URL, "api_jsonrpc.php") {
		// 移除末尾的斜杠（如果有）
		c.URL = strings.TrimRight(c.URL, "/")
		// 添加API路径
		c.URL = c.URL + "/api_jsonrpc.php"
	}

	resp, err := c.HTTPClient.Post(c.URL, "application/json", bytes.NewBuffer(requestData))
	if err != nil {
		return nil, fmt.Errorf("HTTP请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var response JSONRPCResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if response.Error != nil {
		return nil, response.Error
	}

	return response.Result, nil
}

// 使用HTTP头部认证
func (c *ZabbixClient) callWithHeaderAuth(method string, params interface{}, auth string) (interface{}, error) {
	// Zabbix 7.0+ 不在请求体中包含auth参数
	request := JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
		// Auth字段为空，不包含在JSON中
	}

	requestData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	// 构建完整的API URL，如果URL中没有包含api_jsonrpc.php则自动添加
	if !strings.Contains(c.URL, "api_jsonrpc.php") {
		// 移除末尾的斜杠（如果有）
		c.URL = strings.TrimRight(c.URL, "/")
		// 添加API路径
		c.URL = c.URL + "/api_jsonrpc.php"
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", c.URL, bytes.NewBuffer(requestData))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// Zabbix 7.0+ 使用Authorization头部进行认证
	if auth != "" {
		req.Header.Set("Authorization", "Bearer "+auth)
	}

	// 执行请求
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var response JSONRPCResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if response.Error != nil {
		return nil, response.Error
	}

	return response.Result, nil
}

func (c *ZabbixClient) Call(method string, params interface{}) (interface{}, error) {
	c.mu.Lock()
	authToken := c.AuthToken
	c.mu.Unlock()

	if authToken == "" && c.AuthType != "token" {
		if err := c.Login(); err != nil {
			return nil, err
		}
		c.mu.Lock()
		authToken = c.AuthToken
		c.mu.Unlock()
	}

	result, err := c.call(method, params, authToken)
	if err != nil {
		// 如果认证失败，尝试重新登录（仅对密码认证）
		if rpcErr, ok := err.(*RPCError); ok && rpcErr.Code == -32602 && c.AuthType != "token" {
			if err := c.Login(); err != nil {
				return nil, err
			}
			c.mu.Lock()
			authToken = c.AuthToken
			c.mu.Unlock()
			return c.call(method, params, authToken)
		}
		return nil, err
	}

	return result, nil
}
