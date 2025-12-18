package handler_test

import (
	"context"
	"testing"

	h "zabbixMcp/handler"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestGetHostsHandler(t *testing.T) {
	mock := &mockHandler{
		callFn: func(method string, params interface{}) (interface{}, error) {
			if method != "host.get" {
				return nil, nil
			}
			return []map[string]interface{}{
				{"hostid": "10101", "host": "web01"},
			}, nil
		},
	}

	h.SetClientPool(mock)

	res, err := h.GetHostsHandler(context.Background(), mcp.CallToolRequest{})
	if err != nil {
		t.Fatalf("GetHostsHandler 返回错误: %v", err)
	}
	if res == nil {
		t.Fatalf("GetHostsHandler 返回 nil 结果")
	}
	// 统一返回结构为 map[string]interface{"ok":bool, "data": ...}
	sc, ok := res.StructuredContent.(map[string]interface{})
	if !ok {
		t.Fatalf("期望 StructuredContent 为 map[string]interface{}, 实际: %T", res.StructuredContent)
	}
	data, exists := sc["data"]
	if !exists {
		t.Fatalf("返回的封装中缺少 data 字段")
	}
	switch d := data.(type) {
	case []map[string]interface{}:
		if len(d) != 1 {
			t.Fatalf("期望 1 个主机，实际: %d", len(d))
		}
	case []interface{}:
		if len(d) != 1 {
			t.Fatalf("期望 1 个主机，实际: %d", len(d))
		}
	default:
		t.Fatalf("未知的 data 类型: %T", data)
	}
}
