package models

// HostGetParams 提供 host.get 常用参数的结构化封装
type HostGetParams struct {
	HostIDs  []string
	GroupIDs []string
	Output   string
}

// BuildParams 将 HostGetParams 转换为 API 参数
func (p HostGetParams) BuildParams() map[string]interface{} {
	params := map[string]interface{}{}
	if len(p.HostIDs) > 0 {
		params["hostids"] = p.HostIDs
	}
	if len(p.GroupIDs) > 0 {
		params["groupids"] = p.GroupIDs
	}
	if p.Output != "" {
		params["output"] = p.Output
	}
	return params
}

func (p HostGetParams) BuildDeleteParams() []string {
	if len(p.HostIDs) > 0 {
		return append([]string(nil), p.HostIDs...)
	}
	return nil
}
