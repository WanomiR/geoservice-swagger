package service

import "net/http"

type ProxyReverser interface {
	ProxyReverse(next http.Handler) http.Handler
}
