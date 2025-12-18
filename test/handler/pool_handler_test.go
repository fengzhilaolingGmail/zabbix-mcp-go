package handler_test

import (
	"context"
	"testing"

	h "zabbixMcp/handler"
	zabbix "zabbixMcp/zabbix"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestGetInstancesInfoHandler_WithNilPool(t *testing.T) {
	// 将 pool 注入为 nil，确保处理器能正常返回空结果而不是 panic
	h.SetClientPool(nil)

	res, err := h.GetInstancesInfoHandler(context.Background(), mcp.CallToolRequest{})
	if err != nil {
		t.Fatalf("GetInstancesInfoHandler 在 nil pool 时返回错误: %v", err)
	}
	if res == nil {
		t.Fatalf("GetInstancesInfoHandler 返回 nil 结果")
	}
	sc, ok := res.StructuredContent.(map[string]interface{})
	if !ok {
		t.Fatalf("期望 StructuredContent 为 map[string]interface{}, 实际: %T", res.StructuredContent)
	}
	// data 应为一个数组（即使为空）
	data, exists := sc["data"]
	if !exists {
		t.Fatalf("返回的封装中缺少 data 字段")
	}
	switch d := data.(type) {
	case []interface{}:
		// ok
		_ = d
	case []map[string]interface{}:
		_ = d
	case []zabbix.ClientInfo:
		_ = d
	default:
		t.Fatalf("unexpected data type: %T", data)
	}
}
