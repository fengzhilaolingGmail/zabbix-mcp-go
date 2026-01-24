package models

type Groups struct {
	GroupID string `json:"groupid,comment:主机组ID"`
}

// ZabbixTag 主机标签参数
type Tag struct {
	Tag   string `json:"tag"`   // 标签字符串
	Value string `json:"value"` // 标签值字符串
}

// ZabbixMacro 主机宏参数
type Macro struct {
	GlobalMacroID string `json:"globalmacroid"` // 全局宏的ID
	Macro         string `json:"macro"`         // 宏字符串
	Value         string `json:"value"`         // 宏的值
	Type          int    `json:"type"`          // 宏的类型 0 - (默认) 文本宏, 1 - 密文宏, 2 - 密钥宏
	Description   string `json:"description"`   // 宏描述信息
}

type Macros struct {
	Macro       string `json:"macro,comment:宏字符串"`                  // 宏字符串
	Value       string `json:"value,comment:宏的值"`                   // 宏的值
	Description string `json:"description,omitempty,comment:宏描述信息"` // 宏描述信息
}

// ZabbixInterface 主机接口参数
type ZabbixInterface struct {
	InterfaceID  string   `json:"interfaceid,comment:接口ID"`                                          // 接口ID.
	Available    int      `json:"available,omitempty,comment:主机接口的可用性:0(默认)-未知;1-可用;2-不可用"`          // 主机 接口的可用性. 0 - (默认) 未知; 1 - 可用; 2 - 不可用.
	HostID       string   `json:"hostid,comment:接口所属的主机ID"`                                          // 接口所属的 主机 ID.
	Type         int      `json:"type,omitempty,comment:接口类型:1-Agent;2-SNMP;3-IPMI;4-JMX"`           // 接口类型. 1 - Agent; 2 - SNMP; 3 - IPMI; 4 - JMX.
	IP           string   `json:"ip,omitempty,comment:接口使用的IP地址.如果通过DNS连接可以为空."`                     // 接口使用的IP地址. 如果通过DNS连接可以为空.
	DNS          string   `json:"dns,omitempty,comment:接口使用的DNS名称.如果通过IP连接可以为空."`                    // 接口使用的DNS名称. 如果通过IP连接可以为空.
	Port         string   `json:"port,omitempty,comment:接口使用的端口号.可包含用户宏."`                           // 接口使用的端口号. 可包含用户宏.
	UseIP        int      `json:"useip,omitempty,comment:是否应通过IP连接:0-使用主机DNS名称连接;1-使用主机IP地址连接.(默认)"` // 是否应通过IP连接. 0 - 使用 主机DNS名称 连接; 1 - 使用 主机IP地址 连接.(默认)
	Main         int      `json:"main,omitempty,comment:接口是否作为主机的默认接口:0-非默认;1-默认.(默认)"`              // 接口是否作为 主机 的默认接口. 每种类型只能有一个接口在 一个主机 上设置为默认. 0 - 非默认; 1 - 默认. (默认)
	Details      []string `json:"details,omitempty,comment:接口的额外详情object."`                          // 接口的额外详情 object.
	DisableUntil int      `json:"disable_until,omitempty,comment:接口不可用的下次轮询时间.0-默认(立即可用)"`           // 不可用 主机 接口的下次轮询时间.
	Error        string   `json:"error,omitempty,comment:主机接口不可用时的错误文本."`                            // 当 主机 接口不可用时的错误文本.
	ErrorsFrom   int      `json:"errors_from,omitempty,comment:主机接口变为不可用的时间."`                       // 主机 接口变为不可用的时间.
}

// ZabbixTemplate 主机模板参数
type Template struct {
	TemplateID    string `json:"templateid,comment:模版的ID"`                                             // 模版的ID。只读。必需 对于更新操作
	Host          string `json:"host,comment:模版的技术名称"`                                                 // 必需 对于创建操作
	Description   string `json:"description,omitempty,comment:模版的描述"`                                  // 模版的描述
	Name          string `json:"name,omitempty,comment:模板的可见名称;默认:host属性值"`                            // 模板的可见名称，默认使用 Host
	UUID          string `json:"uuid,omitempty,comment:通用唯一标识符(UUID),用于将导入的模板链接到已存在的模板.如果未提供,则会自动生成。"` // UUID 用于导入关联
	VendorName    string `json:"vendor_name,omitempty,comment:模板的供应商名称.创建时vendor_name和vendor_version应同时设置或同时留空。"`
	VendorVersion string `json:"vendor_version,omitempty,comment:模板的供应商版本.创建时vendor_name和vendor_version应同时设置或同时留空。"`
}

type Templates struct {
	TemplateID string `json:"templateid,comment:模版的ID"` // 模版的ID。只读。必需 对于更新操作
}
