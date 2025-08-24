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
var rlCreate = map[string]*rate.Limiter{}
var rlJoin = map[string]*rate.Limiter{}

func getCreateLimiter(ip string) *rate.Limiter {
	rlMu.Lock()
	defer rlMu.Unlock()

	lim, ok := rlCreate[ip]
	if !ok {
		lim = rate.NewLimiter(rate.Every(10*time.Second), 3)
		rlCreate[ip] = lim
	}
	return lim
}

func getJoinLimiter(ip string) *rate.Limiter {
	rlMu.Lock()
	defer rlMu.Unlock()

	lim, ok := rlJoin[ip]
	if !ok {
		lim = rate.NewLimiter(rate.Limit(5), 20)
		rlJoin[ip] = lim
	}
	return lim
}

func rateLimitCreate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := strings.TrimSpace(clientIP(r))
		if !getCreateLimiter(ip).Allow() {
			w.Header().Set("Retry-After", "10")
			w.WriteHeader(http.StatusTooManyRequests)
			if cfg.Debug {
				log.Printf("Limiter -> RLCreate | ip = {%s} , allow = false", ip)
			}
			return
		}
		next.ServeHTTP(w, r)
	})
}

func rateLimitJoin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := strings.TrimSpace(clientIP(r))
		if !getJoinLimiter(ip).Allow() {
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(http.StatusTooManyRequests)
			if cfg.Debug {
				log.Printf("Limiter -> RLJoin | ip = {%s} , allow = false", ip)
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
	if host == "::1" {
		return "127.0.0.1"
	}
	return host
}
