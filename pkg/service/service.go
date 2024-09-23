package service

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Service interface {
	GetURL() *url.URL
	Serve(http.ResponseWriter, *http.Request)
}

type service struct {
	url   *url.URL
	Alive bool

	ReverseProxy *httputil.ReverseProxy
}

func (s *service) Serve(w http.ResponseWriter, r *http.Request) {
	s.ReverseProxy.ServeHTTP(w, r)
}

func (s *service) GetURL() *url.URL {
	return s.url
}

func NewService(url *url.URL) service {
	rp := httputil.NewSingleHostReverseProxy(url)

	return service{
		url:          url,
		Alive:        true,
		ReverseProxy: rp,
	}
}
