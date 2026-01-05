package api

import (
	"log"
	"net/http"
	"sync"
	"time"
)


var	mu      sync.Mutex
var	clients = make(map[string][]time.Time)

func RateLimit(next http.HandlerFunc, limit int, windowSec int) http.HandlerFunc{
	return func (w http.ResponseWriter, r *http.Request)  {
		ip := r.RemoteAddr
		now := time.Now()
		window := time.Duration(windowSec) *time.Second
		mu.Lock() //locking here 
		times := clients[ip]

		var fresh []time.Time
		for _, t := range times {
			if now.Sub(t)< window {
				fresh = append(fresh, t)
			}
		}
		if len(fresh) >= limit {
			mu.Unlock()
			log.Printf("RATE LIMIT HIT blocked %s on %s ", ip, r.URL.Path)
			http.Error(w, `{"error" : "too many requests"}`, 429)
			return 
		}
		fresh = append(fresh, now)
		clients[ip]  = fresh 
		mu.Unlock()
		log.Printf("REQ OK %s %s (%v)", ip, r.URL.Path, len(fresh))
		next(w,r)
		
	}
}