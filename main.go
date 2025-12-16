package main

import (
	"flag"
	"fmt"
	"zabbixMcp/logger"

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
	logger.SetLogLevel(*level)
	if err := logger.InitLogger(); err != nil {
		panic("初始化日志失败: " + err.Error())
	}
	defer logger.Sync()

	logger.L().Info("启动Zabbix MCP服务器")
	// 加载配置
	if err := LoadConfig(); err != nil {
		logger.L().Fatalf("加载配置失败: %v", err)
	}

	// for _, instance := range AppConfig.Instances {
	// 	// 尝试连接Zabbix实例
	// 	client := zabbix.NewZabbixClient(instance.URL, instance.User, instance.Pass)
	// }

	// 创建MCP服务器
	s := server.NewMCPServer(
		"zabbix-mcp-server",
		"1.0.0",
	)
	logger.L().Info("MCP服务器创建成功")

	// 注册工具
	Registers(s)
	logger.L().Info("工具注册完成")

	// 根据参数选择传输方式
	if *stdioMode {
		// 启动stdio服务器
		logger.L().Info("启动stdio传输方式的MCP服务器...")
		if err := server.ServeStdio(s); err != nil {
			logger.L().Fatalf("stdio服务器启动失败: %v", err)
		}
	} else if *httpMode {
		// 启动HTTP/SSE服务器
		startHTTPServer(s, *port)
	} else {
		// 默认同时启动两种方式（在不同的goroutine中）
		logger.L().Info("同时启动stdio和HTTP/SSE传输方式的MCP服务器...")

		// 在后台启动HTTP服务器
		go startHTTPServer(s, *port)

		// 在主线程启动stdio服务器
		if err := server.ServeStdio(s); err != nil {
			logger.L().Fatalf("stdio服务器启动失败: %v", err)
		}
	}
}

// startHTTPServer 启动HTTP传输服务器（使用SSE）
func startHTTPServer(s *server.MCPServer, port int) {
	addr := fmt.Sprintf(":%d", port)
	logger.L().Infof("启动HTTP/SSE传输服务器，监听端口: %d", port)
	logger.L().Infof("MCP端点: http://localhost:%d", port)

	// 使用v0.9.0版本支持的API：创建SSE服务器
	sseServer := server.NewSSEServer(s)
	if err := sseServer.Start(addr); err != nil {
		logger.L().Fatalf("HTTP/SSE服务器启动失败: %v", err)
	}
}
