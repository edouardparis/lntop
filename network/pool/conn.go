package pool

import (
	"time"

	"google.golang.org/grpc"
)

// Conn is the wrapper for a grpc client conn.
type Conn struct {
	*grpc.ClientConn
	pool   *Pool
	usedAt time.Time
}

func (c *Conn) Close() error {
	if c == nil {
		return nil
	}
	if c.ClientConn == nil {
		return ErrAlreadyClosed
	}
	if c.pool.IsClosed() {
		return ErrClosed
	}

	conn := Conn{
		pool:       c.pool,
		ClientConn: c.ClientConn,
	}
	select {
	case c.pool.conns <- conn:
	default:
		return ErrFullPool
	}
	return nil
}
