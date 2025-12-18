/*
  - @Author: fengzhilaoling fengzhilaoling@gmail.com
  - @Date: 2025-12-16 20:19:03
    }
  - @LastEditTime: 2025-12-18 21:39:06
  - @FilePath: \zabbix-mcp-go\zabbix\client.go
  - @Description: Zabbix客户端相关功能
  - Copyright (c) 2025 by fengzhilaoling@gmail.com, All Rights Reserved.
*/
package zabbix

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"zabbixMcp/logger"
	"zabbixMcp/models"
)

type ZabbixClient struct {
	Instance         string
	URL              string
	apiURL           string
	User             string
	Pass             string
	AuthToken        string
	AuthType         string // "password" 或 "token"
	ServerTZ         string
	HTTPClient       *http.Client
	mu               sync.Mutex
	preferHeaderAuth bool
	// 缓存检测到的版本（防止频繁请求）
	cachedVersion *VersionInfo
	cacheLock     sync.RWMutex
}

// NewZabbixClient 创建新的Zabbix客户端
func NewZabbixClient(name, rawURL, user, pass string, timeout int) (*ZabbixClient, error) {
	var HTTPTimeout time.Duration = 120 * time.Second // 默认超时时间为120秒
	if timeout > 0 {
		HTTPTimeout = time.Duration(timeout) * time.Second
	}
	apiURL, err := buildAPIEndpoint(rawURL)
	if err != nil {
		return nil, err
	}
	return &ZabbixClient{
		Instance: name,
		URL:      rawURL,
		apiURL:   apiURL,
		User:     user,
		Pass:     pass,
		ServerTZ: "",
		AuthType: "password", // 默认为密码认证
		HTTPClient: &http.Client{
			Timeout: HTTPTimeout,
		},
	}, nil
}

// ClientConfig 是用于创建 ZabbixClient 的工厂配置结构体，便于在一处集中管理实例化逻辑
type ClientConfig struct {
	Instance string
	URL      string
	User     string
	Pass     string
	Token    string
	AuthType string // "password" 或 "token"
	Timeout  int    // HTTP 超时（秒），0 表示使用默认值
	ServerTZ string // 可选，设置服务器时区，空则保持默认
}

// NewZabbixClientFromConfig 根据 ClientConfig 创建并初始化一个 *ZabbixClient。
// 这样可以把实例化逻辑集中到工厂里，调用方（例如 main）只需传入配置即可；同时便于测试替换。
func NewZabbixClientFromConfig(cfg ClientConfig) (*ZabbixClient, error) {
	cli, err := NewZabbixClient(cfg.Instance, cfg.URL, cfg.User, cfg.Pass, cfg.Timeout)
	if err != nil {
		return nil, err
	}
	if cfg.AuthType != "" {
		cli.SetAuthType(cfg.AuthType)
	}
	if cfg.Token != "" {
		cli.SetAuthToken(cfg.Token)
	}
	// 时区使用配置中的值，如果为空则使用本地时区
	cli.SetServerTimezone(cfg.ServerTZ)
	if err := cli.Login(context.Background()); err != nil {
		return nil, err
	}
	if ver, err := NewVersionDetector(cli).DetectVersion(context.Background()); err == nil {
		cli.preferHeaderAuth = ver.Major >= 7
	} else {
		// 版本探测失败不阻塞流程，但记录日志由调用方负责
		logger.L().Warnf("探测 %s 版本失败: %v", cli.Instance, err)
	}
	return cli, nil
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

func (c *ZabbixClient) Call(ctx context.Context, method string, params interface{}, result interface{}) error {
	authToken, err := c.ensureAuthToken(ctx)
	if err != nil {
		return err
	}

	payload, err := c.call(ctx, method, params, authToken)
	if err != nil {
		if isAuthError(err) && c.getAuthType() != "token" {
			if err := c.Login(ctx); err != nil {
				return err
			}
			if payload, err = c.call(ctx, method, params, c.getAuthToken()); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	if result == nil {
		return nil
	}
	return json.Unmarshal(payload, result)
}

func (c *ZabbixClient) ensureAuthToken(ctx context.Context) (string, error) {
	token := c.getAuthToken()
	if token == "" && c.getAuthType() != "token" {
		if err := c.Login(ctx); err != nil {
			return "", err
		}
		token = c.getAuthToken()
	}
	return token, nil
}

func (c *ZabbixClient) call(ctx context.Context, method string, params interface{}, auth string) (json.RawMessage, error) {
	primaryHeader := c.prefersHeaderAuth()
	var first func(context.Context, string, interface{}, string) (json.RawMessage, error)
	var second func(context.Context, string, interface{}, string) (json.RawMessage, error)
	if primaryHeader {
		first = c.callWithHeaderAuth
		second = c.callWithAuth
	} else {
		first = c.callWithAuth
		second = c.callWithHeaderAuth
	}

	if payload, err := first(ctx, method, params, auth); err == nil {
		return payload, nil
	} else if _, ok := err.(*models.RPCError); ok && second != nil {
		if altPayload, altErr := second(ctx, method, params, auth); altErr == nil {
			c.setHeaderPreference(!primaryHeader)
			return altPayload, nil
		}
		return nil, err
	} else {
		return nil, err
	}
}

func (c *ZabbixClient) callWithAuth(ctx context.Context, method string, params interface{}, auth string) (json.RawMessage, error) {
	request := models.JSONRPCRequest{
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

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.apiURL, bytes.NewReader(requestData))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	return c.doRequest(req)
}

func (c *ZabbixClient) callWithHeaderAuth(ctx context.Context, method string, params interface{}, auth string) (json.RawMessage, error) {
	request := models.JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
	}

	requestData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.apiURL, bytes.NewReader(requestData))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if auth != "" {
		req.Header.Set("Authorization", "Bearer "+auth)
	}

	return c.doRequest(req)
}

func (c *ZabbixClient) doRequest(req *http.Request) (json.RawMessage, error) {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var response models.JSONRPCResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if response.Error != nil {
		return nil, response.Error
	}
	if response.Result == nil {
		return json.RawMessage("null"), nil
	}
	return response.Result, nil
}

func (c *ZabbixClient) prefersHeaderAuth() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.preferHeaderAuth
}

func (c *ZabbixClient) setHeaderPreference(header bool) {
	c.mu.Lock()
	c.preferHeaderAuth = header
	c.mu.Unlock()
}

func (c *ZabbixClient) getAuthToken() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.AuthToken
}

func (c *ZabbixClient) getAuthType() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.AuthType
}

func isAuthError(err error) bool {
	if err == nil {
		return false
	}
	if rpcErr, ok := err.(*models.RPCError); ok {
		return rpcErr.Code == -32602 || rpcErr.Code == -32500
	}
	return false
}

func (c *ZabbixClient) GetDetailedVersionFeatures() map[string]interface{} {
	return NewVersionDetector(c).GetDetailedVersionFeatures()
}

func (c *ZabbixClient) AdaptAPIParams(method string, spec models.ParamSpec) map[string]interface{} {
	return NewVersionDetector(c).AdaptAPIParams(method, spec)
}

func buildAPIEndpoint(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", errors.New("zabbix url 不能为空")
	}
	if strings.Contains(trimmed, "api_jsonrpc.php") {
		return trimmed, nil
	}
	return strings.TrimRight(trimmed, "/") + "/api_jsonrpc.php", nil
}
