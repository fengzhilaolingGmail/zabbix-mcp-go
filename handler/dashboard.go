/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2026-01-06 18:59:21
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2026-01-08 17:29:55
 * @FilePath: \zabbix-mcp-go\handler\dashboard.go
 * @Description: 仪表盘相关功能
 * @Copyright: Copyright (c) 2025 by fengzhilaoling@gmail.com, All Rights Reserved.
 */
package handler

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"zabbixMcp/logger"
	"zabbixMcp/models"
	"zabbixMcp/server"
	"zabbixMcp/zabbix"

	"github.com/mark3labs/mcp-go/mcp"
)

// buildDashboardSpecFromArgs 将 mcp 参数（map[string]interface{}）解析为 models.DashboardParams
// 包含校验：
//   - widget.type 使用 models.WidgetType.IsValid() 验证，非法则跳过该 widget
//   - widget.fields 中 type 必须在 0-13 范围内，否则跳过该 field
func buildDashboardSpecFromArgs(args map[string]interface{}) (models.DashboardParams, error) {
	if args == nil {
		return models.DashboardParams{}, fmt.Errorf("args is nil")
	}

	spec := models.DashboardParams{}

	if v, ok := args["name"].(string); ok {
		spec.Name = v
	}
	if v, ok := args["userid"].(string); ok {
		spec.UserID = v
	}
	if v, ok := args["private"].(float64); ok {
		n := int(v)
		spec.Private = &n
	} else if v, ok := args["private"].(int); ok {
		n := v
		spec.Private = &n
	}
	if v, ok := args["display_period"].(float64); ok {
		n := int(v)
		spec.DisplayPeriod = &n
	} else if v, ok := args["display_period"].(int); ok {
		n := v
		spec.DisplayPeriod = &n
	}
	if v, ok := args["auto_start"].(float64); ok {
		n := int(v)
		spec.AutoStart = &n
	} else if v, ok := args["auto_start"].(int); ok {
		n := v
		spec.AutoStart = &n
	}

	// parse users
	if arr, ok := args["users"].([]interface{}); ok {
		users := make([]models.DashboardUser, 0, len(arr))
		for _, u := range arr {
			if um, ok := u.(map[string]interface{}); ok {
				du := models.DashboardUser{}
				if id, ok2 := um["userid"].(string); ok2 {
					du.UserID = id
				}
				if p, ok2 := um["permission"].(float64); ok2 {
					du.Permission = int(p)
				}
				users = append(users, du)
			}
		}
		spec.Users = users
	}

	// parse userGroups
	if arr, ok := args["userGroups"].([]interface{}); ok {
		ugs := make([]models.DashboardUserGroup, 0, len(arr))
		for _, g := range arr {
			if gm, ok := g.(map[string]interface{}); ok {
				dug := models.DashboardUserGroup{}
				if id, ok2 := gm["usrgrpid"].(string); ok2 {
					dug.Usrgrpid = id
				}
				if p, ok2 := gm["permission"].(float64); ok2 {
					dug.Permission = int(p)
				}
				ugs = append(ugs, dug)
			}
		}
		spec.UserGroups = ugs
	}

	// parse pages
	if parr, ok := args["pages"].([]interface{}); ok {
		pages := make([]models.DashboardPage, 0, len(parr))
		for _, p := range parr {
			if pm, ok := p.(map[string]interface{}); ok {
				page := models.DashboardPage{}
				if v, ok2 := pm["name"].(string); ok2 {
					page.Name = v
				}
				if v, ok2 := pm["display_period"].(float64); ok2 {
					n := int(v)
					page.DisplayPeriod = &n
				} else if v, ok2 := pm["display_period"].(int); ok2 {
					n := v
					page.DisplayPeriod = &n
				}

				// widgets
				if wArr, ok2 := pm["widgets"].([]interface{}); ok2 {
					widgets := make([]models.DashboardWidget, 0, len(wArr))
					for _, w := range wArr {
						if wm, ok3 := w.(map[string]interface{}); ok3 {
							widget := models.DashboardWidget{}
							if v, ok4 := wm["type"].(string); ok4 {
								widget.Type = models.WidgetType(v)
							}
							if !widget.Type.IsValid() {
								logger.L().Warnf("跳过未知 widget type: %v", widget.Type)
								continue
							}
							if v, ok4 := wm["name"].(string); ok4 {
								widget.Name = v
							}
							if v, ok4 := wm["x"].(float64); ok4 {
								n := int(v)
								widget.X = &n
							}
							if v, ok4 := wm["y"].(float64); ok4 {
								n := int(v)
								widget.Y = &n
							}
							if v, ok4 := wm["width"].(float64); ok4 {
								n := int(v)
								widget.Width = &n
							}
							if v, ok4 := wm["height"].(float64); ok4 {
								n := int(v)
								widget.Height = &n
							}
							if v, ok4 := wm["view_mode"].(float64); ok4 {
								n := int(v)
								widget.ViewMode = &n
							}

							// fields
							if fArr, ok4 := wm["fields"].([]interface{}); ok4 {
								fields := make([]models.DashboardWidgetField, 0, len(fArr))
								for _, f := range fArr {
									if fm, ok5 := f.(map[string]interface{}); ok5 {
										df := models.DashboardWidgetField{}
										if t, ok6 := fm["type"].(float64); ok6 {
											df.Type = int(t)
										}
										if !models.IsValidWidgetFieldType(df.Type) {
											logger.L().Warnf("跳过无效 widget field.type: %d", df.Type)
											continue
										}
										if n, ok6 := fm["name"].(string); ok6 {
											df.Name = n
										}
										if val, ok6 := fm["value"]; ok6 {
											df.Value = val
										}
										fields = append(fields, df)
									}
								}
								if len(fields) > 0 {
									widget.Fields = fields
								}
							}
							widgets = append(widgets, widget)
						}
					}
					if len(widgets) > 0 {
						page.Widgets = widgets
					}
				}
				pages = append(pages, page)
			}
		}
		spec.Pages = pages
	}

	// if top-level pages not provided but single page provided as map using 'page' key
	if spec.Name == "" && args["name"] != nil {
		// nothing special; leave empty
	}

	return spec, nil
}

