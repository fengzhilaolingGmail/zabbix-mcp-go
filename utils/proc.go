/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2025-12-16 09:41:25
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2025-12-20 13:58:19
 * @FilePath: \zabbix-mcp-go\utils\proc.go
 * @Description: 文件详情
 * @Copyright: Copyright (c) 2025 by fengzhilaoling@gmail.com, All Rights Reserved.
 */
package utils

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
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

var passwordCharset = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!$^&*()-_=+[]{}|;:,.<>?/~")

// GenerateSecurePassword 随机生成指定长度的高强度密码，包含大小写、数字和特殊字符。
func GenerateSecurePassword(length int) (string, error) {
	if length <= 0 {
		return "", nil
	}
	result := make([]byte, length)
	charsetLen := big.NewInt(int64(len(passwordCharset)))
	for i := 0; i < length; i++ {
		rnd, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", err
		}
		result[i] = passwordCharset[rnd.Int64()]
	}
	return string(result), nil
}
