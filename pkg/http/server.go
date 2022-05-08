package http

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/HCY2315/chaoyue-golib/pkg/errors"
	"github.com/HCY2315/chaoyue-golib/pkg/log"
	"github.com/HCY2315/chaoyue-golib/pkg/tcp"
)

type HTTPServer struct {
	host            string
	port            int
	server          *http.Server
	shutDownTimeout time.Duration
	listenTimeout   time.Duration
}

func NewHTTPServer(handler http.Handler, host string, port int) *HTTPServer {
	addr := addrImpl(host, port)
	shutDownTO := 100 * time.Millisecond // 随着请求量上涨，一开始设置过高不利于重启
	runListenTO := 1 * time.Second
	s := &HTTPServer{
		host:            host,
		port:            port,
		server:          &http.Server{Addr: addr, Handler: handler},
		shutDownTimeout: shutDownTO,
		listenTimeout:   runListenTO,
	}
	return s
}

func (s *HTTPServer) Run() error {
	addr := s.addr()
	errCh := make(chan error, 1)
	go func() {
		if errRun := s.server.ListenAndServe(); errRun != nil { //REF
			errCh <- errRun
			return
		}
	}()
	if err := s.checkRunErr(errCh); err != nil {
		return err
	}
	log.Infof("服务已启动:%s", addr)
	return nil
}
func (s *HTTPServer) Stop(ctx context.Context) error {
	log.Infof("开始关闭服务")
	shutDownCtx, cancel := context.WithTimeout(ctx, s.shutDownTimeout)
	defer cancel()
	if errShut := s.server.Shutdown(shutDownCtx); errShut != nil {
		return errShut
	}
	return nil
}

func (s *HTTPServer) checkRunErr(errCh chan error) error {
	addr := s.addr()
	select {
	case err := <-errCh:
		return errors.Wrap(err, "启动监听[%s]", addr)
	case <-time.After(s.listenTimeout):
		if opened := tcp.IsTCPAddrOpen(s.host, strconv.Itoa(s.port)); !opened {
			return fmt.Errorf("长时间未成功监听:[%s]", addr)
		}
		return nil
	}
}

func (s *HTTPServer) addr() string {
	return addrImpl(s.host, s.port)
}

func addrImpl(host string, port int) string {
	return fmt.Sprintf("%s:%d", host, port)
}