// CreateDashboardHandler 通过注入的 ClientProvider 调用 dashboard.create 并返回结果
func CreateDashboardHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	instanceName := ""
	if args, ok := req.Params.Arguments.(map[string]interface{}); ok {
		if v, ok2 := args["instance"].(string); ok2 {
			instanceName = v
		}
	}

	var argsMap map[string]interface{}
	if a, ok := req.Params.Arguments.(map[string]interface{}); ok {
		argsMap = a
	}

	spec, err := buildDashboardSpecFromArgs(argsMap)
	if err != nil {
		logger.L().Warnf("解析 dashboard 参数失败: %v", err)
		return nil, fmt.Errorf("解析 dashboard 参数失败: %w", err)
	}

	// override users/userGroups if explicit arrays passed (maintain backward compat)
	if argsMap != nil {
		if arr, ok := argsMap["users"].([]interface{}); ok && len(arr) > 0 {
			users := make([]models.DashboardUser, 0, len(arr))
			for _, u := range arr {
				if um, ok := u.(map[string]interface{}); ok {
					du := models.DashboardUser{}
					if id, ok2 := um["userid"].(string); ok2 {
						du.UserID = id
					}
					if p, ok2 := um["permission"].(float64); ok2 {
						du.Permission = int(p)
					}
					users = append(users, du)
				}
			}
			spec.Users = users
		}
		if arr, ok := argsMap["userGroups"].([]interface{}); ok && len(arr) > 0 {
			ugs := make([]models.DashboardUserGroup, 0, len(arr))
			for _, g := range arr {
				if gm, ok := g.(map[string]interface{}); ok {
					dug := models.DashboardUserGroup{}
					if id, ok2 := gm["usrgrpid"].(string); ok2 {
						dug.Usrgrpid = id
					}
					if p, ok2 := gm["permission"].(float64); ok2 {
						dug.Permission = int(p)
					}
					ugs = append(ugs, dug)
				}
			}
			spec.UserGroups = ugs
		}
	}

	if clientPool == nil {
		return mcp.NewToolResultStructuredOnly(makeResult(map[string]interface{}{})), nil
	}

	// 批量解析 widget 内传入的 hostnames（如果有），并把值替换为数字 hostid 或数字数组
	// 收集需要解析的主机名
	hostNamesSet := map[string]struct{}{}
	for pi := range spec.Pages {
		pg := &spec.Pages[pi]
		for wi := range pg.Widgets {
			w := &pg.Widgets[wi]
			for fi := range w.Fields {
				f := &w.Fields[fi]
				if f.Name == "hostids" {
					switch v := f.Value.(type) {
					case string:
						if v != "" {
							hostNamesSet[v] = struct{}{}
						}
					case []interface{}:
						for _, e := range v {
							if s, ok := e.(string); ok && s != "" {
								hostNamesSet[s] = struct{}{}
							}
						}
					}
				}
			}
		}
	}

	if len(hostNamesSet) > 0 {
		hostNames := make([]string, 0, len(hostNamesSet))
		for hn := range hostNamesSet {
			hostNames = append(hostNames, hn)
		}
		specHosts := models.HostParams{Output: []string{"hostid", "host"}}
		specHosts.Filter = map[string]interface{}{"host": hostNames}
		hostList, err := server.GetHosts(ctx, clientPool, specHosts, instanceName)
		hostNameToID := map[string]string{}
		if err == nil {
			for _, h := range hostList {
				hn := ""
				hid := ""
				if v, ok := h["host"].(string); ok {
					hn = v
				}
				if v, ok := h["hostid"].(string); ok {
					hid = v
				}
				if hn != "" && hid != "" {
					hostNameToID[hn] = hid
				}
			}
		} else {
			logger.L().Warnf("批量解析 widget hostnames 失败: %v", err)
		}

		// 替换字段值
		for pi := range spec.Pages {
			pg := &spec.Pages[pi]
			for wi := range pg.Widgets {
				w := &pg.Widgets[wi]
				for fi := range w.Fields {
					f := &w.Fields[fi]
					if f.Name != "hostids" {
						continue
					}
					switch v := f.Value.(type) {
					case string:
						if v == "" {
							continue
						}
						// 如果是数字字符串直接转换
						if idNum, err := strconv.Atoi(v); err == nil {
							f.Value = idNum
							continue
						}
						if hid, ok := hostNameToID[v]; ok {
							if idNum, err := strconv.Atoi(hid); err == nil {
								f.Value = idNum
							} else {
								logger.L().Warnf("解析到的 hostid %s 不是数字，保留原值", hid)
							}
						} else {
							logger.L().Warnf("未解析到主机 %s 的 hostid，保留原值", v)
						}
					case []interface{}:
						outNums := make([]interface{}, 0, len(v))
						for _, e := range v {
							switch ev := e.(type) {
							case string:
								if ev == "" {
									continue
								}
								if idNum, err := strconv.Atoi(ev); err == nil {
									outNums = append(outNums, idNum)
									continue
								}
								if hid, ok := hostNameToID[ev]; ok {
									if idNum, err := strconv.Atoi(hid); err == nil {
										outNums = append(outNums, idNum)
									} else {
										logger.L().Warnf("解析到的 hostid %s 不是数字，跳过", hid)
									}
								} else {
									logger.L().Warnf("未解析到主机 %s 的 hostid，跳过", ev)
								}
							case float64:
								outNums = append(outNums, int(ev))
							case int:
								outNums = append(outNums, ev)
							}
						}
						if len(outNums) == 1 {
							f.Value = outNums[0]
						} else if len(outNums) > 1 {
							f.Value = outNums
						} else {
							logger.L().Warnf("widget %s 中 hostids 未解析出有效 id，保留原值", w.Name)
						}
					}
				}
			}
		}
	}

	dashboard, err := server.CreateDashboard(ctx, clientPool, spec, instanceName)
	if err != nil {
		logger.L().Errorf("调用 dashboard.create 失败: %v", err)
		return nil, fmt.Errorf("调用 dashboard.create 失败: %w", err)
	}
	return mcp.NewToolResultStructuredOnly(makeResult(dashboard)), nil
}

