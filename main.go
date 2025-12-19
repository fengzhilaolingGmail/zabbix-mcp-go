/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2025-12-16 20:14:53
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2025-12-18 17:57:45
 * @FilePath: \zabbix-mcp-go\main.go
 * @Description: 文件解释
 * Copyright (c) 2025 by fengzhilaoling@gmail.com, All Rights Reserved.
 */
package main

import (
	"flag"
	"fmt"
	"zabbixMcp/handler"
	lg "zabbixMcp/logger"
	"zabbixMcp/register"
	zabbix "zabbixMcp/zabbix"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// 定义命令行参数
	var (
		stdioMode = flag.Bool("stdio", false, "使用stdio传输方式")
		httpMode  = flag.Bool("http", false, "使用HTTP/SSE传输方式")
		port      = flag.Int("port", 5443, "HTTP/SSE监听端口")
		level     = flag.String("loglevel", "info", "日志等级 (debug, info, warn, error, panic, fatal)")
	)
	flag.Parse()
	// 初始化日志
	lg.SetLogLevel(*level)
	if err := lg.InitLogger(); err != nil {
		panic("初始化日志失败: " + err.Error())
	}
	defer lg.Sync()

	lg.L().Info("启动Zabbix MCP服务器")
	// 加载配置
	if err := LoadConfig(); err != nil {
		lg.L().Fatalf("加载配置失败: %v", err)
	}

	// 根据配置创建 Zabbix 客户端池（通过接口方式，不直接暴露底层类型）
	poolHandler, err := InitPoolsFromConfig()
	if err != nil {
		lg.L().Fatalf("初始化 Zabbix 客户端池失败: %v", err)
	}
	if poolHandler != nil {
		infos := poolHandler.Info("")
		lg.L().Infof("已初始化 Zabbix 客户端池，容量=%d", len(infos))
		for _, info := range infos {
			lg.L().Infof("客户端: %s 连接方式: %s 登录状态: %v 使用中: %v 版本: %v", info.Instance, info.AuthType, info.Connected, info.InUse, info.Version)
		}
	}

	// 创建MCP服务器
	s := server.NewMCPServer(
		"zabbix-mcp-server",
		"1.0.0",
	)
	lg.L().Info("MCP服务器创建成功")

	// 注册工具
	register.Registers(s)
	lg.L().Info("工具注册完成")

	// 根据参数选择传输方式
	if *stdioMode {
		// 启动stdio服务器
		lg.L().Info("启动stdio传输方式的MCP服务器...")
		if err := server.ServeStdio(s); err != nil {
			lg.L().Fatalf("stdio服务器启动失败: %v", err)
		}
	} else if *httpMode {
		// 启动HTTP/SSE服务器
		startHTTPServer(s, *port)
	} else {
		// 默认同时启动两种方式（在不同的goroutine中）
		lg.L().Info("同时启动stdio和HTTP/SSE传输方式的MCP服务器...")

		// 在后台启动HTTP服务器
		go startHTTPServer(s, *port)

		// 在主线程启动stdio服务器
		if err := server.ServeStdio(s); err != nil {
			lg.L().Fatalf("stdio服务器启动失败: %v", err)
		}
	}
}

// startHTTPServer 启动HTTP传输服务器（使用SSE）
func startHTTPServer(s *server.MCPServer, port int) {
	addr := fmt.Sprintf(":%d", port)
	lg.L().Infof("启动HTTP/SSE传输服务器，监听端口: %d", port)
	lg.L().Infof("MCP端点: http://localhost:%d", port)

	// 使用v0.9.0版本支持的API：创建SSE服务器
	sseServer := server.NewSSEServer(s)
	if err := sseServer.Start(addr); err != nil {
		lg.L().Fatalf("HTTP/SSE服务器启动失败: %v", err)
	}
}

// InitPoolsFromConfig 根据全局 AppConfig 创建并返回一个客户端池，池容量等于实例数量
func InitPoolsFromConfig() (zabbix.ClientProvider, error) {
	n := len(AppConfig.Instances)
	if n == 0 {
		return nil, nil
	}

	cfgs := make([]zabbix.ClientConfig, 0, n)
	for _, inst := range AppConfig.Instances {
		cfgs = append(cfgs, zabbix.ClientConfig{
			Instance: inst.Name,
			URL:      inst.URL,
			User:     inst.User,
			Pass:     inst.Pass,
			Token:    inst.Token,
			AuthType: inst.AuthType,
			Timeout:  30,
			ServerTZ: "",
		})
	}

	// 使用 zabbix 包提供的工厂，返回接口类型，隐藏内部 ClientPool
	handlerObj, err := zabbix.NewClientProviderFromConfigs(cfgs)
	if err != nil {
		return nil, err
	}

	handler.SetClientPool(handlerObj)
	return handlerObj, nil
}
