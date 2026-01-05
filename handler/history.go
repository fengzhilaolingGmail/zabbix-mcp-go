/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2026-01-04 10:56:53
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2026-01-05 11:27:58
 * @FilePath: \zabbix-mcp-go\handler\history.go
 * @Description: 历史数据相关功能
 * @Copyright: Copyright (c) 2025 by fengzhilaoling@gmail.com, All Rights Reserved.
 */

package handler

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"zabbixMcp/logger"
	"zabbixMcp/models"
	"zabbixMcp/server"

	"github.com/mark3labs/mcp-go/mcp"
)

// 解析类似 "7d 15m 3h" 的持续时间字符串为秒
func parseDurationString(s string) (int64, error) {
	s = strings.ToLower(strings.TrimSpace(s))
	if s == "" {
		return 0, nil
	}
	re := regexp.MustCompile(`(?i)(\d+)\s*([dhms])`)
	matches := re.FindAllStringSubmatch(s, -1)
	if len(matches) == 0 {
		return 0, fmt.Errorf("invalid duration: %s", s)
	}
	var total int64
	for _, m := range matches {
		if len(m) < 3 {
			continue
		}
		num, err := strconv.ParseInt(m[1], 10, 64)
		if err != nil {
			return 0, err
		}
		unit := strings.ToLower(m[2])
		switch unit {
		case "d":
			total += num * 86400
		case "h":
			total += num * 3600
		case "m":
			total += num * 60
		case "s":
			total += num
		default:
			return 0, fmt.Errorf("unknown duration unit: %s", unit)
		}
	}
	return total, nil
}

// parseTimeParam 支持整型/浮点数时间戳或字符串时间（RFC3339, "2006-01-02 15:04:05",
// 以及常见的带/不带前导零的日期格式，例如 "2026-1-2 15:00:00" 或 "2026/1/2 15:00:00"）。
// 返回 Unix 秒数。需要传入时区 location，用于解析无时区的时间字符串以及获取当前时间计算相对时间时使用。
func parseTimeParam(v interface{}, loc *time.Location) (int, error) {
	if v == nil {
		return 0, nil
	}
	switch tv := v.(type) {
	case float64:
		return int(tv), nil
	case int:
		return tv, nil
	case int64:
		return int(tv), nil
	case string:
		s := strings.TrimSpace(tv)
		if s == "" {
			return 0, nil
		}
		// 纯数字字符串，当作时间戳
		if _, err := strconv.ParseInt(s, 10, 64); err == nil {
			v64, _ := strconv.ParseInt(s, 10, 64)
			return int(v64), nil
		}
		// 尝试按常见时间格式解析（在指定时区），包含带/不带前导零及斜杠分隔的常见格式
		layouts := []string{
			time.RFC3339,
			"2006-01-02 15:04:05",
			"2006/01/02 15:04:05",
			"2006-1-2 15:04:05",
			"2006/1/2 15:04:05",
			"2006-01-02",
			"2006/01/02",
			"2006-1-2",
			"2006/1/2",
		}
		for _, layout := range layouts {
			if t, err := time.ParseInLocation(layout, s, loc); err == nil {
				return int(t.Unix()), nil
			}
		}
		return 0, fmt.Errorf("unsupported time string format: %s", s)
	default:
		return 0, fmt.Errorf("unsupported time param type: %T", v)
	}
}

