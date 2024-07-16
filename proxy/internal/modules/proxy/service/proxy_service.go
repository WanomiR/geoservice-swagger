package service

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type ProxyReverse struct {
	host string
	port string
}

func NewProxyReverse(host, port string) *ProxyReverse {
	return &ProxyReverse{
		host: host,
		port: port,
	}
}

// localhost:1313/static -> hugo
// localhost:1313/api -> api

func (rp *ProxyReverse) ProxyReverse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api") || strings.HasPrefix(r.URL.Path, "/swagger") {
			next.ServeHTTP(w, r)
			return
		}
		link := fmt.Sprintf("http://%s:%s", rp.host, rp.port)
		uri, _ := url.Parse(link)

		if uri.Host == r.Host {
			next.ServeHTTP(w, r)
			return
		}
		r.Header.Set("Reverse-Proxy", "true")

		proxy := httputil.ReverseProxy{Director: func(r *http.Request) {
			r.URL.Scheme = uri.Scheme
			r.URL.Host = uri.Host
			r.URL.Path = uri.Path + r.URL.Path
			r.Host = uri.Host
		}}

		proxy.ServeHTTP(w, r)
	})
}
