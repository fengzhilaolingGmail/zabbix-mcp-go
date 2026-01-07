/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2026-01-06 18:59:21
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2026-01-07 12:02:46
 * @FilePath: \zabbix-mcp-go\handler\dashboard.go
 * @Description: 仪表盘相关功能
 * @Copyright: Copyright (c) 2025 by fengzhilaoling@gmail.com, All Rights Reserved.
 */
package handler

import (
	"context"
	"fmt"
	"strconv"

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

// CreateGraphDashboardHandler 使用 hosts + graphids 自动创建布局良好的仪表盘
func CreateGraphDashboardHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	instanceName := ""
	name := ""
	hosts := []string{}
	graphidsRaw := []interface{}{}
	cols := 2
	widgetWidth := 36
	widgetHeight := 5
	var private *int
	var displayPeriod *int
	var autoStart *int

	if args, ok := req.Params.Arguments.(map[string]interface{}); ok {
		if v, ok2 := args["instance"].(string); ok2 {
			instanceName = v
		}
		if v, ok2 := args["name"].(string); ok2 {
			name = v
		}
		if arr, ok2 := args["hosts"].([]interface{}); ok2 {
			for _, h := range arr {
				if s, ok3 := h.(string); ok3 && s != "" {
					hosts = append(hosts, s)
				}
			}
		}
		if arr, ok2 := args["graphids"].([]interface{}); ok2 {
			graphidsRaw = arr
		}
		if v, ok2 := args["cols"].(float64); ok2 && int(v) > 0 {
			cols = int(v)
		}
		if v, ok2 := args["widgetWidth"].(float64); ok2 && int(v) > 0 {
			widgetWidth = int(v)
		}
		if v, ok2 := args["widgetHeight"].(float64); ok2 && int(v) > 0 {
			widgetHeight = int(v)
		}
		if v, ok2 := args["private"].(float64); ok2 {
			n := int(v)
			private = &n
		}
		if v, ok2 := args["display_period"].(float64); ok2 {
			n := int(v)
			displayPeriod = &n
		}
		if v, ok2 := args["auto_start"].(float64); ok2 {
			n := int(v)
			autoStart = &n
		}
	}

	if clientPool == nil {
		return mcp.NewToolResultStructuredOnly(makeResult(map[string]interface{}{})), nil
	}

	// Prepare graph id list as strings (models will normalize)
	graphids := make([]string, 0, len(graphidsRaw))
	for _, g := range graphidsRaw {
		switch v := g.(type) {
		case string:
			if v != "" {
				graphids = append(graphids, v)
			}
		case float64:
			graphids = append(graphids, strconv.Itoa(int(v)))
		case int:
			graphids = append(graphids, strconv.Itoa(v))
		default:
			graphids = append(graphids, fmt.Sprintf("%v", v))
		}
	}

	pairwise := false
	if len(hosts) > 0 && len(hosts) == len(graphids) {
		pairwise = true
	}

	// Batch-resolve graph names for nicer widget labels
	graphIDToName := map[string]string{}
	if len(graphids) > 0 {
		if clientPool != nil {
			var lease zabbix.ClientLease
			var err error
			if instanceName != "" {
				lease, err = clientPool.AcquireByInstance(ctx, instanceName)
			} else {
				lease, err = clientPool.Acquire(ctx)
			}
			if err == nil && lease != nil {
				client := lease.Client()
				specGraph := models.GraphParams{GraphIDs: graphids, Output: []string{"graphid", "name"}}
				adapted := client.AdaptAPIParams("graph.get", specGraph)
				var graphList []map[string]interface{}
				callErr := client.Call(ctx, "graph.get", adapted, &graphList)
				lease.Release(callErr)
				if callErr == nil {
					for _, g := range graphList {
						gid := ""
						gname := ""
						if v, ok := g["graphid"].(string); ok {
							gid = v
						}
						if v, ok := g["name"].(string); ok {
							gname = v
						}
						if gid != "" && gname != "" {
							graphIDToName[gid] = gname
						}
					}
				} else {
					logger.L().Warnf("查询 graph 名称失败: %v", callErr)
				}
			} else if err != nil {
				logger.L().Warnf("无法获取 zabbix client 来查询 graph 名称: %v", err)
			}
		}
	}

	// Always attempt to batch-resolve hosts (not only in pairwise mode)
	hostNameToID := map[string]string{}
	if len(hosts) > 0 {
		specHosts := models.HostParams{Output: []string{"hostid", "host"}}
		specHosts.Filter = map[string]interface{}{"host": hosts}
		hostList, err := server.GetHosts(ctx, clientPool, specHosts, instanceName)
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
			logger.L().Warnf("查询主机 hostid 失败，将不注入 hostids：%v", err)
		}
	}

	widgets := make([]models.DashboardWidget, 0, len(graphids))
	for idx, gid := range graphids {
		// Prefer using graph name when available (fallback to id)
		nameOrID := gid
		if n, ok := graphIDToName[gid]; ok && n != "" {
			nameOrID = n
		}
		w := models.DashboardWidget{
			Type: models.WidgetTypeGraph,
			Name: nameOrID,
		}
		{
			ww := widgetWidth
			hh := widgetHeight
			w.Width = &ww
			w.Height = &hh
		}
		xv := (idx % cols) * widgetWidth
		yv := (idx / cols) * widgetHeight
		w.X = &xv
		w.Y = &yv

		var graphVal interface{} = gid
		if n, err := strconv.Atoi(gid); err == nil {
			graphVal = n
		}
		// graphid: use field type for graph (6) per Zabbix dashboard widget fields
		fields := []models.DashboardWidgetField{{
			Type:  6,
			Name:  "graphid",
			Value: graphVal,
		}}
		if pairwise {
			h := hosts[idx]
			w.Name = fmt.Sprintf("%s - %s", h, nameOrID)
			if hid, ok := hostNameToID[h]; ok && hid != "" {
				// hostids: use field type for host (3) and pass array of numeric ids
				if idNum, err := strconv.Atoi(hid); err == nil {
					// Zabbix expects a single numeric host id value for this field (not an array)
					fields = append(fields, models.DashboardWidgetField{Type: 3, Name: "hostids", Value: idNum})
				} else {
					logger.L().Warnf("主机 id %s 不是数字，跳过注入 hostids", hid)
				}
			} else {
				logger.L().Warnf("未解析到主机 %s 的 hostid，跳过注入 hostids", h)
			}
		}
		w.Fields = fields
		widgets = append(widgets, w)
	}

	dp := 0
	if displayPeriod != nil {
		dp = *displayPeriod
	}
	page := models.DashboardPage{
		Name:          "Default Page",
		DisplayPeriod: &dp,
		Widgets:       widgets,
	}

	spec := models.DashboardParams{
		Name:  name,
		Pages: []models.DashboardPage{page},
	}
	if private != nil {
		spec.Private = private
	}
	if displayPeriod != nil {
		spec.DisplayPeriod = displayPeriod
	}
	if autoStart != nil {
		spec.AutoStart = autoStart
	}

	dashboard, err := server.CreateDashboard(ctx, clientPool, spec, instanceName)
	if err != nil {
		logger.L().Errorf("调用 dashboard.create 失败: %v", err)
		return nil, fmt.Errorf("调用 dashboard.create 失败: %w", err)
	}
	return mcp.NewToolResultStructuredOnly(makeResult(dashboard)), nil
}
