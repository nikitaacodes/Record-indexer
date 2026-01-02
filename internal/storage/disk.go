package storage

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"record-indexer/internal/model"
	"time"
)

// VerifyHash
//checks if the stored hash matches actual data hash
func verifyHash(data, hash string) bool {
	h := sha256.Sum256([]byte(data))
	expected := hex.EncodeToString(h[:])
	return expected == hash
}

type DiskStore struct {
	filepath string //  private field
}


func NewDiskStore(path string) *DiskStore {
	return &DiskStore{filepath: path}
}

//for write operations
func (ds *DiskStore) Append(record model.Record) error {
	record.Status = "valid"
	record.Timestamp = time.Now().Format(time.RFC3339)

	file, err := os.OpenFile(ds.filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644) //0644 file permission for only owner can read 
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	bytes, err := json.Marshal(record)
	if err != nil {
		return err
	}

	_, err = writer.WriteString(string(bytes) + "\n")
	if err != nil {
		return err
	}

	return writer.Flush() //pushed on disk 
}

//for read operations
func (ds *DiskStore) ReadAllSafe() ([]model.Record, error) {
	file, err := os.Open(ds.filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var records []model.Record
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		var r model.Record
		err := json.Unmarshal(scanner.Bytes(), &r)
		if err != nil {
			return nil, err
		}
		if !verifyHash(r.Data,r.Hash){
			r.Status = "corrupted"
		}else{
			r.Status = "valid"
		}
		records = append(records, r)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	lastID := 0

	for _, r := range records {
		if r.ID != lastID+1 {
			return nil, errors.New("missing record")
		}
		if r.ID <= lastID{
			return nil, errors.New("records are out of order")
		}
		lastID = r.ID
	
	}
	
	return records, nil
}

func(ds *DiskStore) IntegrityStatus() map[string]interface{}{
 records, _ := ds.ReadAllSafe()
 total := len(records)
 valid,corrupted := 0, 0
for _, r := range records{
	if r.Status == "valid"{
		valid++
	}else{
		corrupted++
	}
}
return map[string]interface{}{
	"total_records" : total,
	"valid_records" : valid,
	"corrupted_records": corrupted,
	"last_checked": time.Now().Format(time.RFC3339),
}
}
