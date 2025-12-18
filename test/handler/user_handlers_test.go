package handler_test

import (
	"context"
	"testing"

	h "zabbixMcp/handler"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestGetUsersHandler(t *testing.T) {
	mock := &mockHandler{
		callFn: func(method string, params interface{}) (interface{}, error) {
			if method != "user.get" {
				return nil, nil
			}
			return []interface{}{
				map[string]interface{}{"userid": "1", "alias": "admin"},
				map[string]interface{}{"userid": "2", "alias": "guest"},
			}, nil
		},
	}

	h.SetClientPool(mock)

	res, err := h.GetUsersHandler(context.Background(), mcp.CallToolRequest{})
	if err != nil {
		t.Fatalf("GetUsersHandler 返回错误: %v", err)
	}
	if res == nil {
		t.Fatalf("GetUsersHandler 返回 nil 结果")
	}
	// StructuredContent 是 interface{}，先做类型断言再检查长度
	switch sc := res.StructuredContent.(type) {
	case []map[string]interface{}:
		if len(sc) != 2 {
			t.Fatalf("期望 2 个用户，实际: %d", len(sc))
		}
	case []interface{}:
		if len(sc) != 2 {
			t.Fatalf("期望 2 个用户，实际: %d", len(sc))
		}
	default:
		t.Fatalf("未知的 StructuredContent 类型: %T", res.StructuredContent)
	}
}
