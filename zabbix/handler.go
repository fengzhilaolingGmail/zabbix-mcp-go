package zabbix

import (
	"sync"
	"time"

	"zabbixMcp/logger"
)

// ZabbixClientHandler 是对外可见的抽象，隐藏内部具体的 ZabbixClient 或 ClientPool 实现
type ZabbixClientHandler interface {
	SetServerTimezone(tz string)
	GetCachedVersion() *VersionInfo
	SetCachedVersion(v *VersionInfo)
	ClearCachedVersion()
	// GetDetailedVersionFeatures 返回针对当前底层实例（或池）的详细能力映射
	GetDetailedVersionFeatures() map[string]interface{}
	// AdaptAPIParams 根据底层 Zabbix 版本调整传入的 API 参数，返回调整后的参数拷贝
	AdaptAPIParams(method string, params map[string]interface{}) map[string]interface{}
	// Info 返回关于底层实例的元信息（对单实例返回长度为1的切片）
	Info(string) []ClientInfo
	Call(method string, params interface{}) (interface{}, error)
}

// singleClientHandler 将单个 *ZabbixClient 封装为 ZabbixClientHandler。
// 注意：此类型为包内不可导出，实现细节不会暴露给包外使用者。
type singleClientHandler struct {
	client *ZabbixClient
}

// poolClientHandler 将 *ClientPool 封装为 ZabbixClientHandler，向外只暴露统一的调用接口。
type poolClientHandler struct {
	pool *ClientPool
}

// NewSingleClientHandlerFromConfig 使用 ClientConfig 创建一个单一客户端的 handler。
// 返回类型为接口，调用者无法直接访问内部的 *ZabbixClient。
func NewSingleClientHandlerFromConfig(cfg ClientConfig) ZabbixClientHandler {
	cli := NewZabbixClientFromConfig(cfg)
	return &singleClientHandler{client: cli}
}

// NewPoolClientHandlerFromConfigs 使用一组 ClientConfig 创建一个 ClientPool 并返回对应的 handler。
// 每个配置会产生一个 pool 成员。返回的接口隐藏了 *ClientPool 的具体实现。
func NewPoolClientHandlerFromConfigs(cfgs []ClientConfig) (ZabbixClientHandler, error) {
	if len(cfgs) == 0 {
		return nil, nil
	}
	p := NewClientPool(len(cfgs))
	for _, cfg := range cfgs {
		cli := NewZabbixClientFromConfig(cfg)
		if err := p.Add(cli); err != nil {
			return nil, err
		}
	}
	// 预热每个客户端的版本缓存，避免延迟检测导致只有部分实例被探测到的情况。
	var wg sync.WaitGroup
	for _, c := range p.all {
		if c == nil {
			continue
		}
		wg.Add(1)
		go func(cli *ZabbixClient) {
			defer wg.Done()
			if ver, err := NewVersionDetector(cli).DetectVersion(); err == nil {
				logger.L().Info(cli.InstenceName + " API版本信息: " + ver.Full)
			} else {
				logger.L().Error("版本检测失败: " + err.Error())
			}
		}(c)
	}
	wg.Wait()
	return &poolClientHandler{pool: p}, nil
}

// --- singleClientHandler 方法实现 ---
func (h *singleClientHandler) SetServerTimezone(tz string) {
	if h == nil || h.client == nil {
		return
	}
	h.client.SetServerTimezone(tz)
}

func (h *singleClientHandler) GetCachedVersion() *VersionInfo {
	if h == nil || h.client == nil {
		return nil
	}
	return h.client.GetCachedVersion()
}

func (h *singleClientHandler) SetCachedVersion(v *VersionInfo) {
	if h == nil || h.client == nil {
		return
	}
	h.client.SetCachedVersion(v)
}

func (h *singleClientHandler) ClearCachedVersion() {
	if h == nil || h.client == nil {
		return
	}
	h.client.ClearCachedVersion()
}

func (h *singleClientHandler) Call(method string, params interface{}) (interface{}, error) {
	if h == nil || h.client == nil {
		return nil, nil
	}
	return h.client.Call(method, params)
}

