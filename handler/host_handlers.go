package handler

import (
	"context"
	"fmt"
	"zabbixMcp/server"

	"github.com/mark3labs/mcp-go/mcp"
)

// GetHostsHandler 通过注入的 ZabbixClientHandler 调用 host.get 并返回结果
func GetHostsHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if clientPool == nil {
		return &mcp.CallToolResult{StructuredContent: makeResult([]map[string]interface{}{})}, nil
	}

	params := map[string]interface{}{"output": "extend"}
	hosts, err := server.GetHosts(clientPool, params)
	if err != nil {
		return nil, fmt.Errorf("调用 host.get 失败: %w", err)
	}
	return &mcp.CallToolResult{StructuredContent: makeResult(hosts)}, nil
}
