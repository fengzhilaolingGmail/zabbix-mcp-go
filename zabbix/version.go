/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2025-12-16 20:54:52
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2026-01-09 15:39:35
 * @FilePath: \zabbix-mcp-go\zabbix\version.go
 * @Description: 版本检测相关功能
 * Copyright (c) 2025 by fengzhilaoling@gmail.com, All Rights Reserved.
 */
package zabbix

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"zabbixMcp/logger"
	"zabbixMcp/models"
)

// VersionInfo Zabbix版本信息
type VersionInfo struct {
	Major int    // 主版本号
	Minor int    // 次版本号
	Patch int    // 补丁版本号解析失败时为0
	Full  string // 完整版本字符串
}

// VersionDetector 版本检测器
type VersionDetector struct {
	client *ZabbixClient
}

// NewVersionDetector 创建版本检测器
func NewVersionDetector(client *ZabbixClient) *VersionDetector {
	return &VersionDetector{client: client}
}

// DetectVersion 检测Zabbix版本
func (vd *VersionDetector) DetectVersion(ctx context.Context) (*VersionInfo, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	// 先尝试从 client 缓存读取
	if vd.client != nil {
		if cached := vd.client.GetCachedVersion(); cached != nil {
			return cached, nil
		}
	}

	// 获取API版本信息 - 使用内部调用避免循环依赖
	// Zabbix API要求params为空数组[]而不是nil
	result, err := vd.client.callWithAuth(ctx, "apiinfo.version", []interface{}{}, "") // 使用旧方式获取API版本信息
	if err != nil {
		logger.L().Warnf("获取API版本失败: %v, 尝试新方法", err)
		result, err = vd.client.callWithHeaderAuth(ctx, "apiinfo.version", nil, "") // 使用新方式获取API版本信息
		if err != nil {
			return nil, fmt.Errorf("获取API版本失败: %w", err)
		}
	}
	var apiVersion string
	if err := json.Unmarshal(result, &apiVersion); err != nil {
		return nil, fmt.Errorf("API版本响应格式错误: %w", err)
	}

	// 解析版本号
	version, err := vd.parseVersion(apiVersion)
	if err != nil {
		return nil, fmt.Errorf("解析版本号失败: %w", err)
	}
	version.Full = apiVersion

	// 将结果写入 client 缓存
	if vd.client != nil {
		vd.client.SetCachedVersion(version)
	}

	return version, nil
}

// parseVersion 解析版本字符串
func (vd *VersionDetector) parseVersion(versionStr string) (*VersionInfo, error) {
	// 移除前缀
	versionStr = strings.TrimPrefix(versionStr, "v")

	// 分割版本号
	parts := strings.Split(versionStr, ".")
	if len(parts) < 2 {
		return nil, fmt.Errorf("版本格式不正确: %s", versionStr)
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("解析主版本号失败: %w", err)
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("解析次版本号失败: %w", err)
	}

	patch := 0
	if len(parts) > 2 {
		patch, err = strconv.Atoi(parts[2])
		if err != nil {
			patch = 0 // 如果解析失败，默认为0
		}
	}

	return &VersionInfo{
		Major: major,
		Minor: minor,
		Patch: patch,
		Full:  versionStr,
	}, nil
}

// ParseVersion 对 parseVersion 的导出包装，便于测试
func (vd *VersionDetector) ParseVersion(versionStr string) (*VersionInfo, error) {
	return vd.parseVersion(versionStr)
}

// getDefaultFeatures 获取默认功能集
func (vd *VersionDetector) getDefaultFeatures() map[string]bool {
	return map[string]bool{
		"host_management":      true,
		"item_management":      true,
		"trigger_management":   true,
		"template_management":  true,
		"event_acknowledgment": true,
	}
}

