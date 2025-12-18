package models

import (
	"encoding/json"
	"fmt"
)

// JSONRPCRequest 描述一次 JSON-RPC 调用请求
type JSONRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      int         `json:"id"`
	Auth    string      `json:"auth,omitempty"`
}

// JSONRPCResponse 描述 JSON-RPC 响应
type JSONRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	Error   *RPCError       `json:"error"`
	ID      int             `json:"id"`
}

// RPCError 表示 Zabbix API 返回的错误信息
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

// Error 实现 error 接口
func (e *RPCError) Error() string {
	return fmt.Sprintf("Zabbix API Error %d: %s (%s)", e.Code, e.Message, e.Data)
}
