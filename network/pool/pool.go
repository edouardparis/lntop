package pool

import (
	"context"
	"errors"
	"sync"
	"time"

	"google.golang.org/grpc"
)

var (
	// ErrClosed is the error when the client pool is closed
	ErrClosed = errors.New("grpc pool: client pool is closed")
	// ErrTimeout is the error when the client pool timed out
	ErrTimeout = errors.New("grpc pool: client pool timed out")
	// ErrAlreadyClosed is the error when the client conn was already closed
	ErrAlreadyClosed = errors.New("grpc pool: the connection was already closed")
	// ErrFullPool is the error when the pool is already full
	ErrFullPool = errors.New("grpc pool: closing a ClientConn into a full pool")
)

// Factory is a function type creating a grpc client
type Factory func() (*grpc.ClientConn, error)

type Pool struct {
	conns   chan Conn
	factory Factory
	mu      sync.RWMutex
	timeout time.Duration
}

func (p *Pool) getConns() chan Conn {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.conns
}

func (p *Pool) Get(ctx context.Context) (*Conn, error) {
	conns := p.getConns()
	if conns == nil {
		return nil, ErrClosed
	}

	conn := Conn{
		pool: p,
	}

	select {
	case conn = <-conns:
	case <-ctx.Done():
		return nil, ErrTimeout
	}

	if conn.ClientConn != nil &&
		p.timeout > 0 &&
		conn.usedAt.Add(p.timeout).Before(time.Now()) {
		conn.ClientConn.Close()
		conn.ClientConn = nil
	}

	var err error
	if conn.ClientConn == nil {
		conn.ClientConn, err = p.factory()
		if err != nil {
			conns <- Conn{
				pool: p,
			}
		}
	}

	return &conn, err
}

func (p *Pool) Close() {
	p.mu.Lock()
	conns := p.conns
	p.conns = nil
	p.mu.Unlock()

	if conns == nil {
		return
	}

	close(conns)
	for i := 0; i < p.Capacity(); i++ {
		client := <-conns
		if client.ClientConn == nil {
			continue
		}
		client.ClientConn.Close()
	}
}

func (p *Pool) Capacity() int {
	if p.IsClosed() {
		return 0
	}
	return cap(p.conns)
}

func (p *Pool) IsClosed() bool {
	return p == nil || p.getConns() == nil
}

func New(factory Factory, capacity int, timeout time.Duration) (*Pool, error) {
	p := &Pool{
		conns:   make(chan Conn, capacity),
		factory: factory,
		timeout: timeout,
	}

	// Fill the pool with empty clients
	for i := 0; i < capacity; i++ {
		p.conns <- Conn{
			pool: p,
		}
	}
	return p, nil
}
