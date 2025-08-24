package http

import (
	"goQuiz/server/internal/cfg"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

var rlMu sync.Mutex
var rl = map[string]*rate.Limiter{}

func getLimiter(ip string) *rate.Limiter {
	rlMu.Lock()
	defer rlMu.Unlock()

	lim, ok := rl[ip]
	if !ok {
		lim = rate.NewLimiter(rate.Every(10*time.Second), 3)
		rl[ip] = lim
	}
	return lim
}

func rateLimitCreate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := clientIP(r)
		if !getLimiter(ip).Allow() {
			w.Header().Set("Retry-After", "10")
			w.WriteHeader(http.StatusTooManyRequests)
			if cfg.Debug {
				log.Printf("Limiter -> RLCreate | ip = {%s} , allow = false")
			}
			return
		}
		next.ServeHTTP(w, r)
	})
}

func clientIP(r *http.Request) string {
	if xf := r.Header.Get("X-Forwarded-For"); xf != "" {
		parts := strings.Split(xf, ",")
		return strings.TrimSpace(parts[0])
	}
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	if host == "" {
		return r.RemoteAddr
	}
	return host
}
