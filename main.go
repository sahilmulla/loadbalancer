package main

import (
	"log"
	"net/http"
	"net/url"

	"github.com/sahilmulla/loadbalancer/pkg/balancer"
	"github.com/sahilmulla/loadbalancer/pkg/pool"
	"github.com/sahilmulla/loadbalancer/pkg/service"
)

func main() {
	url1, err := url.Parse("http://localhost:8081")
	if err != nil {
		log.Fatalln(err)
	}
	service1 := service.NewService(url1)

	url2, err := url.Parse("http://localhost:8082")
	if err != nil {
		log.Fatalln(err)
	}
	service2 := service.NewService(url2)

	p := pool.NewRoundRobinPool()
	p.AddService(&service1)
	p.AddService(&service2)

	lb := balancer.NewBalancer(p)

	server := http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(lb.Serve),
	}

	log.Fatal(server.ListenAndServe())
}