// CreateGraphDashboardHandler 自动创建图形仪表盘，支持两种模式
func CreateGraphDashboardHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	instanceName := ""
	name := ""
	mode := 0
	hosts := []string{}
	graphNames := []string{}
	maxCols := 0
	maxRows := 10

	if args, ok := req.Params.Arguments.(map[string]interface{}); ok {
		if v, ok2 := args["instance"].(string); ok2 {
			instanceName = v
		}
		if v, ok2 := args["name"].(string); ok2 {
			name = v
		}
		if v, ok2 := args["mode"].(float64); ok2 {
			mode = int(v)
		}
		if arr, ok2 := args["hosts"].([]interface{}); ok2 {
			for _, h := range arr {
				if s, ok3 := h.(string); ok3 && s != "" {
					hosts = append(hosts, s)
				}
			}
		}
		if arr, ok2 := args["graph_names"].([]interface{}); ok2 {
			for _, g := range arr {
				if s, ok3 := g.(string); ok3 && s != "" {
					graphNames = append(graphNames, s)
				}
			}
		}
		if v, ok2 := args["max_cols"].(float64); ok2 && int(v) > 0 {
			maxCols = int(v)
		}
		if v, ok2 := args["max_rows"].(float64); ok2 && int(v) > 0 {
			maxRows = int(v)
		}
	}

	if clientPool == nil {
		return mcp.NewToolResultStructuredOnly(makeResult(map[string]interface{}{})), nil
	}

	// 查询主机ID映射
	hostNameToID := map[string]string{}
	if len(hosts) > 0 {
		specHosts := models.HostParams{Output: []string{"hostid", "host"}}
		specHosts.Filter = map[string]interface{}{
			"host":   hosts,
			"status": "0", // 只查询主机，排除模板（模板 status = 3）
		}
		hostList, err := server.GetHosts(ctx, clientPool, specHosts, instanceName)
		if err == nil {
			for _, h := range hostList {
				if hn, ok := h["host"].(string); ok {
					if hid, ok := h["hostid"].(string); ok {
						hostNameToID[hn] = hid
					}
				}
			}
		}
	}

	var widgets []models.DashboardWidget

	switch mode {
	case 1:
		// 模式1：聚合模式 - 多台主机相同图形名称
		widgets = buildAggregateModeWidgets(hosts, hostNameToID, graphNames, instanceName, ctx, maxCols, maxRows)
	case 2:
		// 模式2：分列模式 - 每台主机不同图形
		widgets = buildSeparateModeWidgets(hosts, hostNameToID, graphNames, instanceName, ctx, maxCols, maxRows)
	default:
		return nil, fmt.Errorf("无效的模式: %d，请使用 1 (聚合模式) 或 2 (分列模式)", mode)
	}

	page := models.DashboardPage{
		Name:    "Default Page",
		Widgets: widgets,
	}

	spec := models.DashboardParams{
		Name:  name,
		Pages: []models.DashboardPage{page},
	}

	dashboard, err := server.CreateDashboard(ctx, clientPool, spec, instanceName)
	if err != nil {
		logger.L().Errorf("调用 dashboard.create 失败: %v", err)
		return nil, fmt.Errorf("调用 dashboard.create 失败: %w", err)
	}
	return mcp.NewToolResultStructuredOnly(makeResult(dashboard)), nil
}

