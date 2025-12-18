package handler_test

import (
	"context"
	"testing"

	h "zabbixMcp/handler"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestGetInstancesInfoHandler_WithNilPool(t *testing.T) {
	// 将 pool 注入为 nil，确保处理器能正常返回空结果而不是 panic
	h.SetClientPool(nil)

	_, err := h.GetInstancesInfoHandler(context.Background(), mcp.CallToolRequest{})
	if err != nil {
		t.Fatalf("GetInstancesInfoHandler 在 nil pool 时返回错误: %v", err)
	}
}
