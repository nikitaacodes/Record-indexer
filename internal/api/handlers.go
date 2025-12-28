package api

import (
	"encoding/json"
	"net/http"
	"record-indexer/internal/model"
	"strconv"
)

func GetRecordHandler(store interface{ ReadAllSafe() ([]model.Record, error) }) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Only GET allowed", 405)
			return
		}

		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", 400)
			return
		}

		records, err := store.ReadAllSafe()
		if err != nil {
			http.Error(w, "data not trustworthy "+err.Error(), 500)
			return
		}

		rec, found := model.Record{}, false
		for _, r := range records {
			if r.ID == id {
				rec, found = r, true
				break
			}
		}
		if !found {
			http.Error(w, "Record not found", 404)
			return
		}

		json.NewEncoder(w).Encode(rec)
	}
}

func ListRecordsHandler(store interface{ ReadAllSafe() ([]model.Record, error) }) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Only GET allowed", 405)
			return
		}

		records, err := store.ReadAllSafe()
		if err != nil {
			http.Error(w, "data not trustworthy "+err.Error(), 500)
			return
		}

		json.NewEncoder(w).Encode(records)
	}
}

func HealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Only GET allowed", 405)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok","msg":"service is alive"}`))
	}
}

func RegisterHandlers(store interface{ ReadAllSafe() ([]model.Record, error) }) {
	http.HandleFunc("/record", GetRecordHandler(store))
	http.HandleFunc("/records", ListRecordsHandler(store))
	http.HandleFunc("/health", HealthHandler())
}
