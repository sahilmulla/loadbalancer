package pool

import (
	"errors"
	"sync"

	"github.com/sahilmulla/loadbalancer/pkg/service"
)

var (
	ErrAllDown = errors.New("all services down")
)

type Pool interface {
	AddService(service.Service)
	GetNextService() (service.Service, error)
}

type roundRobinPool struct {
	services []service.Service
	current  uint64

	mux sync.RWMutex
}

func (p *roundRobinPool) AddService(s service.Service) {
	p.services = append(p.services, s)
}

func (p *roundRobinPool) GetNextService() (service.Service, error) {
	currIdx := p.current

	for {
		ns := p.services[p.current]
		if ns.IsAlive() {
			p.nextIndex()
			return ns, nil
		}

		p.nextIndex()
		if p.current == currIdx {
			break
		}
	}

	return nil, ErrAllDown
}

func (p *roundRobinPool) nextIndex() {
	p.mux.Lock()
	defer p.mux.Unlock()

	p.current = (p.current + 1) % uint64(len(p.services))
}

func NewRoundRobinPool() *roundRobinPool {
	return &roundRobinPool{
		services: make([]service.Service, 0),
		current:  0,
		mux:      sync.RWMutex{},
	}
}
