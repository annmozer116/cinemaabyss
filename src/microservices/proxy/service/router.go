package proxy

import (
	"log"
	"math/rand/v2"
	"net/http"
	"strings"

	config "github.com/cinemaabyss/microservices/proxy/configs"
)

type Router struct {
	config *config.Config
	proxy  *Proxy
}

func NewRouter(cfg *config.Config, proxy *Proxy) *Router {
	return &Router{
		config: cfg,
		proxy:  proxy,
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Проброс health-check
	if req.URL.Path == "/health" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Поиск подходящего маршрута
	for _, route := range r.config.Routes {
		if strings.HasPrefix(req.URL.Path, route.PathPrefix) {
			if route.Target != "monolith" {

				if r.config.GradualMigration && r.config.MigrationPercent > 0 {
					randomNum := rand.IntN(100)

					if randomNum < r.config.MigrationPercent {
						log.Printf("A/B routing, N is %d, route to %s", randomNum, route.Target)
						r.proxy.Serve(w, req, route.Target)
					} else {
						log.Printf("A/B routing, N is %d, route to %s", randomNum, "monolith")
						r.proxy.Serve(w, req, "monolith")
					}
					return
				} else {
					r.proxy.Serve(w, req, route.Target)
				}
			} else {
				r.proxy.Serve(w, req, route.Target)
			}
			return
		}
	}
	// Фоллбэк на монолит
	r.proxy.Serve(w, req, "monolith")
}