func (h *singleClientHandler) Info(instenceName string) []ClientInfo {
	if h == nil || h.client == nil {
		return []ClientInfo{}
	}
	info := ClientInfo{
		InstenceName: h.client.InstenceName,
		User:         h.client.User,
		AuthType:     h.client.AuthType,
		ServerTZ:     h.client.ServerTZ,
		InUse:        false,
		AddedAt:      time.Now(),
	}
	return []ClientInfo{info}
}

// GetDetailedVersionFeatures 实现在 singleClientHandler 上，优先使用当前客户端检测到的版本
func (h *singleClientHandler) GetDetailedVersionFeatures() map[string]interface{} {
	if h == nil || h.client == nil {
		return make(map[string]interface{})
	}
	vd := NewVersionDetector(h.client)
	return vd.GetDetailedVersionFeatures()
}

// AdaptAPIParams 在 singleClientHandler 上实现：根据底层版本适配参数
func (h *singleClientHandler) AdaptAPIParams(method string, params map[string]interface{}) map[string]interface{} {
	if h == nil || h.client == nil {
		return params
	}
	vd := NewVersionDetector(h.client)
	return vd.AdaptAPIParams(method, params)
}

// --- poolClientHandler 方法实现 ---
func (h *poolClientHandler) SetServerTimezone(tz string) {
	if h == nil || h.pool == nil {
		return
	}
	// 对池中所有客户端逐个设置
	for _, c := range h.pool.all {
		if c != nil {
			c.SetServerTimezone(tz)
		}
	}
}

func (h *poolClientHandler) GetCachedVersion() *VersionInfo {
	if h == nil || h.pool == nil {
		return nil
	}
	// 返回第一个有缓存的版本（若无则返回nil）
	for _, c := range h.pool.all {
		if c == nil {
			continue
		}
		if v := c.GetCachedVersion(); v != nil {
			return v
		}
	}
	return nil
}

func (h *poolClientHandler) SetCachedVersion(v *VersionInfo) {
	if h == nil || h.pool == nil {
		return
	}
	for _, c := range h.pool.all {
		if c != nil {
			c.SetCachedVersion(v)
		}
	}
}

func (h *poolClientHandler) ClearCachedVersion() {
	if h == nil || h.pool == nil {
		return
	}
	for _, c := range h.pool.all {
		if c != nil {
			c.ClearCachedVersion()
		}
	}
}

func (h *poolClientHandler) Call(method string, params interface{}) (interface{}, error) {
	if h == nil || h.pool == nil {
		return nil, nil
	}
	client, err := h.pool.Get()
	if err != nil {
		return nil, err
	}
	// 确保在返回前将 client 归还到池
	defer func() {
		_ = h.pool.Release(client)
	}()
	return client.Call(method, params)
}

func (h *poolClientHandler) Info(instenceName string) []ClientInfo {
	if h == nil || h.pool == nil {
		return []ClientInfo{}
	}
	return h.pool.Info(instenceName)
}

// GetDetailedVersionFeatures 在 poolClientHandler 上实现：使用池中任意一个可用客户端进行版本检测，若无可用客户端返回默认功能集
func (h *poolClientHandler) GetDetailedVersionFeatures() map[string]interface{} {
	// 尝试找到一个非nil的客户端用于检测
	for _, c := range h.pool.all {
		if c != nil {
			vd := NewVersionDetector(c)
			return vd.GetDetailedVersionFeatures()
		}
	}
	res := make(map[string]interface{})
	return res
}

// AdaptAPIParams 在 poolClientHandler 上实现：尝试用池中任意客户端的版本信息来适配参数
func (h *poolClientHandler) AdaptAPIParams(method string, params map[string]interface{}) map[string]interface{} {
	if h == nil || h.pool == nil {
		return params
	}
	for _, c := range h.pool.all {
		if c != nil {
			vd := NewVersionDetector(c)
			return vd.AdaptAPIParams(method, params)
		}
	}
	return params
}
