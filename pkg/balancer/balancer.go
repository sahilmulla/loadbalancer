package balancer

import (
	"log"
	"net/http"

	"github.com/sahilmulla/loadbalancer/pkg/pool"
)

type Balancer interface {
	Serve(http.ResponseWriter, *http.Request)
}

type balancer struct {
	pool pool.Pool
}

func (b *balancer) Serve(w http.ResponseWriter, r *http.Request) {
	log.Println("serving: ", r.URL)
	service := b.pool.GetNextService()
	service.Serve(w, r)
}

func NewBalancer(pool pool.Pool) balancer {
	return balancer{
		pool: pool,
	}
}
