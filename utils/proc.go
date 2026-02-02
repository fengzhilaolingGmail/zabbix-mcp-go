/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2025-12-16 09:41:25
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2026-01-24 17:18:31
 * @FilePath: \zabbix-mcp-go\utils\proc.go
 * @Description: 文件详情
 * @Copyright: Copyright (c) 2025 by fengzhilaoling@gmail.com, All Rights Reserved.
 */
package utils

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"
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

func JsonInt(v interface{}, def int) int {
	switch val := v.(type) {
	case string:
		if val == "" {
			return def
		}
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	case float64:
		return int(val)
	case int:
		return val
	case int64:
		return int(val)
	}
	return def
}

// ParseZabbixSelectParam 解析Zabbix API的select类参数，兼容多类型输入并遵循Zabbix约定
// params: 原始参数map（通常是JSON解码后的map[string]interface{}）
// key: 要解析的参数键名（如"selectGraphs"）
// defaultValue: 解析失败/空值/未知类型时的回退默认值（如"extend"）
// return: 解析后的结果，类型为interface{}（实际为string或[]string）
func ParseZabbixSelectParam(params map[string]interface{}, key string, defaultValue string) interface{} {
	// 1. 从map中获取指定键值，不存在则返回默认值
	raw, ok := params[key]
	if !ok {
		return defaultValue
	}

	// 2. 处理JSON解码常见的[]interface{}类型（最常用场景）
	if arr, ok := raw.([]interface{}); ok {
		tmp := make([]string, 0, len(arr))
		// 遍历转换为字符串切片，过滤非字符串元素
		for _, it := range arr {
			if s, ok := it.(string); ok {
				tmp = append(tmp, s)
			}
		}
		return handleParsedSlice(tmp, defaultValue)
	}

	// 3. 处理原生[]string类型
	if arrS, ok := raw.([]string); ok {
		return handleParsedSlice(arrS, defaultValue)
	}

	// 4. 处理直接传入的字符串类型（如直接传"extend"）
	if s, ok := raw.(string); ok {
		return s
	}

	// 5. 未知类型，返回默认值
	return defaultValue
}

// handleParsedSlice 处理解析后的字符串切片，遵循Zabbix API约定
// slice: 解析后的纯字符串切片
// defaultValue: 空切片时的回退值
// return: 处理结果（string或[]string）
func handleParsedSlice(slice []string, defaultValue string) interface{} {
	// 空切片返回默认值
	if len(slice) == 0 {
		return defaultValue
	}
	// 首元素为默认值则返回字符串类型，否则返回原切片
	if slice[0] == defaultValue {
		return defaultValue
	}
	return slice
}

// ParseSliceFromMap 从map[string]interface{}中解析指定key为目标切片类型（泛型版，Go1.18+）
// 参数：
//
//	m: 原始map[string]interface{}（如JSON解析后的args）
//	key: 要解析的map的key（如"groups"）
//	target: 目标切片的指针（如&[]models.Groups{}）
//
// 返回值：
//
//	error: 解析过程中的错误（参数非法/key不存在/序列化/反序列化失败）
func ParseSliceFromMap(m map[string]interface{}, key string, target interface{}) error {
	// 1. 校验原始map是否为nil
	if m == nil {
		return errors.New("原始map为nil，无法解析")
	}

	// 2. 校验目标指针是否合法（必须是切片的指针，否则反序列化失败）
	if target == nil {
		return errors.New("目标切片指针不能为nil")
	}

	// 3. 检查map中是否存在指定key
	val, ok := m[key]
	if !ok {
		return fmt.Errorf("map中不存在指定key: %s", key)
	}

	// 4. 校验key对应值是否为nil（避免Marshal nil导致空值）
	if val == nil {
		return fmt.Errorf("key: %s 对应值为nil", key)
	}

	// 5. 将interface{}序列化为JSON字节流
	jsonBytes, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("序列化为JSON失败: %w, key: %s", err, key)
	}

	// 6. 将JSON字节流反序列化为目标切片类型
	if err := json.Unmarshal(jsonBytes, target); err != nil {
		return fmt.Errorf("反序列化为目标类型失败: %w, key: %s", err, key)
	}

	return nil
}
