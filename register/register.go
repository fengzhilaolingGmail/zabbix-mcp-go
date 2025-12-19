/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2025-12-18 10:06:45
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2025-12-18 12:39:25
 * @FilePath: \zabbix-mcp-go\register\register.go
 * @Description: 文件详情
 * @Copyright: Copyright (c) 2025 by fengzhilaoling@gmail.com, All Rights Reserved.
 */
package register

import (
	"github.com/mark3labs/mcp-go/server"
)

func Registers(s *server.MCPServer) {
	// 注册 ClientPool 相关工具
	registerInstances(s)
	registerUser(s)
}
