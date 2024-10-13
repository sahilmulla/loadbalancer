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
	p := pool.NewRoundRobinPool()

	lb := balancer.NewBalancer(p)

	server := http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(lb.Serve),
	}

	go func() {
		internal := http.NewServeMux()

		internal.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				urlStr := r.URL.Query().Get("url")
				serviceUrl, err := url.Parse(urlStr)
				if err != nil {
					log.Println(err)
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				s := service.NewService(serviceUrl)

				if err := p.AddService(s); err != nil {
					log.Println(err)
					w.WriteHeader(http.StatusUnprocessableEntity)
					return
				}

				w.WriteHeader(http.StatusNoContent)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		})

		log.Fatalln(http.ListenAndServe(":9000", internal))
	}()

	log.Fatal(server.ListenAndServe())
}
