package model

import (
	"crypto/sha256"
	"encoding/hex"
)

type Record struct {
	ID   int    `json:"id"`
	Data string `json:"data"`
	Hash string `json:"hash"`
}

func NewRecord(id int, data string) Record {
	h := sha256.Sum256([] byte(data))
	return Record{
		ID: id, 
		Data: data,
		Hash: hex.EncodeToString(h[:]),
	}
}