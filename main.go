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
	"strings"
	"zabbixMcp/handler"
	lg "zabbixMcp/logger"
	"zabbixMcp/register"
	zabbix "zabbixMcp/zabbix"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// 定义命令行参数
	var (
		stdioMode     = flag.Bool("stdio", false, "使用stdio传输方式")
		httpMode      = flag.Bool("http", false, "使用HTTP Streamable传输方式")
		sseMode       = flag.Bool("sse", false, "使用HTTP SSE传输方式")
		port          = flag.Int("port", 5443, "HTTP Streamable监听端口")
		ssePort       = flag.Int("sse-port", 5444, "HTTP SSE监听端口")
		httpTransport = flag.String("http-transport", "streamable", "[兼容参数] HTTP传输类型: streamable 或 sse")
		httpEndpoint  = flag.String("http-endpoint", "/mcp", "streamable-http 端点路径")
		httpStateless = flag.Bool("http-stateless", true, "streamable-http 是否使用无状态会话")
		level         = flag.String("loglevel", "info", "日志等级 (debug, info, warn, error, panic, fatal)")
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

	// 兼容旧参数：-http -http-transport=sse 等价于 -sse
	if *httpMode && strings.EqualFold(strings.TrimSpace(*httpTransport), "sse") {
		*sseMode = true
		*httpMode = false
	}

	// 若未显式指定模式，默认同时启动三种方式
	if !*stdioMode && !*httpMode && !*sseMode {
		*stdioMode = true
		*httpMode = true
		*sseMode = true
	}

	// 端口冲突检查（仅在同时启用 HTTP 与 SSE 时）
	if *httpMode && *sseMode && *port == *ssePort {
		lg.L().Fatalf("HTTP Streamable 端口(%d)与 SSE 端口(%d)冲突，请通过 -port/-sse-port 调整", *port, *ssePort)
	}

	// 启动 HTTP Streamable
	if *httpMode {
		go startStreamableHTTPServer(s, *port, *httpEndpoint, *httpStateless)
	}

	// 启动 SSE
	if *sseMode {
		go startSSEServer(s, *ssePort)
	}

	// 启动 stdio（阻塞）
	if *stdioMode {
		lg.L().Info("启动stdio传输方式的MCP服务器...")
		if err := server.ServeStdio(s); err != nil {
			lg.L().Fatalf("stdio服务器启动失败: %v", err)
		}
	} else {
		// 无 stdio 时阻塞主线程，保持 HTTP/SSE 服务存活
		lg.L().Info("stdio 未启用，服务将以 HTTP/SSE 模式持续运行")
		select {}
	}
}

// startStreamableHTTPServer 启动 HTTP Streamable 传输服务器
func startStreamableHTTPServer(s *server.MCPServer, port int, endpoint string, stateless bool) {
	addr := fmt.Sprintf(":%d", port)
	normalizedEndpoint := strings.TrimSpace(endpoint)
	if normalizedEndpoint == "" {
		normalizedEndpoint = "/mcp"
	}
	if !strings.HasPrefix(normalizedEndpoint, "/") {
		normalizedEndpoint = "/" + normalizedEndpoint
	}

	lg.L().Infof("启动HTTP/Streamable传输服务器，监听端口: %d", port)
	lg.L().Infof("MCP端点: http://localhost:%d%s", port, normalizedEndpoint)
	lg.L().Infof("Streamable会话模式: stateless=%v", stateless)

	httpServer := server.NewStreamableHTTPServer(
		s,
		server.WithEndpointPath(normalizedEndpoint),
		server.WithStateLess(stateless),
	)
	if err := httpServer.Start(addr); err != nil {
		lg.L().Fatalf("HTTP/Streamable服务器启动失败: %v", err)
	}
}

// startSSEServer 启动 HTTP SSE 传输服务器
func startSSEServer(s *server.MCPServer, port int) {
	addr := fmt.Sprintf(":%d", port)
	lg.L().Infof("启动HTTP/SSE传输服务器，监听端口: %d", port)
	lg.L().Infof("SSE端点: http://localhost:%d/sse", port)
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
