package zabbix

import (
	"context"
	"errors"
	"sync"
	"time"
)

// ClientInfo 描述连接池中客户端的详细信息
type ClientInfo struct {
	Instance  string    `json:"instance"`
	URL       string    `json:"url"`
	User      string    `json:"user"`
	AuthType  string    `json:"auth_type"`
	ServerTZ  string    `json:"server_tz"`
	InUse     bool      `json:"in_use"`
	Connected bool      `json:"connected"`
	AddedAt   time.Time `json:"added_at"`
	Version   string    `json:"version"`
}

var (
	// ErrPoolFull 当池已满时返回
	ErrPoolFull = errors.New("client pool is full")
	// ErrPoolEmpty 当池为空且无法获取时返回
	ErrPoolEmpty = errors.New("client pool is empty")
	// ErrPoolClosed 指示池已经关闭
	ErrPoolClosed = errors.New("client pool is closed")
)

type clientMeta struct {
	addedAt   time.Time
	inUse     bool
	lastError error
}

// ClientPool 管理一组可复用的 ZabbixClient
type ClientPool struct {
	idle      chan *ZabbixClient
	mu        sync.Mutex
	order     []*ZabbixClient
	meta      map[*ZabbixClient]*clientMeta
	capacity  int
	closed    bool
	closeOnce sync.Once
}

// NewClientPool 创建一个容量为 capacity 的连接池（capacity 必须 >=1）
func NewClientPool(capacity int) *ClientPool {
	if capacity <= 0 {
		capacity = 1
	}
	return &ClientPool{
		idle:     make(chan *ZabbixClient, capacity),
		order:    make([]*ZabbixClient, 0, capacity),
		meta:     make(map[*ZabbixClient]*clientMeta),
		capacity: capacity,
	}
}

// NewClientPoolWithFactory 使用工厂函数初始化并填充池
func NewClientPoolWithFactory(factory func() (*ZabbixClient, error), capacity int) (*ClientPool, error) {
	p := NewClientPool(capacity)
	for i := 0; i < capacity; i++ {
		c, err := factory()
		if err != nil {
			return nil, err
		}
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
		return ErrPoolClosed
	}
	if len(p.order) >= p.capacity {
		return ErrPoolFull
	}
	for _, existing := range p.order {
		if existing == client {
			return errors.New("client already added to pool")
		}
	}
	p.order = append(p.order, client)
	p.meta[client] = &clientMeta{addedAt: time.Now()}
	p.idle <- client
	return nil
}

// Acquire 获取一个租借句柄，实现 ClientProvider 接口
func (p *ClientPool) Acquire(ctx context.Context) (ClientLease, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case client, ok := <-p.idle:
		if !ok {
			return nil, ErrPoolClosed
		}
		p.markInUse(client, true)
		return newPoolLease(p, client), nil
	}
}

func (p *ClientPool) markInUse(client *ZabbixClient, inUse bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if meta, ok := p.meta[client]; ok {
		meta.inUse = inUse
	}
}

func (p *ClientPool) releaseClient(client *ZabbixClient, lastErr error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return
	}
	if meta, ok := p.meta[client]; ok {
		meta.inUse = false
		meta.lastError = lastErr
	}
	select {
	case p.idle <- client:
	default:
		// Should not happen; fallback to blocking send to avoid leak
		p.idle <- client
	}
}

// Info 返回每个实例的详细信息
func (p *ClientPool) Info(Instance string) []ClientInfo {
	p.mu.Lock()
	defer p.mu.Unlock()
	out := make([]ClientInfo, 0, len(p.order))
	for _, c := range p.order {
		if Instance != "" && c.Instance != Instance {
			continue
		}
		meta := p.meta[c]
		version := ""
		if v := c.GetCachedVersion(); v != nil {
			version = v.Full
		}
		info := ClientInfo{
			Instance:  c.Instance,
			URL:       c.URL,
			User:      c.User,
			AuthType:  c.AuthType,
			ServerTZ:  c.ServerTZ,
			InUse:     meta != nil && meta.inUse,
			Connected: c.IsConnected(),
			AddedAt:   meta.addedAt,
			Version:   version,
		}
		out = append(out, info)
	}
	return out
}

// Total 返回池中注册的客户端总数（capacity 内实际添加的）
func (p *ClientPool) Total() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.order)
}

// Capacity 返回池容量
func (p *ClientPool) Capacity() int {
	return p.capacity
}

// HealthCheck 对池中所有空闲实例执行简单检查
func (p *ClientPool) HealthCheck(ctx context.Context, timeout time.Duration) map[string]bool {
	if ctx == nil {
		ctx = context.Background()
	}
	results := make(map[string]bool)
	var wg sync.WaitGroup
	var resMu sync.Mutex

	for {
		lease, ok := p.tryAcquire()
		if !ok {
			break
		}
		wg.Add(1)
		go func(l *poolLease) {
			defer wg.Done()
			checkCtx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()
			err := l.Client().Call(checkCtx, "apiinfo.version", []interface{}{}, nil)
			resMu.Lock()
			results[l.client.URL] = err == nil
			resMu.Unlock()
			l.Release(err)
		}(lease)
	}
	wg.Wait()
	return results
}

func (p *ClientPool) tryAcquire() (*poolLease, bool) {
	select {
	case client, ok := <-p.idle:
		if !ok {
			return nil, false
		}
		p.markInUse(client, true)
		return newPoolLease(p, client), true
	default:
		return nil, false
	}
}

// Close 关闭连接池并释放资源，关闭后不能再 Add 或 Acquire
func (p *ClientPool) Close() {
	p.closeOnce.Do(func() {
		p.mu.Lock()
		p.closed = true
		close(p.idle)
		p.mu.Unlock()
	})
}

// 确保 ClientPool 实现 ClientProvider
var _ ClientProvider = (*ClientPool)(nil)

type poolLease struct {
	pool   *ClientPool
	client *ZabbixClient
	once   sync.Once
}

func newPoolLease(pool *ClientPool, client *ZabbixClient) *poolLease {
	return &poolLease{pool: pool, client: client}
}

func (l *poolLease) Client() APIClient {
	return l.client
}

func (l *poolLease) Release(err error) {
	l.once.Do(func() {
		l.pool.releaseClient(l.client, err)
	})
}
