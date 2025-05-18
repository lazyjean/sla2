package test

import (
	"net"

	"github.com/lazyjean/sla2/internal/infrastructure/listen"
	"google.golang.org/grpc/test/bufconn"
)

type Listener struct {
	grpcListener *bufconn.Listener
	httpListener *bufconn.Listener
}

// NewTestListener 创建新的测试监听器实例
func NewTestListener() listen.Listener {
	return &Listener{
		grpcListener: bufconn.Listen(1024 * 1024),
		httpListener: bufconn.Listen(1024 * 1024),
	}
}

func (l *Listener) ListenGRPC(port int) error {
	return nil
}

func (l *Listener) ListenHTTP(port int) error {
	return nil
}

func (l *Listener) GetGRPCListener() net.Listener {
	return l.grpcListener
}

func (l *Listener) GetHTTPListener() net.Listener {
	return l.httpListener
}

func (l *Listener) Stop() error {
	if l.httpListener != nil {
		l.httpListener.Close()
	}
	if l.grpcListener != nil {
		l.grpcListener.Close()
	}
	return nil
}

var _ listen.Listener = (*Listener)(nil)