// GetHistoryHandler 处理 MCP 请求并调用 server.GetHistory
func GetHistoryHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	instance := ""
	var hostIDs []string
	var itemIDs []string
	var historyPtr *int
	// 是否返回 summary，默认 true
	summary := true
	var startTimeStr string
	var stopTimeStr string
	var sortField interface{} = "clock" // 默认 clock
	limit := 0
	var output interface{}
	timezone := "Asia/Shanghai"
	var timeRangeStr string

	if args, ok := req.Params.Arguments.(map[string]interface{}); ok {
		if v, ok2 := args["instance"].(string); ok2 {
			instance = v
		}
		if arr, ok := args["host_ids"].([]interface{}); ok {
			for _, e := range arr {
				if s, ok := e.(string); ok && s != "" {
					hostIDs = append(hostIDs, s)
				} else if n, ok := e.(float64); ok {
					hostIDs = append(hostIDs, strconv.FormatInt(int64(n), 10))
				}
			}
		} else if v, ok := args["host_ids"].(string); ok {
			if v != "" {
				hostIDs = append(hostIDs, v)
			}
		}
		if arr, ok := args["item_ids"].([]interface{}); ok {
			for _, e := range arr {
				if s, ok := e.(string); ok && s != "" {
					itemIDs = append(itemIDs, s)
				} else if n, ok := e.(float64); ok {
					itemIDs = append(itemIDs, strconv.FormatInt(int64(n), 10))
				}
			}
		} else if v, ok := args["item_ids"].(string); ok {
			if v != "" {
				itemIDs = append(itemIDs, v)
			}
		}

		// history 参数，可为数字或字符串数字
		if hv, ok := args["history"]; ok {
			switch hh := hv.(type) {
			case float64:
				hi := int(hh)
				historyPtr = &hi
			case int:
				hi := hh
				historyPtr = &hi
			case string:
				if hh != "" {
					if n, err := strconv.ParseInt(hh, 10, 64); err == nil {
						ni := int(n)
						historyPtr = &ni
					}
				}
			}
		}

		// 获取时区（可选），默认 Asia/Shanghai
		if tz, ok := args["timezone"].(string); ok && tz != "" {
			timezone = tz
		}
		// 时区解析将在后面统一完成（避免在 args 解析期间重复加载）
		// time_range 支持类似 "7d 3h 15m" 的持续时间字符串（相对于当前时间）
		if tr, ok := args["time_range"].(string); ok && strings.TrimSpace(tr) != "" {
			timeRangeStr = tr
		}
		// 支持 startTime/stopTime 或 snake_case time_from/time_till，均以字符串形式接收
		if v, ok := args["start_time"]; ok {
			startTimeStr = fmt.Sprintf("%v", v)
		} else if v, ok := args["time_from"]; ok {
			startTimeStr = fmt.Sprintf("%v", v)
		}
		if v, ok := args["end_time"]; ok {
			stopTimeStr = fmt.Sprintf("%v", v)
		} else if v, ok := args["time_till"]; ok {
			stopTimeStr = fmt.Sprintf("%v", v)
		}

		if sf, ok := args["sortfield"]; ok {
			sortField = sf
		}
		if l, ok := args["limit"].(float64); ok {
			limit = int(l)
		} else if l, ok := args["limit"].(int); ok {
			limit = l
		}
		if out, ok := args["output"]; ok {
			output = out
		}
		// summary 开关，支持 bool/string/number
		if sv, ok := args["summary"]; ok {
			summary = false
			switch vv := sv.(type) {
			case bool:
				summary = vv
			case string:
				s := strings.ToLower(strings.TrimSpace(vv))
				if s == "true" || s == "1" || s == "yes" {
					summary = true
				} else if s == "false" || s == "0" || s == "no" {
					summary = false
				}
			case float64:
				if int(vv) != 0 {
					summary = true
				} else {
					summary = false
				}
			}
		}
	}

	// 必须至少提供 hostids 或 itemids
	if len(hostIDs) == 0 && len(itemIDs) == 0 {
		logger.L().Warnf("GetHistoryHandler: hostids or itemids is required")
		return nil, fmt.Errorf("hostids or itemids is required")
	}

	// 默认 history 为 0（数值型）
	if historyPtr == nil {
		def := 0
		historyPtr = &def
	}

	// 解析时区
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		logger.L().Warnf("invalid timezone '%s', fallback to Asia/Shanghai: %v", timezone, err)
		loc, _ = time.LoadLocation("Asia/Shanghai")
	}

	// time_range 与 start/stop 二选一
	var parsedStart int
	var parsedStop int

	if timeRangeStr != "" {
		logger.L().Infof("%s", timeRangeStr)
		// 如果同时提供了 start/stop，则认为冲突
		if startTimeStr != "" || stopTimeStr != "" {
			return nil, fmt.Errorf("provide either time_range or startTime/stopTime, not both")
		}
		seconds, err := parseDurationString(timeRangeStr)
		if err != nil {
			return nil, fmt.Errorf("invalid time_range: %w", err)
		}
		now := time.Now().In(loc).Unix()
		parsedStop = int(now)
		parsedStart = int(now - seconds)
		output = "extend"
	}

	// 若提供了 startTime/stopTime 字符串，则解析为 Unix 秒；若未提供则保持为 0（API 支持范围查询）
	if timeRangeStr == "" {
		if startTimeStr != "" {
			if t, err := parseTimeParam(startTimeStr, loc); err == nil {
				parsedStart = t
			} else {
				return nil, fmt.Errorf("invalid startTime: %w", err)
			}
		}
		if stopTimeStr != "" {
			if t, err := parseTimeParam(stopTimeStr, loc); err == nil {
				parsedStop = t
			} else {
				return nil, fmt.Errorf("invalid stopTime: %w", err)
			}
		}
		output = "extend"
	}

	if clientPool == nil {
		return mcp.NewToolResultStructuredOnly(makeResult([]map[string]interface{}{})), nil
	}

	spec := models.ParamsHistory{
		History:   historyPtr,
		HostIDs:   hostIDs,
		ItemIDs:   itemIDs,
		TimeFrom:  parsedStart,
		TimeTill:  parsedStop,
		SortField: sortField,
		Limit:     limit,
		Output:    output,
	}

	logger.L().Infof("GetHistoryHandler: instance=%s hostids=%v itemids=%v history=%v start=%d stop=%d sortfield=%v limit=%d",
		instance, hostIDs, itemIDs, historyPtr, parsedStart, parsedStop, sortField, limit)

	histories, err := server.GetHistory(ctx, clientPool, spec, instance)
	if err != nil {
		logger.L().Errorf("调用 history.get 失败: %v", err)
		return nil, err
	}
	// 0 (默认)- 数值型float; 1 - 字符型; 2 - 日志型; 3 - 无符号数值型; 4 - 文本型; 5 - 二进制型
	// 如果是数值类型（0=float 或 3=unsigned）且 summary 开启，按 host:item 分组计算 summary 并返回
	if summary && historyPtr != nil && (*historyPtr == 0 || *historyPtr == 3) {
		type stat struct {
			count int
			sum   float64
			max   float64
			min   float64
		}
		stats := make(map[string]*stat)
		// 聚合统计
		for _, h := range histories {
			valRaw, ok := h["value"]
			if !ok || valRaw == nil {
				continue
			}
			var fv float64
			switch v := valRaw.(type) {
			case float64:
				fv = v
			case string:
				if parsed, err := strconv.ParseFloat(v, 64); err == nil {
					fv = parsed
				} else {
					continue
				}
			case int:
				fv = float64(v)
			case int64:
				fv = float64(v)
			default:
				continue
			}
			itemid := ""
			if id, ok := h["itemid"]; ok {
				itemid = fmt.Sprintf("%v", id)
			}
			hostid := ""
			if hid, ok := h["hostid"]; ok {
				hostid = fmt.Sprintf("%v", hid)
			}
			key := itemid
			if hostid != "" {
				key = hostid + ":" + itemid
			}
			s, ok := stats[key]
			if !ok {
				stats[key] = &stat{count: 1, sum: fv, max: fv, min: fv}
			} else {
				s.count++
				s.sum += fv
				if fv > s.max {
					s.max = fv
				}
				if fv < s.min {
					s.min = fv
				}
			}
		}
		// 构造嵌套 map: map[itemid][hostid] -> summary
		result := make(map[string]map[string]map[string]interface{})
		for key, s := range stats {
			if s.count == 0 {
				continue
			}
			hostid := ""
			itemid := ""
			if strings.Contains(key, ":") {
				parts := strings.SplitN(key, ":", 2)
				hostid = parts[0]
				itemid = parts[1]
			} else {
				itemid = key
			}
			if _, ok := result[itemid]; !ok {
				result[itemid] = make(map[string]map[string]interface{})
			}
			result[itemid][hostid] = map[string]interface{}{
				"count": s.count,
				"avg":   s.sum / float64(s.count),
				"max":   s.max,
				"min":   s.min,
			}
		}
		return mcp.NewToolResultStructuredOnly(makeResult(result)), nil
	}

	// 非数值类型直接返回原始 histories
	return mcp.NewToolResultStructuredOnly(makeResult(histories)), nil
}
