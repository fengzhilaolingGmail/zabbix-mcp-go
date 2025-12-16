/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2025-12-16 20:29:27
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2025-12-16 21:11:31
 * @FilePath: \zabbix-mcp-go\zabbix\models.go
 * @Description: Zabbix API数据模型
 * Copyright (c) 2025 by fengzhilaoling@gmail.com, All Rights Reserved.
 */

package zabbix

import "fmt"

// JSONRPCRequest JSON-RPC请求结构
type JSONRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`        // JSON-RPC 协议的版本号，本例中为 2.0
	Method  string      `json:"method"`         // 要调用的API方法名称
	Params  interface{} `json:"params"`         // 方法参数，可以是任何类型
	ID      int         `json:"id"`             // 请求的唯一标识符，用于响应请求
	Auth    string      `json:"auth,omitempty"` // 用户认证令牌，用于API认证
}

type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"` // JSON-RPC 协议的版本号，本例中为 2.0
	Result  interface{} `json:"result"`  // 方法调用的结果，可以是任何类型
	Error   *RPCError   `json:"error"`   // 如果调用失败，则包含错误信息
	ID      int         `json:"id"`      // 请求的唯一标识符，用于响应请求
}

// RPCError RPC错误结构
type RPCError struct {
	Code    int    `json:"code"`    // 错误代码
	Message string `json:"message"` // 错误消息
	Data    string `json:"data"`    // 额外的错误数据
}

// Error 实现error接口
func (e *RPCError) Error() string {
	return fmt.Sprintf("Zabbix API Error %d: %s (%s)", e.Code, e.Message, e.Data)
}
