package balancer

import (
	"errors"
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
	service, err := b.pool.GetNextService()
	if err != nil {
		if errors.Is(err, pool.ErrAllDown) {
			http.Error(w, "Service not available", http.StatusServiceUnavailable)
		} else {
			log.Fatalln(err)
		}
		return
	}
	service.Serve(w, r)
}

func NewBalancer(pool pool.Pool) balancer {
	return balancer{
		pool: pool,
	}
}