// 在 version.go 中添加更详细的版本特性映射
func (vd *VersionDetector) GetDetailedVersionFeatures() map[string]interface{} {
	version, err := vd.DetectVersion(context.Background())
	if err != nil {
		// 将 map[string]bool 转换为 map[string]interface{}
		defaultFeatures := vd.getDefaultFeatures()
		result := make(map[string]interface{})
		for k, v := range defaultFeatures {
			result[k] = v
		}
		return result
	}

	features := make(map[string]interface{})

	// API端点支持
	features["endpoints"] = map[string]bool{
		"problem.get":    version.Major >= 4,
		"sla.get":        version.Major >= 5,
		"authentication": version.Major >= 7,
		"connector":      version.Major >= 6,
		"proxygroup":     version.Major >= 7,
	}

	// 参数支持
	features["parameters"] = map[string]bool{
		"selectTags":          version.Major >= 4,
		"selectDependencies":  version.Major >= 4,
		"selectPreprocessing": version.Major >= 4,
		"templateSelectTags":  version.Major >= 5, // template.get
	}

	return features
}

// AdaptAPIParams 根据版本适配API参数
func (vd *VersionDetector) AdaptAPIParams(method string, spec models.ParamSpec) map[string]interface{} {
	version, err := vd.DetectVersion(context.Background())
	logger.L().Info(version.Full)
	var params map[string]interface{}
	if spec != nil {
		params = spec.BuildParams()
	} else {
		params = map[string]interface{}{}
	}

	if err != nil {
		return params
	}

	adaptedParams := make(map[string]interface{}, len(params))
	for k, v := range params {
		adaptedParams[k] = v
	}

	// 根据版本调整参数
	switch method {
	case "host.get":
		if version.Major < 4 {
			// 旧版本不支持某些参数
			delete(adaptedParams, "selectTags")
			// adaptedParams["output"] = []string{"hostid", "name"}
		}
		if version.Major < 6 {
			// Older Zabbix versions expect different parameter name for requesting host groups
			// Map both camelCase and snake_case variants to the legacy `selectGroups` key
			if v, ok := adaptedParams["selectHostGroups"]; ok {
				adaptedParams["selectGroups"] = v
				delete(adaptedParams, "selectHostGroups")
			}
			if v, ok := adaptedParams["select_host_groups"]; ok {
				adaptedParams["selectGroups"] = v
				delete(adaptedParams, "select_host_groups")
			}
		}
	// ========================= User API =========================
	case "user.get":
		if version.Major > 5 {
			if f, ok := adaptedParams["filter"].(map[string]interface{}); ok {
				delete(f, "alias")
				if len(f) == 0 {
					delete(adaptedParams, "filter")
				}
			}
		} else {
			if f, ok := adaptedParams["filter"].(map[string]interface{}); ok {
				delete(f, "username")
				if len(f) == 0 {
					delete(adaptedParams, "filter")
				}
			}
		}
	case "user.create":
		// Zabbix 5.x 及更早版本不使用 `username` 字段，使用 `alias`。
		// 对于 6.x/7.x，使用 `username`。
		if version.Major <= 5 {
			if un, ok := adaptedParams["username"]; ok {
				if s, ok2 := un.(string); ok2 {
					adaptedParams["alias"] = s
				}
				delete(adaptedParams, "username")
			}
			delete(adaptedParams, "roleid")
		} else {
			// >=6: 尽量保留 username；如果只有 alias 提供也不会出错
			delete(adaptedParams, "alias")
		}
	case "usergroup.get":
		if version.Major < 6 {
			if v, ok := adaptedParams["selectUsers"]; ok {
				if fields, ok := v.([]string); ok {
					filtered := fields[:0] // 复用底层数组
					for _, f := range fields {
						if f != "username" {
							filtered = append(filtered, f)
						}
					}
					if len(filtered) == 0 {
						delete(adaptedParams, "selectUsers")
					} else {
						adaptedParams["selectUsers"] = filtered
					}
				}
			}
		} else {
			if v, ok := adaptedParams["selectUsers"]; ok {
				if fields, ok := v.([]string); ok {
					filtered := fields[:0] // 复用底层数组
					for _, f := range fields {
						if f != "alias" {
							filtered = append(filtered, f)
						}
					}
					if len(filtered) == 0 {
						delete(adaptedParams, "selectUsers")
					} else {
						adaptedParams["selectUsers"] = filtered
					}
				}
			}
		}
	// ========================= Item API =========================
	case "item.get":
		if version.Major < 4 {
			delete(adaptedParams, "selectTags")
			delete(adaptedParams, "selectPreprocessing")
		}
	case "trigger.get":
		if version.Major < 4 {
			delete(adaptedParams, "selectTags")
			delete(adaptedParams, "selectDependencies")
		}
	case "template.get":
		if version.Major < 5 {
			delete(adaptedParams, "selectTags")
		}
	// ========================= dashboard 仪表盘 =========================
	case "dashboard.get":
		if version.Major < 6 {
			delete(adaptedParams, "selectPages")
			adaptedParams["selectWidgets"] = "extend"
		}
	case "dashboard.create":
		// helper to convert various numeric types to int
		toInt := func(v interface{}) (int, bool) {
			switch t := v.(type) {
			case int:
				return t, true
			case int32:
				return int(t), true
			case int64:
				return int(t), true
			case float64:
				return int(t), true
			case float32:
				return int(t), true
			case string:
				if n, err := strconv.Atoi(t); err == nil {
					return n, true
				}
			}
			return 0, false
		}

		if version.Major < 6 {
			// Zabbix 5.x: 将 pages 转换为 widgets，并支持自动分页
			// Zabbix 5.x 限制：x:0-23, y:0-63, width:1-24, height:1-32
			if pages, ok := adaptedParams["pages"]; ok {
				var allWidgets []interface{}
				// 提取所有 widgets
				extract := func(w interface{}) {
					switch ws := w.(type) {
					case []interface{}:
						allWidgets = append(allWidgets, ws...)
					case []map[string]interface{}:
						for _, it := range ws {
							allWidgets = append(allWidgets, it)
						}
					case map[string]interface{}:
						allWidgets = append(allWidgets, ws)
					}
				}

				switch p := pages.(type) {
				case []interface{}:
					for _, page := range p {
						if pm, ok := page.(map[string]interface{}); ok {
							if w, ok := pm["widgets"]; ok {
								extract(w)
							}
						}
					}
				case []map[string]interface{}:
					for _, pm := range p {
						if w, ok := pm["widgets"]; ok {
							extract(w)
						}
					}
				case map[string]interface{}:
					if w, ok := p["widgets"]; ok {
						extract(w)
					}
				}

				// 处理 widgets：维度限制 + 避免重叠 + 重新计算坐标
				// Zabbix 5.x 是单页结构，需要确保 widgets 不重叠且坐标在范围内
				occupiedCells := make(map[string]bool) // 记录已占用的单元格 "x,y"
				gridCols := 24                         // Zabbix 5.x 网格列数
				gridRows := 64                         // Zabbix 5.x 网格行数

				// 重新计算布局：按列分组，每列垂直排列
				// 统计不同 x 值的数量（列数）
				xValues := make(map[int]bool)
				for _, wi := range allWidgets {
					if wm, ok := wi.(map[string]interface{}); ok {
						if xv, ok := wm["x"]; ok {
							if n, ok2 := toInt(xv); ok2 {
								xValues[n] = true
							}
						}
					}
				}

				// 计算新的列数和 widget 尺寸
				numCols := len(xValues)
				if numCols == 0 {
					numCols = 1
				}
				if numCols > gridCols {
					numCols = gridCols
				}

				widgetWidth := gridCols / numCols
				if widgetWidth < 1 {
					widgetWidth = 1
				}

				// 按列重新排列 widgets
				columnWidgets := make(map[int][]map[string]interface{})
				for _, wi := range allWidgets {
					if wm, ok := wi.(map[string]interface{}); ok {
						colIdx := 0
						if xv, ok := wm["x"]; ok {
							if n, ok2 := toInt(xv); ok2 {
								// 将原始 x 映射到列索引（假设原始使用 72 列计算）
								colIdx = n / (72 / numCols)
								if colIdx < 0 {
									colIdx = 0
								} else if colIdx >= numCols {
									colIdx = numCols - 1
								}
							}
						}
						columnWidgets[colIdx] = append(columnWidgets[colIdx], wm)
					}
				}

				// 重新计算所有 widgets 的位置
				newWidgets := make([]interface{}, 0)
				for colIdx := 0; colIdx < numCols; colIdx++ {
					widgetsInCol := columnWidgets[colIdx]
					// 按 y 排序
					for i := 0; i < len(widgetsInCol); i++ {
						for j := i + 1; j < len(widgetsInCol); j++ {
							y1, y2 := 0, 0
							if yv, ok := widgetsInCol[i]["y"]; ok {
								if n, ok2 := toInt(yv); ok2 {
									y1 = n
								}
							}
							if yv, ok := widgetsInCol[j]["y"]; ok {
								if n, ok2 := toInt(yv); ok2 {
									y2 = n
								}
							}
							if y1 > y2 {
								widgetsInCol[i], widgetsInCol[j] = widgetsInCol[j], widgetsInCol[i]
							}
						}
					}

					// 重新计算该列中每个 widget 的位置
					x := colIdx * widgetWidth
					for rowIdx, wm := range widgetsInCol {
						// 获取并限制维度
						height := 5
						if hv, ok := wm["height"]; ok {
							if n, ok2 := toInt(hv); ok2 {
								height = n
							}
						}

						// Clamp height to 1-32
						if height < 1 {
							height = 1
						} else if height > 32 {
							height = 32
						}
						wm["height"] = height

						// Clamp width to 1-24
						wm["width"] = widgetWidth

						// 计算 y 坐标（垂直排列）
						y := rowIdx * height
						if y > 63 {
							break // 超出最大行数，停止添加该列的剩余 widgets
						}

						// 检查并避免位置冲突
						conflict := false
						for dy := 0; dy < height; dy++ {
							for dx := 0; dx < widgetWidth; dx++ {
								key := fmt.Sprintf("%d,%d", x+dx, y+dy)
								if occupiedCells[key] {
									conflict = true
									break
								}
							}
							if conflict {
								break
							}
						}

						if conflict {
							// 找到新的位置（简单策略：顺序查找空位）
							found := false
							for newY := 0; newY <= gridRows-height && !found; newY++ {
								for newX := 0; newX <= gridCols-widgetWidth && !found; newX++ {
									valid := true
									for dy := 0; dy < height; dy++ {
										for dx := 0; dx < widgetWidth; dx++ {
											key := fmt.Sprintf("%d,%d", newX+dx, newY+dy)
											if occupiedCells[key] {
												valid = false
												break
											}
										}
										if !valid {
											break
										}
									}
									if valid {
										x, y = newX, newY
										found = true
									}
								}
							}
							if !found {
								// 如果找不到位置，跳过此 widget
								continue
							}
						}

						// 更新 widget 位置
						wm["x"] = x
						wm["y"] = y

						// 标记占用的单元格
						for dy := 0; dy < height; dy++ {
							for dx := 0; dx < widgetWidth; dx++ {
								key := fmt.Sprintf("%d,%d", x+dx, y+dy)
								occupiedCells[key] = true
							}
						}

						newWidgets = append(newWidgets, wm)
					}
				}

				adaptedParams["widgets"] = newWidgets
				delete(adaptedParams, "pages")
			}
		} else {
			// Zabbix 6.x+ 版本：x 范围 0-71，y 范围 0-6，超出时自动分页
			if pages, ok := adaptedParams["pages"]; ok {
				// helper to convert various numeric types to int
				toInt := func(v interface{}) (int, bool) {
					switch t := v.(type) {
					case int:
						return t, true
					case int32:
						return int(t), true
					case int64:
						return int(t), true
					case float64:
						return int(t), true
					case float32:
						return int(t), true
					case string:
						if n, err := strconv.Atoi(t); err == nil {
							return n, true
						}
					}
					return 0, false
				}

				// Process pages and split widgets that exceed y=6
				switch p := pages.(type) {
				case []interface{}:
					var newPages []interface{}
					for _, page := range p {
						if pm, ok := page.(map[string]interface{}); ok {
							newPages = append(newPages, pm)
						}
					}
					adaptedParams["pages"] = newPages
				case []map[string]interface{}:
					var newPages []interface{}
					for _, pm := range p {
						newPages = append(newPages, pm)
					}
					adaptedParams["pages"] = newPages
				case map[string]interface{}:
					adaptedParams["pages"] = []interface{}{p}
				}

				// Auto-pagination: split widgets that exceed y=6 into new pages
				if pagesArr, ok := adaptedParams["pages"].([]interface{}); ok {
					var allPages []interface{}
					pageCounter := 0

					for _, page := range pagesArr {
						if pm, ok := page.(map[string]interface{}); ok {
							if widgets, ok := pm["widgets"]; ok {
								// Extract all widgets from this page
								var allWidgets []interface{}
								switch ws := widgets.(type) {
								case []interface{}:
									allWidgets = ws
								case []map[string]interface{}:
									for _, w := range ws {
										allWidgets = append(allWidgets, w)
									}
								case map[string]interface{}:
									allWidgets = []interface{}{ws}
								}

								// Split widgets into pages based on y position
								pageWidgets := make([][]interface{}, 0)
								currentPage := make([]interface{}, 0)

								for _, wi := range allWidgets {
									if wm, ok := wi.(map[string]interface{}); ok {
										// Get widget y position and height
										y := 0
										if yv, ok := wm["y"]; ok {
											if n, ok2 := toInt(yv); ok2 {
												y = n
											}
										}

										height := 5 // default height
										if hv, ok := wm["height"]; ok {
											if n, ok2 := toInt(hv); ok2 {
												height = n
											}
										}

										// Clamp width to 1-72
										width := 72
										if wv, ok := wm["width"]; ok {
											if n, ok2 := toInt(wv); ok2 {
												if n < 1 {
													n = 1
												} else if n > 72 {
													n = 72
												}
												wm["width"] = n
												width = n
											}
										}

										// Clamp x to 0-(72-width)
										if xv, ok := wm["x"]; ok {
											if n, ok2 := toInt(xv); ok2 {
												maxX := 72 - width
												if n < 0 {
													n = 0
												} else if n > maxX {
													n = maxX
												}
												wm["x"] = n
											}
										}

										// Clamp y to 0-63
										if yv, ok := wm["y"]; ok {
											if n, ok2 := toInt(yv); ok2 {
												if n < 0 {
													n = 0
												} else if n > 63 {
													n = 63
												}
												wm["y"] = n
												y = n // Update y with clamped value
											}
										} // Clamp height to 1-64
										if hv, ok := wm["height"]; ok {
											if n, ok2 := toInt(hv); ok2 {
												if n < 1 {
													n = 1
												} else if n > 64 {
													n = 64
												}
												wm["height"] = n
												height = n // Update height with clamped value
											}
										}

										// Check if widget fits in current page (y <= 63)
										if y <= 63 && y+height-1 <= 63 {
											currentPage = append(currentPage, wm)
										} else {
											// Start a new page
											if len(currentPage) > 0 {
												pageWidgets = append(pageWidgets, currentPage)
												currentPage = make([]interface{}, 0)
											}
											// Adjust y position for new page
											wm["y"] = 0
											currentPage = append(currentPage, wm)
										}
									}
								}

								if len(currentPage) > 0 {
									pageWidgets = append(pageWidgets, currentPage)
								}

								// Create new pages with the split widgets
								for _, widgets := range pageWidgets {
									newPage := make(map[string]interface{})
									for k, v := range pm {
										if k != "widgets" {
											newPage[k] = v
										}
									}
									newPage["widgets"] = widgets
									allPages = append(allPages, newPage)
									pageCounter++
								}
							}
						}
					}
					adaptedParams["pages"] = allPages
				}
			}
		}
	}

	return adaptedParams
}
