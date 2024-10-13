package pool

import (
	"errors"
	"fmt"
	"sync"

	"github.com/sahilmulla/loadbalancer/pkg/service"
)

var (
	ErrAllDown   = errors.New("all services down")
	ErrDuplicate = errors.New("already added")
)

type Pool interface {
	AddService(service.Service) error
	GetNextService() (service.Service, error)
}

type roundRobinPool struct {
	services []service.Service
	current  uint64

	mux sync.RWMutex
}

func (p *roundRobinPool) AddService(s service.Service) error {
	p.mux.Lock()
	defer p.mux.Unlock()
	for _, service := range p.services {
		fmt.Println("s", service.GetURL().Host)
		if service.GetURL().Host == s.GetURL().Host {
			if !service.IsAlive() {
				service.SetAlive(true)
				return nil
			}

			return ErrDuplicate
		}
	}

	p.services = append(p.services, s)

	return nil
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
