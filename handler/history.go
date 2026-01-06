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
	"math"
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

// GetHistoryCompareHandler 实现同比（current vs previous）比较
// 支持按小时或按天的同比（period: "hour" or "day"，默认 "day"）
// 计算每个 host+item 的 max/min/avg
// 请求数据时单次不超过一天（86400s），若范围超过一天则分批次请求并汇总
func GetHistoryCompareHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	instance := ""
	var hostIDs []string
	var itemIDs []string
	var historyPtr *int
	period := "day" // hour 或 day
	timezone := "Asia/Shanghai"
	var startTimeStr string
	var stopTimeStr string
	var timeRangeStr string
	pctFormat := "number" // "number" (default) or "string" to return percent string like "12.34%"

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
		if tz, ok := args["timezone"].(string); ok && tz != "" {
			timezone = tz
		}
		if p, ok := args["period"].(string); ok && (p == "hour" || p == "day") {
			period = p
		}
		if tr, ok := args["time_range"].(string); ok && strings.TrimSpace(tr) != "" {
			timeRangeStr = tr
		}
		if pf, ok := args["pct_format"].(string); ok && pf != "" {
			pf = strings.ToLower(strings.TrimSpace(pf))
			if pf == "number" || pf == "string" {
				pctFormat = pf
			}
		}
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
	}

	if len(hostIDs) == 0 && len(itemIDs) == 0 {
		return nil, fmt.Errorf("hostids or itemids is required")
	}
	if historyPtr == nil {
		def := 0
		historyPtr = &def
	}
	if loc, err := time.LoadLocation(timezone); err != nil {
		logger.L().Warnf("invalid timezone '%s', fallback to Asia/Shanghai: %v", timezone, err)
		timezone = "Asia/Shanghai"
	} else {
		_ = loc
	}

	// 解析时间范围
	loc, _ := time.LoadLocation(timezone)
	var parsedStart int
	var parsedStop int
	if timeRangeStr != "" {
		seconds, err := parseDurationString(timeRangeStr)
		if err != nil {
			return nil, fmt.Errorf("invalid time_range: %w", err)
		}
		now := time.Now().In(loc).Unix()
		parsedStop = int(now)
		parsedStart = int(now - seconds)
	} else {
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
	}

	if clientPool == nil {
		return mcp.NewToolResultStructuredOnly(makeResult([]map[string]interface{}{})), nil
	}

	// helper: fetch histories in batches not exceeding one day
	const maxBatchSeconds = 86400
	fetchBatched := func(ctx context.Context, baseSpec models.ParamsHistory, start, stop int) ([]map[string]interface{}, error) {
		var all []map[string]interface{}
		if stop <= start {
			return all, nil
		}
		cursor := start
		for cursor < stop {
			end := cursor + maxBatchSeconds
			if end > stop {
				end = stop
			}
			spec := baseSpec
			spec.TimeFrom = cursor
			spec.TimeTill = end
			h, err := server.GetHistory(ctx, clientPool, spec, instance)
			if err != nil {
				return nil, err
			}
			all = append(all, h...)
			cursor = end
		}
		return all, nil
	}

	// stat struct and aggregation
	type stat struct {
		count int
		sum   float64
		max   float64
		min   float64
	}
	// (aggregation replaced by per-bucket aggregation below)

	baseSpec := models.ParamsHistory{
		History: historyPtr,
		HostIDs: hostIDs,
		ItemIDs: itemIDs,
		Output:  "extend",
	}

	// For time-series per-bucket comparison: build buckets of size bucketSeconds
	bucketSeconds := 86400
	if period == "hour" {
		bucketSeconds = 3600
	}
	length := parsedStop - parsedStart
	if length <= 0 {
		return nil, fmt.Errorf("invalid time range")
	}
	numBuckets := (length + bucketSeconds - 1) / bucketSeconds

	// Instead of fetching all data at once, fetch per-bucket (including previous bucket) and aggregate immediately.
	statsPerBucket := make(map[int]map[string]*stat)
	for i := -1; i < numBuckets; i++ {
		statsPerBucket[i] = make(map[string]*stat)
		var bstart, bend int
		if i == -1 {
			bstart = parsedStart - bucketSeconds
			if bstart < 0 {
				bstart = 0
			}
			bend = parsedStart
		} else {
			bstart = parsedStart + i*bucketSeconds
			bend = bstart + bucketSeconds
			if bend > parsedStop {
				bend = parsedStop
			}
		}
		if bend <= bstart {
			continue
		}
		// fetch this bucket (fetchBatched will respect maxBatchSeconds)
		hist, err := fetchBatched(ctx, baseSpec, bstart, bend)
		if err != nil {
			return nil, err
		}
		// aggregate into statsPerBucket[i]
		for _, h := range hist {
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
			s, ok := statsPerBucket[i][key]
			if !ok {
				statsPerBucket[i][key] = &stat{count: 1, sum: fv, max: fv, min: fv}
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
	}

	// build final: map[itemid][hostid] -> []entries per bucket (each with current/previous/delta/pct_change)
	final := make(map[string]map[string][]map[string]interface{})
	// collect all keys across buckets
	allKeys := make(map[string]struct{})
	for i := -1; i < numBuckets; i++ {
		for k := range statsPerBucket[i] {
			allKeys[k] = struct{}{}
		}
	}

	for key := range allKeys {
		hostid := ""
		itemid := ""
		if strings.Contains(key, ":") {
			parts := strings.SplitN(key, ":", 2)
			hostid = parts[0]
			itemid = parts[1]
		} else {
			itemid = key
		}
		if _, ok := final[itemid]; !ok {
			final[itemid] = make(map[string][]map[string]interface{})
		}
		// build entries for each bucket
		entries := make([]map[string]interface{}, 0, numBuckets)
		for i := 0; i < numBuckets; i++ {
			bucketStart := parsedStart + i*bucketSeconds
			bucketEnd := bucketStart + bucketSeconds
			if bucketEnd > parsedStop {
				bucketEnd = parsedStop
			}
			curStat := statsPerBucket[i][key]
			prevStat := statsPerBucket[i-1][key]
			cur := map[string]interface{}{"count": 0, "avg": 0.0, "max": 0.0, "min": 0.0}
			prev := map[string]interface{}{"count": 0, "avg": 0.0, "max": 0.0, "min": 0.0}
			var curAvg, prevAvg float64
			if curStat != nil && curStat.count > 0 {
				cur["count"] = curStat.count
				curAvg = curStat.sum / float64(curStat.count)
				cur["avg"] = curAvg
				cur["max"] = curStat.max
				cur["min"] = curStat.min
			}
			if prevStat != nil && prevStat.count > 0 {
				prev["count"] = prevStat.count
				prevAvg = prevStat.sum / float64(prevStat.count)
				prev["avg"] = prevAvg
				prev["max"] = prevStat.max
				prev["min"] = prevStat.min
			}
			// compute delta and pct_change (pct_change is 0 when previous avg is zero or missing)
			var delta float64
			var pctChange interface{}
			delta = curAvg - prevAvg
			var pctVal float64
			if prevStat == nil || prevStat.count == 0 || prevAvg == 0 {
				pctVal = 0.0
			} else {
				pctVal = (delta / prevAvg) * 100.0
			}
			// round to 2 decimal places
			pctVal = math.Round(pctVal*100) / 100
			if pctFormat == "string" {
				pctChange = fmt.Sprintf("%.2f%%", pctVal)
			} else {
				pctChange = pctVal
			}

			entry := map[string]interface{}{
				"time_from":  bucketStart,
				"time_till":  bucketEnd,
				"current":    cur,
				"previous":   prev,
				"delta":      delta,
				"pct_change": pctChange,
			}
			entries = append(entries, entry)
		}
		final[itemid][hostid] = entries
	}

	wrapper := map[string]interface{}{"period": period, "data": final}
	return mcp.NewToolResultStructuredOnly(makeResult(wrapper)), nil
}
