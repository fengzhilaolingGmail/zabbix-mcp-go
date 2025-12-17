package zabbix

import (
	"context"
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
	// closed indicates whether the pool has been closed
	closed    bool
	closeOnce sync.Once
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
	if client == nil {
		return errors.New("nil client")
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return errors.New("pool is closed")
	}

	if len(p.all) >= p.capacity {
		return ErrPoolFull
	}

	// 防止重复添加同一个客户端
	for _, existing := range p.all {
		if existing == client {
			return errors.New("client already added to pool")
		}
	}

	// 记录元信息并放入可用通道（此处可能会阻塞直到通道有空间）
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

// GetWithContext 尝试在给定的 context 内获取可用客户端，context 取消或超时时返回错误
func (p *ClientPool) GetWithContext(ctx context.Context) (*ZabbixClient, error) {
	select {
	case client, ok := <-p.ch:
		if !ok {
			return nil, ErrPoolEmpty
		}
		p.mu.Lock()
		p.inUse[client] = true
		p.mu.Unlock()
		return client, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// TryGet 在 timeout 时限内尝试获取客户端，超时返回 ErrPoolEmpty
func (p *ClientPool) TryGet(timeout time.Duration) (*ZabbixClient, error) {
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case client, ok := <-p.ch:
		if !ok {
			return nil, ErrPoolEmpty
		}
		p.mu.Lock()
		p.inUse[client] = true
		p.mu.Unlock()
		return client, nil
	case <-timer.C:
		return nil, ErrPoolEmpty
	}
}

// Release 将 client 归还到池中；如果池已满返回 ErrPoolFull
func (p *ClientPool) Release(client *ZabbixClient) error {
	if client == nil {
		return errors.New("nil client")
	}

	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return errors.New("pool is closed")
	}
	_, known := p.addedAt[client]
	if !known {
		p.mu.Unlock()
		return errors.New("client does not belong to pool")
	}
	// 检查是否已经被归还
	if inUse, ok := p.inUse[client]; ok && !inUse {
		p.mu.Unlock()
		return errors.New("client already released")
	}
	p.mu.Unlock()

	// 将 client 放回通道；我们选择阻塞直到成功归还（避免丢失客户端）
	select {
	case p.ch <- client:
		p.mu.Lock()
		p.inUse[client] = false
		p.mu.Unlock()
		return nil
	default:
		// 如果默认分支发生，意味着通道暂时已满（非常罕见，但我们仍然阻塞以确保归还）
		p.ch <- client
		p.mu.Lock()
		p.inUse[client] = false
		p.mu.Unlock()
		return nil
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

// Close 关闭连接池并释放资源，关闭后不能再 Add 或 Release
func (p *ClientPool) Close() {
	p.closeOnce.Do(func() {
		p.mu.Lock()
		defer p.mu.Unlock()
		p.closed = true
		close(p.ch)
		// 可在此处对所有客户端执行额外清理（例如释放连接/句柄）
	})
}
