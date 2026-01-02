package api

import (
	"encoding/json"
	"net/http"
	"record-indexer/internal/model"
	"record-indexer/internal/storage"
	"strconv"
)
const errOnlyGET = "only GET allowed"
func GetRecordHandler(store storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, errOnlyGET, 405)
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
		for _, record := range records {
			if record.ID == id {
				rec, found = record, true
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

func ListRecordsHandler(store storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, errOnlyGET, 405)
			return
		}

		records, err := store.ReadAllSafe()
		if err != nil {
			http.Error(w, "data not trustworthy "+err.Error(), 500)
			return
		}

		limitStr := r.URL.Query().Get("limit")
		offsetStr := r.URL.Query().Get("offset")

		limit, _ := strconv.Atoi(limitStr)
		offset, _ := strconv.Atoi(offsetStr)

		if limit <= 0 {
			limit = len(records)
		}
		if offset < 0 {
			offset = 0
		}
		end := offset + limit
		if end > len(records) {
			end = len(records)
		}

		json.NewEncoder(w).Encode(records[offset:end])
	}
}
 func HealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, errOnlyGET, 405)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok","msg":"service is alive"}`))
	}
}
func IntegrityHandler(store storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if ds, ok := store.(*storage.DiskStore); ok {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(ds.IntegrityStatus())
			return
		}
		http.Error(w, "integrity route not supported", 500)
	}
}

func RegisterHandlers(store storage.Store) {
	http.HandleFunc("/record", GetRecordHandler(store))
	http.HandleFunc("/records", ListRecordsHandler(store))
	http.HandleFunc("/health", HealthHandler())
	http.HandleFunc("/integrity", IntegrityHandler(store))
}
