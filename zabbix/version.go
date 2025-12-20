/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2025-12-16 20:54:52
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2025-12-20 13:16:30
 * @FilePath: \zabbix-mcp-go\zabbix\version.go
 * @Description: 版本检测相关功能
 * Copyright (c) 2025 by fengzhilaoling@gmail.com, All Rights Reserved.
 */
package zabbix

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"zabbixMcp/logger"
	"zabbixMcp/models"
)

// VersionInfo Zabbix版本信息
type VersionInfo struct {
	Major int    // 主版本号
	Minor int    // 次版本号
	Patch int    // 补丁版本号解析失败时为0
	Full  string // 完整版本字符串
}

// VersionDetector 版本检测器
type VersionDetector struct {
	client *ZabbixClient
}

// NewVersionDetector 创建版本检测器
func NewVersionDetector(client *ZabbixClient) *VersionDetector {
	return &VersionDetector{client: client}
}

// DetectVersion 检测Zabbix版本
func (vd *VersionDetector) DetectVersion(ctx context.Context) (*VersionInfo, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	// 先尝试从 client 缓存读取
	if vd.client != nil {
		if cached := vd.client.GetCachedVersion(); cached != nil {
			return cached, nil
		}
	}

	// 获取API版本信息 - 使用内部调用避免循环依赖
	// Zabbix API要求params为空数组[]而不是nil
	result, err := vd.client.callWithAuth(ctx, "apiinfo.version", []interface{}{}, "") // 使用旧方式获取API版本信息
	if err != nil {
		logger.L().Warnf("获取API版本失败: %v, 尝试新方法", err)
		result, err = vd.client.callWithHeaderAuth(ctx, "apiinfo.version", nil, "") // 使用新方式获取API版本信息
		if err != nil {
			return nil, fmt.Errorf("获取API版本失败: %w", err)
		}
	}
	var apiVersion string
	if err := json.Unmarshal(result, &apiVersion); err != nil {
		return nil, fmt.Errorf("API版本响应格式错误: %w", err)
	}

	// 解析版本号
	version, err := vd.parseVersion(apiVersion)
	if err != nil {
		return nil, fmt.Errorf("解析版本号失败: %w", err)
	}
	version.Full = apiVersion

	// 将结果写入 client 缓存
	if vd.client != nil {
		vd.client.SetCachedVersion(version)
	}

	return version, nil
}

// parseVersion 解析版本字符串
func (vd *VersionDetector) parseVersion(versionStr string) (*VersionInfo, error) {
	// 移除前缀
	versionStr = strings.TrimPrefix(versionStr, "v")

	// 分割版本号
	parts := strings.Split(versionStr, ".")
	if len(parts) < 2 {
		return nil, fmt.Errorf("版本格式不正确: %s", versionStr)
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("解析主版本号失败: %w", err)
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("解析次版本号失败: %w", err)
	}

	patch := 0
	if len(parts) > 2 {
		patch, err = strconv.Atoi(parts[2])
		if err != nil {
			patch = 0 // 如果解析失败，默认为0
		}
	}

	return &VersionInfo{
		Major: major,
		Minor: minor,
		Patch: patch,
		Full:  versionStr,
	}, nil
}

// ParseVersion 对 parseVersion 的导出包装，便于测试
func (vd *VersionDetector) ParseVersion(versionStr string) (*VersionInfo, error) {
	return vd.parseVersion(versionStr)
}

// getDefaultFeatures 获取默认功能集
func (vd *VersionDetector) getDefaultFeatures() map[string]bool {
	return map[string]bool{
		"host_management":      true,
		"item_management":      true,
		"trigger_management":   true,
		"template_management":  true,
		"event_acknowledgment": true,
	}
}

// 在 version.go 中添加更详细的版本特性映射
func (vd *VersionDetector) GetDetailedVersionFeatures() map[string]interface{} {
	version, err := vd.DetectVersion(context.Background())
	if err != nil {
		// 将 map[string]bool 转换为 map[string]interface{}
		defaultFeatures := vd.getDefaultFeatures()
		result := make(map[string]interface{})
		for k, v := range defaultFeatures {
			result[k] = v
		}
		return result
	}

	features := make(map[string]interface{})

	// API端点支持
	features["endpoints"] = map[string]bool{
		"problem.get":    version.Major >= 4,
		"sla.get":        version.Major >= 5,
		"authentication": version.Major >= 7,
		"connector":      version.Major >= 6,
		"proxygroup":     version.Major >= 7,
	}

	// 参数支持
	features["parameters"] = map[string]bool{
		"selectTags":          version.Major >= 4,
		"selectDependencies":  version.Major >= 4,
		"selectPreprocessing": version.Major >= 4,
		"templateSelectTags":  version.Major >= 5, // template.get
	}

	return features
}

// AdaptAPIParams 根据版本适配API参数
func (vd *VersionDetector) AdaptAPIParams(method string, spec models.ParamSpec) map[string]interface{} {
	version, err := vd.DetectVersion(context.Background())
	logger.L().Info(version.Full)
	var params map[string]interface{}
	if spec != nil {
		params = spec.BuildParams()
	} else {
		params = map[string]interface{}{}
	}

	if err != nil {
		return params
	}

	adaptedParams := make(map[string]interface{}, len(params))
	for k, v := range params {
		adaptedParams[k] = v
	}

	// 根据版本调整参数
	switch method {
	case "host.get":
		if version.Major < 4 {
			// 旧版本不支持某些参数
			delete(adaptedParams, "selectTags")
			// adaptedParams["output"] = []string{"hostid", "name"}
		}
	case "user.get":
		if version.Major > 5 {
			if f, ok := adaptedParams["filter"].(map[string]interface{}); ok {
				delete(f, "alias")
				if len(f) == 0 {
					delete(adaptedParams, "filter")
				}
			}
		} else {
			if f, ok := adaptedParams["filter"].(map[string]interface{}); ok {
				delete(f, "username")
				if len(f) == 0 {
					delete(adaptedParams, "filter")
				}
			}
		}
	case "item.get":
		if version.Major < 4 {
			delete(adaptedParams, "selectTags")
			delete(adaptedParams, "selectPreprocessing")
		}
	case "trigger.get":
		if version.Major < 4 {
			delete(adaptedParams, "selectTags")
			delete(adaptedParams, "selectDependencies")
		}
	case "template.get":
		if version.Major < 5 {
			delete(adaptedParams, "selectTags")
		}
	case "user.create":
		// Zabbix 5.x 及更早版本不使用 `username` 字段，使用 `alias`。
		// 对于 6.x/7.x，使用 `username`。
		if version.Major <= 5 {
			if un, ok := adaptedParams["username"]; ok {
				if s, ok2 := un.(string); ok2 {
					adaptedParams["alias"] = s
				}
				delete(adaptedParams, "username")
			}
		} else {
			// >=6: 尽量保留 username；如果只有 alias 提供也不会出错
		}
	}

	return adaptedParams
}
