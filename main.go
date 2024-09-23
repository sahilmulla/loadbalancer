package main

import (
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/sahilmulla/loadbalancer/pkg/balancer"
	"github.com/sahilmulla/loadbalancer/pkg/pool"
	"github.com/sahilmulla/loadbalancer/pkg/service"
)

func main() {
	serverUrlStrs := "http://localhost:8081,http://localhost:8082,http://localhost:8083"

	p := pool.NewRoundRobinPool()

	for _, urlStr := range strings.Split(serverUrlStrs, ",") {
		url, err := url.Parse(urlStr)
		if err != nil {
			log.Fatalln(err)
		}
		s := service.NewService(url)
		p.AddService(s)
	}

	lb := balancer.NewBalancer(p)
	server := http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(lb.Serve),
	}
	log.Fatal(server.ListenAndServe())
}
