package model

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

type Record struct {
	ID   int    `json:"id"`
	Timestamp string `json:"timestamp"`
	Data string `json:"data"`
	Hash string `json:"hash"`
    Status string `json:"status"`


}

func NewRecord(id int, data string) Record {
	h := sha256.Sum256([]byte(data))
	return Record{
		ID: id, 
		Timestamp: time.Now().Format(time.RFC3339),
		Data: data,
		Hash: hex.EncodeToString(h[:]),
		Status: "valid",
	}
}