package models

// ParamSpec 用于描述一个方法的业务参数，由具体实现转换成 Zabbix API 所需的 map
// 每个具体的 API spec 应当实现 BuildParams 以便统一适配
// MapParams 为旧的 map[string]interface{} 调用方式提供兼容

type ParamSpec interface {
	BuildParams() map[string]interface{}
}

// MapParams 允许沿用 map[string]interface{} 的方式，同时实现 ParamSpec 接口
type MapParams map[string]interface{}

// BuildParams 返回 map 的浅拷贝，避免调用方修改底层存储
func (m MapParams) BuildParams() map[string]interface{} {
	cloned := make(map[string]interface{}, len(m))
	for k, v := range m {
		cloned[k] = v
	}
	return cloned
}
