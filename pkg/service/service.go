package service

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type Service interface {
	GetURL() *url.URL
	Serve(http.ResponseWriter, *http.Request)
	IsAlive() bool
	SetAlive(bool)
}

type service struct {
	url   *url.URL
	alive bool

	reverseProxy *httputil.ReverseProxy

	mux sync.RWMutex
}

func (s *service) SetAlive(alive bool) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.alive = alive
}

func (s *service) IsAlive() bool {
	s.mux.RLock()
	defer s.mux.RUnlock()
	return s.alive
}

func (s *service) Serve(w http.ResponseWriter, r *http.Request) {
	s.reverseProxy.ServeHTTP(w, r)
}

func (s *service) GetURL() *url.URL {
	return s.url
}

func NewService(url *url.URL) *service {
	ns := &service{
		url:   url,
		alive: true,
	}

	rp := httputil.NewSingleHostReverseProxy(url)
	rp.ErrorHandler = ns.errorHandler

	ns.reverseProxy = rp

	return ns
}

func (s *service) errorHandler(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Error in service %s: %v\n", s.url, err)

	s.SetAlive(false)
	http.Error(w, "Service not available", http.StatusServiceUnavailable)
}
