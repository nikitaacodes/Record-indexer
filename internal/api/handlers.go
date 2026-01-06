package api

import (
	"encoding/json"
	"log"
	"net/http"
	"record-indexer/internal/model"
	"record-indexer/internal/storage"
	"strconv"
	"time"
)

const errOnlyGET = "Only GET allowed"

func GetRecordHandler(store storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		if r.Method != "GET" {
			http.Error(w, errOnlyGET, 405)
			log.Printf("%s %s → 405 (%v)", r.Method, r.URL.Path, time.Since(start))
			return
		}

		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", 400)
			log.Printf("%s %s → 400 (%v)", r.Method, r.URL.Path, time.Since(start))
			return
		}

		records, err := store.ReadAllSafe()
		if err != nil {
			http.Error(w, "data not trustworthy "+err.Error(), 500)
			log.Printf("%s %s → 500 (%v)", r.Method, r.URL.Path, time.Since(start))
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
			log.Printf("%s %s → 404 (%v)", r.Method, r.URL.Path, time.Since(start))
			return
		}

		_ = json.NewEncoder(w).Encode(rec)
		log.Printf("%s %s → 200 (%v)", r.Method, r.URL.Path, time.Since(start))
	}
}

func ListRecordsHandler(store storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		if r.Method != "GET" {
			http.Error(w, errOnlyGET, 405)
			log.Printf("%s %s → 405 (%v)", r.Method, r.URL.Path, time.Since(start))
			return
		}

		records, err := store.ReadAllSafe()
		if err != nil {
			http.Error(w, "data not trustworthy "+err.Error(), 500)
			log.Printf("%s %s → 500 (%v)", r.Method, r.URL.Path, time.Since(start))
			return
		}

		// Pagination logic
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

		paged := records[offset:end]
		_ = json.NewEncoder(w).Encode(paged)

		log.Printf("%s %s → 200 (%d returned) (%v)", r.Method, r.URL.Path, len(paged), time.Since(start))
	}
}

func HealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		if r.Method != "GET" {
			http.Error(w, errOnlyGET, 405)
			log.Printf("%s %s → 405 (%v)", r.Method, r.URL.Path, time.Since(start))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok","msg":"service is alive"}`))
		log.Printf("%s %s → 200 (%v)", r.Method, r.URL.Path, time.Since(start))
	}
}

func IntegrityHandler(store storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		if r.Method != "GET" {
			http.Error(w, errOnlyGET, 405)
			log.Printf("%s %s → 405 (%v)", r.Method, r.URL.Path, time.Since(start))
			return
		}

		records, _ := store.ReadAllSafe()
		total := len(records)
		valid := 0
		corrupted := 0

		for _, record := range records {
			if record.Status == "valid" {
				valid++
			} else {
				corrupted++
			}
		}

		resp := map[string]interface{}{
			"total_records":     total,
			"valid_records":     valid,
			"corrupted_records": corrupted,
			"last_checked":      time.Now().Format(time.RFC3339),
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)

		log.Printf("%s %s → 200 (%v) (checked %d records)", r.Method, r.URL.Path, time.Since(start), total)
	}
}

func RegisterHandlers(store storage.Store) {
	const limit = 10 
	const window = 30
	http.HandleFunc("/record", BasicAuth(RateLimit(GetRecordHandler(store), limit, window)))
	http.HandleFunc("/records", BasicAuth(RateLimit(ListRecordsHandler(store), limit, window)))
	http.HandleFunc("/health", HealthHandler())      		//public 
	http.HandleFunc("/integrity/status", BasicAuth(RateLimit(IntegrityHandler(store), limit, window)))
}
