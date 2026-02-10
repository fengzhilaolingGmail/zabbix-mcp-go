/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2026-01-02 16:33:29
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2026-02-10 17:38:11
 * @FilePath: \zabbix-mcp-go\register\item.go
 * @Description: mcp注册监控项相关工具
 * Copyright (c) 2026 by fengzhilaoling@gmail.com, All Rights Reserved.
 */

package register

import (
	"zabbixMcp/handler"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerItem(s *server.MCPServer) {
	s.AddTool(
		mcp.NewTool("get_items",
			mcp.WithDescription("获取实例Zabbix监控项信息,支持所有或模糊匹配. 注意: 必须提供主机过滤(host_ids/hostname)或监控项过滤(item_key/item_name)至少一种"),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Zabbix实例名称必须填")),
			mcp.WithArray("host_ids", mcp.Description("主机ID列表 (与hostname二选一,至少提供一个)")),
			mcp.WithString("hostname", mcp.Description("主机名称 (与host_ids二选一,至少提供一个)")),
			mcp.WithString("item_key", mcp.Description("监控项键(key) (与item_name二选一,至少提供一个)")),
			mcp.WithString("item_name", mcp.Description("监控项名称(name),需要使用通配符(与item_key二选一,至少提供一个)")),
		),
		handler.GetItemsHandler,
	)
	s.AddTool(
		mcp.NewTool("get_host_item_from_item_key",
			mcp.WithDescription("根据主机id和监控项 KEY 查主机监控项"),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Zabbix实例名称必须填")),
			mcp.WithArray("host_ids", mcp.Required(), mcp.Description("主机ID列表")),
			mcp.WithString("item_key", mcp.Description("监控项键(key)")),
		),
		handler.GetItemsHandlerNew,
	)
	s.AddTool(
		mcp.NewTool("get_host_item_from_item_name",
			mcp.WithDescription("根据主机id和监控项 NAME 查主机监控项"),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Zabbix实例名称必须填")),
			mcp.WithArray("host_ids", mcp.Required(), mcp.Description("主机ID列表")),
			mcp.WithString("item_name", mcp.Description("监控项名称(name),需要使用通配符")),
		),
		handler.GetItemsHandlerNew,
	)
}

// 基础场景（主机 / 监控项 / 指标）
// get_metric_from_itemid - 通过监控项 ID 获取指标数据
// get_host_from_groupid - 通过主机组 ID 获取主机列表
// get_trigger_from_item - 通过监控项获取关联告警触发器
// get_history_from_metric - 通过指标名获取历史数据
// get_graph_from_hostitem - 通过主机 + 监控项获取趋势图
// 操作类场景（创建 / 更新 / 删除）
// create_item_for_hostname - 为指定主机名创建监控项
// update_trigger_from_itemid - 通过监控项 ID 更新告警触发器
// delete_host_from_groupid - 从主机组 ID 中删除指定主机
// enable_item_for_hostid - 启用指定主机 ID 的监控项
// disable_alert_from_metric - 禁用指定指标的告警规则
// 统计 / 分析类场景
// count_items_from_hostgroup - 统计主机组下的监控项数量
// sum_metrics_from_hostlist - 汇总主机列表的指标总和
// avg_values_from_itemhistory - 计算监控项历史数据的平均值
// check_health_from_hostmetrics - 通过主机指标检查健康状态
// filter_items_from_hosttags - 通过主机标签筛选监控项
// 扩展场景（关联 / 转换）
// map_hostname_to_itemkey - 映射主机名到监控项 key
// convert_item_to_metricformat - 将监控项数据转换为指标格式
// link_alert_to_hostitem - 关联告警到主机监控项
// export_items_from_hostfilter - 导出符合主机筛选条件的监控项
// import_metrics_to_hostitem - 导入指标数据到主机监控项
