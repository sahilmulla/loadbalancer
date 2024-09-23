package balancer

import (
	"errors"
	"net/http"

	"github.com/sahilmulla/loadbalancer/pkg/pool"
	"github.com/sahilmulla/loadbalancer/pkg/service"
)

type Balancer interface {
	Serve(http.ResponseWriter, *http.Request)
}

type balancer struct {
	pool.Pool
}

func (b *balancer) Serve(w http.ResponseWriter, r *http.Request) {
	service := b.GetNextService()
	service.Serve(w, r)
}

func NewBalancer(services ...service.Service) (*balancer, error) {
	if len(services) < 1 {
		return nil, errors.New("at least one service is required")
	}

	p := pool.NewRoundRobinPool()

	for _, s := range services {
		p.AddService(s)
	}

	return &balancer{
		Pool: p,
	}, nil
}
