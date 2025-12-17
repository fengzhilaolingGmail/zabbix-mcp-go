package zabbix_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	zabbix "zabbixMcp/zabbix"
)

// fakeRoundTripper 返回固定响应体以模拟 apiinfo.version
type fakeRoundTripper struct {
	respBody string
}

func (f *fakeRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(f.respBody)),
		Header:     make(http.Header),
	}, nil
}

func TestPoolAddGetRelease(t *testing.T) {
	p := zabbix.NewClientPool(2)
	c1 := zabbix.NewZabbixClient("http://a", "u1", "p", 1)
	c2 := zabbix.NewZabbixClient("http://b", "u2", "p", 1)

	if err := p.Add(c1); err != nil {
		t.Fatalf("Add c1 failed: %v", err)
	}
	if err := p.Add(c2); err != nil {
		t.Fatalf("Add c2 failed: %v", err)
	}

	if p.Total() != 2 {
		t.Fatalf("expected total 2, got %d", p.Total())
	}
	if p.Available() != 2 {
		t.Fatalf("expected available 2, got %d", p.Available())
	}

	got, err := p.Get()
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if got == nil {
		t.Fatalf("Get returned nil")
	}

	if p.Available() != 1 {
		t.Fatalf("expected available 1 after get, got %d", p.Available())
	}

	if err := p.Release(got); err != nil {
		t.Fatalf("Release failed: %v", err)
	}
	if p.Available() != 2 {
		t.Fatalf("expected available 2 after release, got %d", p.Available())
	}
}

func TestTryGetTimeout(t *testing.T) {
	p := zabbix.NewClientPool(1)
	c := zabbix.NewZabbixClient("http://x", "u", "p", 1)
	if err := p.Add(c); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	got, err := p.Get()
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if got == nil {
		t.Fatalf("expected client")
	}

	// now pool empty; TryGet with short timeout should fail
	_, err = p.TryGet(10 * time.Millisecond)
	if err == nil {
		t.Fatalf("expected timeout error, got nil")
	}

	// release and try again
	_ = p.Release(got)
	g2, err := p.TryGet(10 * time.Millisecond)
	if err != nil {
		t.Fatalf("TryGet after release failed: %v", err)
	}
	_ = p.Release(g2)
}

func TestNewClientPoolWithFactoryAndHealthCheck(t *testing.T) {
	// factory creates clients with fake HTTPClient that returns api version
	idx := 0
	factory := func() *zabbix.ZabbixClient {
		idx++
		cli := zabbix.NewZabbixClient("http://fake", "u", "p", 1)
		cli.URL = fmt.Sprintf("http://fake-%d", idx) // unique URL per instance
		cli.HTTPClient = &http.Client{Transport: &fakeRoundTripper{respBody: `{"jsonrpc":"2.0","result":"6.0.5","id":1}`}}
		return cli
	}

	p, err := zabbix.NewClientPoolWithFactory(factory, 3)
	if err != nil {
		t.Fatalf("NewClientPoolWithFactory failed: %v", err)
	}
	if p.Total() != 3 {
		t.Fatalf("expected total 3, got %d", p.Total())
	}

	// HealthCheck should return true for each
	res := p.HealthCheck(500 * time.Millisecond)
	if len(res) != 3 {
		t.Fatalf("expected 3 results, got %d", len(res))
	}
	for _, ok := range res {
		if !ok {
			t.Fatalf("expected healthy client")
		}
	}
}
