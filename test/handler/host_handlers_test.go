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
	// StructuredContent 是 interface{}，先进行类型断言再检查长度
	switch sc := res.StructuredContent.(type) {
	case []map[string]interface{}:
		if len(sc) != 1 {
			t.Fatalf("期望 1 个主机，实际: %d", len(sc))
		}
	case []interface{}:
		if len(sc) != 1 {
			t.Fatalf("期望 1 个主机，实际: %d", len(sc))
		}
	default:
		t.Fatalf("未知的 StructuredContent 类型: %T", res.StructuredContent)
	}
}
