package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"record-indexer/internal/api"
	"record-indexer/internal/storage"
	"syscall"

	"github.com/joho/godotenv"
)
func main() {
	err := godotenv.Load()
	if err!= nil{
		log.Println("no .env file found")
	}
	store := storage.NewDiskStore("records.log")
// appended at first 
//    store.Append(model.NewRecord(1, "hello"))
//     store.Append(model.NewRecord(2, "world"))
//     store.Append(model.NewRecord(3, "learning go was not that difficult"))
//     store.Append(model.NewRecord(4, "i hope i'll get selected "))
//     store.Append(model.NewRecord(5, "Hope for the best"))
//     log.Println("sample dataset written")

log.Println("loaded user:" , os.Getenv("BASIC_AUTH_USER")) //log to check
    api.RegisterHandlers(store)
	srv := &http.Server{Addr: ":8080"}

	log.Println("server started on :8080") 

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("server shutting down...")
	srv.Close()
	log.Println("server exited cleanly")
}
