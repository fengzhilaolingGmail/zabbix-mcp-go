package zabbix

import (
	"errors"
	"sync"
	"time"
)

// ClientInfo 描述连接池中客户端的详细信息
type ClientInfo struct {
	URL      string    `json:"url"`
	User     string    `json:"user"`
	AuthType string    `json:"auth_type"`
	ServerTZ string    `json:"server_tz"`
	InUse    bool      `json:"in_use"`
	AddedAt  time.Time `json:"added_at"`
}

// ClientPool 管理一组可复用的 ZabbixClient
type ClientPool struct {
	ch       chan *ZabbixClient
	mu       sync.Mutex
	all      []*ZabbixClient
	addedAt  map[*ZabbixClient]time.Time
	inUse    map[*ZabbixClient]bool
	capacity int
}

// ErrPoolFull 当池已满时返回
var ErrPoolFull = errors.New("client pool is full")

// ErrPoolEmpty 当池为空且无法获取时返回
var ErrPoolEmpty = errors.New("client pool is empty")

// NewClientPool 创建一个容量为 capacity 的连接池（capacity 必须 >=1）
func NewClientPool(capacity int) *ClientPool {
	if capacity <= 0 {
		capacity = 1
	}
	return &ClientPool{
		ch:       make(chan *ZabbixClient, capacity),
		all:      make([]*ZabbixClient, 0, capacity),
		addedAt:  make(map[*ZabbixClient]time.Time),
		inUse:    make(map[*ZabbixClient]bool),
		capacity: capacity,
	}
}

// NewClientPoolWithFactory 使用工厂函数初始化并填充池
func NewClientPoolWithFactory(factory func() *ZabbixClient, capacity int) (*ClientPool, error) {
	p := NewClientPool(capacity)
	for i := 0; i < capacity; i++ {
		c := factory()
		if c == nil {
			return nil, errors.New("factory returned nil client")
		}
		if err := p.Add(c); err != nil {
			return nil, err
		}
	}
	return p, nil
}

// Add 将 client 添加到池中，如果已满返回 ErrPoolFull
func (p *ClientPool) Add(client *ZabbixClient) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.all) >= p.capacity {
		return ErrPoolFull
	}

	// 记录元信息并放入可用通道
	p.all = append(p.all, client)
	p.addedAt[client] = time.Now()
	p.inUse[client] = false
	p.ch <- client
	return nil
}

// Get 从池中获取一个客户端（阻塞直到有可用的客户端或池被破坏）
func (p *ClientPool) Get() (*ZabbixClient, error) {
	client, ok := <-p.ch
	if !ok {
		return nil, ErrPoolEmpty
	}
	p.mu.Lock()
	p.inUse[client] = true
	p.mu.Unlock()
	return client, nil
}

// TryGet 在 timeout 时限内尝试获取客户端，超时返回 ErrPoolEmpty
func (p *ClientPool) TryGet(timeout time.Duration) (*ZabbixClient, error) {
	select {
	case client, ok := <-p.ch:
		if !ok {
			return nil, ErrPoolEmpty
		}
		p.mu.Lock()
		p.inUse[client] = true
		p.mu.Unlock()
		return client, nil
	case <-time.After(timeout):
		return nil, ErrPoolEmpty
	}
}

// Release 将 client 归还到池中；如果池已满返回 ErrPoolFull
func (p *ClientPool) Release(client *ZabbixClient) error {
	p.mu.Lock()
	_, known := p.addedAt[client]
	p.mu.Unlock()
	if !known {
		return errors.New("client does not belong to pool")
	}

	select {
	case p.ch <- client:
		p.mu.Lock()
		p.inUse[client] = false
		p.mu.Unlock()
		return nil
	default:
		return ErrPoolFull
	}
}

// Total 返回池中注册的客户端总数（capacity 内实际添加的）
func (p *ClientPool) Total() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.all)
}

// Available 返回当前可立即获取的客户端数（非使用中）
func (p *ClientPool) Available() int {
	return len(p.ch)
}

// Capacity 返回池容量
func (p *ClientPool) Capacity() int {
	return p.capacity
}

// Info 返回每个实例的详细信息
func (p *ClientPool) Info() []ClientInfo {
	p.mu.Lock()
	defer p.mu.Unlock()
	out := make([]ClientInfo, 0, len(p.all))
	for _, c := range p.all {
		inuse := false
		if v, ok := p.inUse[c]; ok {
			inuse = v
		}
		added := p.addedAt[c]
		info := ClientInfo{
			URL:      c.URL,
			User:     c.User,
			AuthType: c.AuthType,
			ServerTZ: c.ServerTZ,
			InUse:    inuse,
			AddedAt:  added,
		}
		out = append(out, info)
	}
	return out
}

// HealthCheck 对池中所有实例并发执行简单检查（调用 apiinfo.version），返回每个实例的健康状态
// 注意：该方法会发起网络请求，调用方需考虑频率与超时
func (p *ClientPool) HealthCheck(timeout time.Duration) map[string]bool {
	p.mu.Lock()
	clients := make([]*ZabbixClient, len(p.all))
	copy(clients, p.all)
	p.mu.Unlock()

	results := make(map[string]bool)
	var wg sync.WaitGroup
	mu := sync.Mutex{}

	for _, c := range clients {
		wg.Add(1)
		go func(cli *ZabbixClient) {
			defer wg.Done()
			// Try to call apiinfo.version with a timeout goroutine
			ok := false
			done := make(chan struct{})
			go func() {
				_, err := cli.Call("apiinfo.version", []interface{}{})
				if err == nil {
					ok = true
				}
				close(done)
			}()
			select {
			case <-done:
				// finished
			case <-time.After(timeout):
				ok = false
			}
			mu.Lock()
			results[cli.URL] = ok
			mu.Unlock()
		}(c)
	}
	wg.Wait()
	return results
}
