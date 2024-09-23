package pool

import (
	"sync"

	"github.com/sahilmulla/loadbalancer/pkg/service"
)

type Pool interface {
	AddService(service.Service)
	GetNextService() service.Service
}

type roundRobinPool struct {
	services []service.Service
	current  uint64

	mux sync.RWMutex
}

func (p *roundRobinPool) AddService(s service.Service) {
	p.services = append(p.services, s)
}

func (p *roundRobinPool) GetNextService() service.Service {
	return p.services[p.nextIndex()]
}

func (p *roundRobinPool) nextIndex() int {
	p.mux.Lock()
	defer p.mux.Unlock()

	p.current = (p.current + 1) % uint64(len(p.services))
	return int(p.current)
}

func NewRoundRobinPool() *roundRobinPool {
	return &roundRobinPool{
		services: make([]service.Service, 0),
		current:  0,
		mux:      sync.RWMutex{},
	}
}
