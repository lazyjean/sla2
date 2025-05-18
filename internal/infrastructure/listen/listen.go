package listen

import (
	"fmt"
	"net"
)

type Listener interface {
	ListenGRPC(port int) error
	ListenHTTP(port int) error
	GetGRPCListener() net.Listener
	GetHTTPListener() net.Listener
	Stop() error
}

type listener struct {
	httpListener net.Listener
	grpcListener net.Listener
}

func NewListener() Listener {
	return &listener{}
}

func (l *listener) ListenGRPC(port int) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	l.grpcListener = listener
	return nil
}

func (l *listener) ListenHTTP(port int) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	l.httpListener = listener
	return nil
}

func (l *listener) GetGRPCListener() net.Listener {
	return l.grpcListener
}

func (l *listener) GetHTTPListener() net.Listener {
	return l.httpListener
}

func (l *listener) Stop() error {
	if l.httpListener != nil {
		return l.httpListener.Close()
	}
	if l.grpcListener != nil {
		return l.grpcListener.Close()
	}
	return nil
}
