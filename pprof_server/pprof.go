/*
*
User: cr-mao
Date: 2024/2/20 11:55
Email: crmao@qq.com
Desc: pprof.go
*/
package pprof_server

import (
	"context"
	"fmt"

	"net/http"
	_ "net/http/pprof"

	"github.com/cr-mao/lori/log"
	"github.com/cr-mao/lori/transport"
)

var _ transport.Server = &pprof{}

type pprof struct {
	addr string
}

func NewPProf(addr string) *pprof {
	return &pprof{
		addr: addr,
	}
}

func (p *pprof) Name() string {
	return "pprof"
}

func (p *pprof) Start(_ context.Context) error {
	log.Debug("pprof addr:", p.addr)
	fmt.Println("pprof addr:", p.addr)
	err := http.ListenAndServe(p.addr, nil)
	if err != nil {
		log.Errorf("pprof server start failed: %v", err)
	}
	return err
}

func (p *pprof) Stop(ctx context.Context) error {
	return nil
}
