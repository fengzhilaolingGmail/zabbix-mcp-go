package utils

import (
	"encoding/json"
)

// ToJSON 将任意接口序列化为JSON字节切片（紧凑格式）。
// 如果传入nil，返回JSON的 null 表示。
func ToJSON(v interface{}) ([]byte, error) {
	if v == nil {
		return []byte("null"), nil
	}
	return json.Marshal(v)
}

// ToJSONString 将任意接口序列化为JSON字符串（紧凑格式）。
func ToJSONString(v interface{}) (string, error) {
	b, err := ToJSON(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ToIndentedJSON 将任意接口序列化为带缩进的JSON字节切片，便于阅读。
// prefix/indent 与 json.MarshalIndent 的参数一致。
func ToIndentedJSON(v interface{}, prefix, indent string) ([]byte, error) {
	if v == nil {
		return []byte("null"), nil
	}
	return json.MarshalIndent(v, prefix, indent)
}
