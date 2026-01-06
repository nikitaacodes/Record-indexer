package api

import (
	"log"
	"net/http"
	"os"
)

func BasicAuth(wrapFunc http.HandlerFunc) http.HandlerFunc{

	user := os.Getenv("BASIC_AUTH_USER")
	pass:= os.Getenv("BASIC_AUTH_PASS")

	return func (w http.ResponseWriter, r *http.Request)  {
		//bypasssing public route here
		if r.URL.Path == "/health" {
			wrapFunc(w,r)
			return
		}
		u, p, ok:= r.BasicAuth() //username, passs, ok bool
		if !ok || u != user || p != pass {
			log.Printf("AUTH FAILED %s %s" , r.Method, r.URL.Path)
			http.Error(w, `{"error" : "unautorized"}`, 401)
			return 

		}
		log.Printf("AUTH SUCCESS %s %s", r.Method, r.URL.Path)
		wrapFunc(w,r)
		
		
	}
}