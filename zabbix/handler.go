/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2025-12-18 20:56:13
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2025-12-18 21:30:15
 * @FilePath: \zabbix-mcp-go\zabbix\handler.go
 * @Description: 文件解释
 * Copyright (c) 2025 by fengzhilaoling@gmail.com, All Rights Reserved.
 */
package zabbix

import (
	"context"
	"sync"
	"zabbixMcp/logger"
)

// APIClient 抽象出最小可用的 Zabbix API 客户端能力
type APIClient interface {
	Call(ctx context.Context, method string, params interface{}, result interface{}) error // 执行一次API调用
	GetDetailedVersionFeatures() map[string]interface{}                                    // 获取详细的版本特性
	AdaptAPIParams(method string, params map[string]interface{}) map[string]interface{}    // 适配API参数
}

// ClientLease 表示一次安全的租借句柄，用于确保归还
type ClientLease interface {
	Client() APIClient
	Release(err error)
}

// ClientProvider 抽象客户端提供方（单实例或连接池）
type ClientProvider interface {
	Acquire(ctx context.Context) (ClientLease, error) // 获取一个客户端句柄
	Info(instanceName string) []ClientInfo            // 获取客户端信息
	Close()                                           // 关闭客户端提供方
}

// NewClientProviderFromConfigs 根据配置创建 ClientProvider
func NewClientProviderFromConfigs(cfgs []ClientConfig) (ClientProvider, error) {
	if len(cfgs) == 0 {
		return nil, nil
	}

	pool := NewClientPool(len(cfgs))
	clients := make([]*ZabbixClient, 0, len(cfgs))
	for _, cfg := range cfgs {
		cli, err := NewZabbixClientFromConfig(cfg)
		if err != nil {
			return nil, err
		}
		clients = append(clients, cli)
		if err := pool.Add(cli); err != nil {
			return nil, err
		}
	}
	prewarmVersions(clients)
	return pool, nil
}

func prewarmVersions(clients []*ZabbixClient) {
	var wg sync.WaitGroup
	for _, c := range clients {
		if c == nil {
			continue
		}
		wg.Add(1)
		go func(cli *ZabbixClient) {
			defer wg.Done()
			if ver, err := NewVersionDetector(cli).DetectVersion(context.Background()); err == nil {
				logger.L().Infof("%s API版本信息: %s", cli.Instance, ver.Full)
			} else {
				logger.L().Errorf("版本检测失败: %v", err)
			}
		}(c)
	}
	wg.Wait()
}