// buildAggregateModeWidgets 构建聚合模式的 widgets
func buildAggregateModeWidgets(hosts []string, hostNameToID map[string]string, graphNames []string, instanceName string, ctx context.Context, maxCols, maxRows int) []models.DashboardWidget {
	widgets := make([]models.DashboardWidget, 0)

	if len(hosts) == 0 || len(graphNames) == 0 {
		return widgets
	}

	// 默认列数为2
	if maxCols <= 0 {
		maxCols = 2
	}

	// 检查是否超限
	if maxCols > 72 {
		maxCols = 72
	}
	if maxRows > 64 {
		maxRows = 64
	}

	// 自动计算 widget 尺寸
	widgetWidth := 72 / maxCols
	if widgetWidth < 1 {
		widgetWidth = 1
	}
	widgetHeight := 64 / maxRows
	if widgetHeight < 1 {
		widgetHeight = 1
	}

	// 一次性查询所有主机和所有图形的映射关系
	hostGraphMap := make(map[string]map[string]string) // host -> graphName -> graphid

	// 收集所有需要查询的主机ID
	hostIDs := make([]string, 0, len(hosts))
	for _, host := range hosts {
		if hidStr, ok := hostNameToID[host]; ok {
			hostIDs = append(hostIDs, hidStr)
			hostGraphMap[host] = make(map[string]string)
		}
	}

	if len(hostIDs) > 0 && len(graphNames) > 0 {
		// 使用 host.get 批量查询所有主机及其图形信息
		// 通过 selectGraphs 参数一次性获取主机的图形列表
		// 注意：需要排除模板（template），只查询真实主机（status = 0）
		specHosts := models.HostParams{
			Output: []string{"hostid", "host"},
			Filter: map[string]interface{}{
				"hostid": hostIDs,
				"status": "0", // 只查询主机，排除模板（模板 status = 3）
			},
			SelectGraphs: []string{"graphid", "name"}, // 获取图形信息
		}

		var lease zabbix.ClientLease
		var err error
		if instanceName != "" {
			lease, err = clientPool.AcquireByInstance(ctx, instanceName)
		} else {
			lease, err = clientPool.Acquire(ctx)
		}
		if err == nil && lease != nil {
			client := lease.Client()
			adapted := client.AdaptAPIParams("host.get", specHosts)
			var hostList []map[string]interface{}
			callErr := client.Call(ctx, "host.get", adapted, &hostList)
			lease.Release(callErr)

			if callErr == nil {
				logger.L().Infof("批量查询到 %d 个主机及其图形信息", len(hostList))
				// 构建映射关系
				for _, host := range hostList {
					hostname, _ := host["host"].(string)

					// 获取该主机的图形列表
					if graphs, ok := host["graphs"].([]interface{}); ok && len(graphs) > 0 {
						for _, g := range graphs {
							if graphMap, ok := g.(map[string]interface{}); ok {
								graphid, _ := graphMap["graphid"].(string)
								graphName, _ := graphMap["name"].(string)

								// 检查图形名称是否在请求的列表中（支持模糊匹配）
								for _, requestedName := range graphNames {
									if strings.Contains(graphName, requestedName) {
										hostGraphMap[hostname][requestedName] = graphid
										logger.L().Debugf("映射: %s -> %s -> %s", hostname, requestedName, graphid)
										break
									}
								}
							}
						}
					}
				}
			} else {
				logger.L().Warnf("批量查询主机及图形失败: %v", callErr)
			}
		}
	}

	idx := 0
	for _, graphName := range graphNames {
		for _, host := range hosts {
			if hidStr, ok := hostNameToID[host]; ok {
				// 将 hostid 转换为数字
				hidNum, err := strconv.Atoi(hidStr)
				if err != nil {
					logger.L().Warnf("主机 %s 的 hostid %s 不是有效数字，跳过", host, hidStr)
					continue
				}

				// 从映射中获取图形ID
				graphid := ""
				if hostMap, ok := hostGraphMap[host]; ok {
					graphid = hostMap[graphName]
				}

				var widgetType models.WidgetType
				var fields []models.DashboardWidgetField

				if graphid != "" {
					// 图形存在，创建 graph widget
					widgetType = models.WidgetTypeGraph
					fields = []models.DashboardWidgetField{
						{Type: 6, Name: "graphid", Value: graphid},
						{Type: 3, Name: "hostids", Value: hidNum},
					}
				} else {
					// 图形不存在，创建 text widget 作为占位符
					widgetType = models.WidgetTypeText
					textContent := fmt.Sprintf("图形不存在: %s\n主机: %s", graphName, host)
					fields = []models.DashboardWidgetField{
						{Type: 0, Name: "text_size", Value: 12},     // text_size 字段类型为 0（整数）
						{Type: 1, Name: "text", Value: textContent}, // text 字段类型为 1（字符串）
					}
					logger.L().Debugf("主机 %s 的图形 %s 不存在，创建文本占位组件", host, graphName)
				}

				x := (idx % maxCols) * widgetWidth
				y := (idx / maxCols) * widgetHeight

				w := models.DashboardWidget{
					Type: widgetType,
					Name: fmt.Sprintf("%s - %s", host, graphName),
					X:    &x,
					Y:    &y,
				}
				w.Width = &widgetWidth
				w.Height = &widgetHeight
				w.Fields = fields
				widgets = append(widgets, w)
				idx++
			}
		}
	}
	return widgets
}

