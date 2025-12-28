package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"record-indexer/internal/api"

	"record-indexer/internal/storage"
	"syscall"
)
func main() {
	store := storage.NewDiskStore("records.log")
//appended at first 
//    store.Append(model.NewRecord(1, "hello"))
//     store.Append(model.NewRecord(2, "world"))
//     store.Append(model.NewRecord(3, "trust-machines-are-boring"))
//     store.Append(model.NewRecord(4, "go-is-my-new-friend"))
//     store.Append(model.NewRecord(5, "immutable-record-vault"))
//     log.Println("sample dataset written")

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
