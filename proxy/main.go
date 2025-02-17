package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/go-chi/chi"
)

type Router struct {
	r *chi.Mux
}

func main() {
	http.ListenAndServe(":8080", getProxyRouter("http://hugo", ":1313").r)
}

func getProxyRouter(host, port string) *Router {
	r := &Router{r: chi.NewRouter()}

	r.r.Use(NewReverseProxy(host, port).ReverseProxy)

	r.r.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from API"))
	})
	return r
}

type ReverseProxy struct {
	host string
	port string
}

func NewReverseProxy(host, port string) *ReverseProxy {
	return &ReverseProxy{
		host: host,
		port: port,
	}
}

// Если ресурс имеет префикс /api/, то запрос должен выдавать текст «Hello from API». Все остальные запросы должны перенаправляться на http://hugo:1313 (сервер hugo).
func (rp *ReverseProxy) ReverseProxy(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/" {
			next.ServeHTTP(w, r)
			return
		}

		targetURL, _ := url.Parse(rp.host + rp.port)

		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		proxy.ServeHTTP(w, r)
	})
}