// buildSeparateModeWidgets 构建分列模式的 widgets
// 模式2：每台主机一列，每列包含该主机的多个图形（垂直排列）
func buildSeparateModeWidgets(hosts []string, hostNameToID map[string]string, graphNames []string, instanceName string, ctx context.Context, maxCols, maxRows int) []models.DashboardWidget {
	widgets := make([]models.DashboardWidget, 0)

	if len(hosts) == 0 {
		return widgets
	}

	// 默认列数为主机数量
	if maxCols <= 0 {
		maxCols = len(hosts)
	}

	// 检查是否超限
	if maxCols > 72 {
		return widgets // 主机数量超限，返回空
	}
	if maxRows > 64 {
		maxRows = 64
	}

	// 自动计算 widget 尺寸
	widgetWidth := 72 / maxCols
	if widgetWidth < 1 {
		widgetWidth = 1
	}
	widgetHeight := 64 / maxRows
	if widgetHeight < 1 {
		widgetHeight = 1
	}

	// 一次性查询所有主机和所有图形的映射关系
	hostGraphMap := make(map[string]map[string]string) // host -> graphName -> graphid

	// 收集所有需要查询的主机ID
	hostIDs := make([]string, 0, len(hosts))
	hostIDToName := make(map[string]string) // hostid -> hostname

	for _, host := range hosts {
		if hidStr, ok := hostNameToID[host]; ok {
			hostIDs = append(hostIDs, hidStr)
			hostIDToName[hidStr] = host
			hostGraphMap[host] = make(map[string]string)
		}
	}

	if len(hostIDs) > 0 && len(graphNames) > 0 {
		// 使用 host.get 批量查询所有主机及其图形信息
		// 通过 selectGraphs 参数一次性获取主机的图形列表
		// 注意：需要排除模板（template），只查询真实主机（status = 0）
		specHosts := models.HostParams{
			Output: []string{"hostid", "host"},
			Filter: map[string]interface{}{
				"hostid": hostIDs,
				"status": "0", // 只查询主机，排除模板（模板 status = 3）
			},
			SelectGraphs: []string{"graphid", "name"}, // 获取图形信息
		}

		var lease zabbix.ClientLease
		var err error
		if instanceName != "" {
			lease, err = clientPool.AcquireByInstance(ctx, instanceName)
		} else {
			lease, err = clientPool.Acquire(ctx)
		}
		if err == nil && lease != nil {
			client := lease.Client()
			adapted := client.AdaptAPIParams("host.get", specHosts)
			var hostList []map[string]interface{}
			callErr := client.Call(ctx, "host.get", adapted, &hostList)
			lease.Release(callErr)

			if callErr == nil {
				logger.L().Infof("批量查询到 %d 个主机及其图形信息", len(hostList))
				// 构建映射关系
				for _, host := range hostList {
					hostname, _ := host["host"].(string)

					// 获取该主机的图形列表
					if graphs, ok := host["graphs"].([]interface{}); ok && len(graphs) > 0 {
						for _, g := range graphs {
							if graphMap, ok := g.(map[string]interface{}); ok {
								graphid, _ := graphMap["graphid"].(string)
								graphName, _ := graphMap["name"].(string)

								// 检查图形名称是否在请求的列表中（支持模糊匹配）
								for _, requestedName := range graphNames {
									if strings.Contains(graphName, requestedName) {
										hostGraphMap[hostname][requestedName] = graphid
										logger.L().Debugf("映射: %s -> %s -> %s", hostname, requestedName, graphid)
										break
									}
								}
							}
						}
					}
				}
			} else {
				logger.L().Warnf("批量查询主机及图形失败: %v", callErr)
			}
		}
	} // 按主机分列，每个主机一列，每列包含多个图形（垂直排列）
	for hostIdx, host := range hosts {
		if hostIdx >= maxCols {
			break // 超过最大列数，停止
		}

		if hidStr, ok := hostNameToID[host]; ok {
			// 将 hostid 转换为数字
			hidNum, err := strconv.Atoi(hidStr)
			if err != nil {
				logger.L().Warnf("主机 %s 的 hostid %s 不是有效数字，跳过", host, hidStr)
				continue
			}

			// 该列的 x 坐标（固定）
			x := hostIdx * widgetWidth

			// 为该主机创建所有图形，垂直排列
			for graphIdx, graphName := range graphNames {
				if graphIdx >= maxRows {
					break // 超过最大行数，停止
				}

				// 从映射中获取图形ID
				graphid := ""
				if hostMap, ok := hostGraphMap[host]; ok {
					graphid = hostMap[graphName]
				}

				// y 坐标根据图形索引递增（垂直排列）
				y := graphIdx * widgetHeight

				w := models.DashboardWidget{
					Name: fmt.Sprintf("%s - %s", host, graphName),
					X:    &x,
					Y:    &y,
				}
				w.Width = &widgetWidth
				w.Height = &widgetHeight

				var fields []models.DashboardWidgetField

				if graphid != "" {
					// 图形存在，创建 graph widget
					w.Type = models.WidgetTypeGraph
					fields = []models.DashboardWidgetField{
						{Type: 6, Name: "graphid", Value: graphid},
						{Type: 3, Name: "hostids", Value: hidNum},
					}
				} else {
					// 图形不存在，创建 text widget 作为占位符
					w.Type = models.WidgetTypeText
					textContent := fmt.Sprintf("图形不存在: %s\n主机: %s", graphName, host)
					fields = []models.DashboardWidgetField{
						{Type: 0, Name: "text_size", Value: 12},     // text_size 字段类型为 0（整数）
						{Type: 1, Name: "text", Value: textContent}, // text 字段类型为 1（字符串）
					}
					logger.L().Debugf("主机 %s 的图形 %s 不存在，创建文本占位组件", host, graphName)
				}

				w.Fields = fields
				widgets = append(widgets, w)
			}
		}
	}
	return widgets
} // queryGraphIDByName 根据主机和图形名称查询 graphid
func queryGraphIDByName(ctx context.Context, hostname, graphName, instanceName string) string {
	if clientPool == nil {
		return ""
	}

	// 先获取主机ID
	specHosts := models.HostParams{Output: []string{"hostid"}}
	specHosts.Filter = map[string]interface{}{"host": []string{hostname}}
	hostList, err := server.GetHosts(ctx, clientPool, specHosts, instanceName)
	if err != nil || len(hostList) == 0 {
		return ""
	}

	hostid := ""
	if v, ok := hostList[0]["hostid"].(string); ok {
		hostid = v
	}

	if hostid == "" {
		return ""
	}

	// 查询该主机下的图形
	specGraphs := models.GraphParams{
		Output: []string{"graphid", "name"},
		Filter: map[string]interface{}{"name": graphName},
	}
	// 添加 hostids 过滤
	specGraphs.Filter["hostids"] = hostid

	var lease zabbix.ClientLease
	if instanceName != "" {
		lease, err = clientPool.AcquireByInstance(ctx, instanceName)
	} else {
		lease, err = clientPool.Acquire(ctx)
	}
	if err != nil || lease == nil {
		return ""
	}

	client := lease.Client()
	adapted := client.AdaptAPIParams("graph.get", specGraphs)
	var graphList []map[string]interface{}
	callErr := client.Call(ctx, "graph.get", adapted, &graphList)
	lease.Release(callErr)

	if callErr != nil || len(graphList) == 0 {
		return ""
	}

	if v, ok := graphList[0]["graphid"].(string); ok {
		return v
	}

	return ""
}

// queryGraphNameByID 根据 graphid 查询图形名称
func queryGraphNameByID(ctx context.Context, graphid, instanceName string) string {
	if clientPool == nil {
		return ""
	}

	specGraphs := models.GraphParams{
		GraphIDs: []string{graphid},
		Output:   []string{"graphid", "name"},
	}

	var lease zabbix.ClientLease
	var err error
	if instanceName != "" {
		lease, err = clientPool.AcquireByInstance(ctx, instanceName)
	} else {
		lease, err = clientPool.Acquire(ctx)
	}
	if err != nil || lease == nil {
		return ""
	}

	client := lease.Client()
	adapted := client.AdaptAPIParams("graph.get", specGraphs)
	var graphList []map[string]interface{}
	callErr := client.Call(ctx, "graph.get", adapted, &graphList)
	lease.Release(callErr)

	if callErr != nil || len(graphList) == 0 {
		return ""
	}

	if v, ok := graphList[0]["name"].(string); ok {
		return v
	}

	return ""
}
