package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type Proxy struct {
	monolithProxy *httputil.ReverseProxy
	moviesProxy   *httputil.ReverseProxy
}

func NewReverseProxy(monolithURL, moviesURL string) *Proxy {
	return &Proxy{
		monolithProxy: newProxy(monolithURL),
		moviesProxy:   newProxy(moviesURL),
	}
}

func newProxy(target string) *httputil.ReverseProxy {
	url, err := url.Parse(target)
	if err != nil {
		log.Fatalf("Failed to parse target URL: %v", err)
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("Proxy error: %v", err)
		w.WriteHeader(http.StatusBadGateway)
	}
	proxy.Director = func(req *http.Request) {
		req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
		req.Host = url.Host
		req.URL.Host = url.Host
		req.URL.Scheme = url.Scheme
		req.URL.Path = singleJoiningSlash(url.Path, req.URL.Path)
	}
	return proxy
}

func (p *Proxy) Serve(w http.ResponseWriter, r *http.Request, target string) {
	log.Printf("Proxy serving: url %s proxy to %s", r.URL, target)
	switch target {
	case "monolith":
		p.monolithProxy.ServeHTTP(w, r)
	case "movies":
		p.moviesProxy.ServeHTTP(w, r)
	}
}

// Функция для корректного объединения путей
func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
