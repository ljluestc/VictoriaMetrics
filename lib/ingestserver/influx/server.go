package influx

import (
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/logger"
	"net"
	"sync"
	"time"
)

type Server struct {
	ln      *net.TCPListener
	connWG  sync.WaitGroup
	stopCh  chan struct{}
	timeout time.Duration
}

func Start(addr string, connTimeout time.Duration) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		logger.Errorf("cannot resolve TCP addr %s: %v", addr, err)
		return err
	}
	ln, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		logger.Errorf("cannot listen on %s: %v", addr, err)
		return err
	}
	_ = ln
	// ...existing code...
	return nil
}

func (s *Server) Stop() {
	if s == nil {
		return // Add nil check to prevent panic
	}
	if s.ln != nil {
		_ = s.ln.Close()
	}
	close(s.stopCh)
	s.connWG.Wait()
}
