package storage

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"record-indexer/internal/model"
)

// VerifyHash
//checks if the stored hash matches actual data hash
func VerifyHash(r model.Record) bool {
	h := sha256.Sum256([]byte(r.Data))
	expected := hex.EncodeToString(h[:])
	return r.Hash == expected
}

type DiskStore struct {
	filepath string //  private field
}


func NewDiskStore(path string) *DiskStore {
	return &DiskStore{filepath: path}
}

//for write operations
func (ds *DiskStore) Append(record model.Record) error {
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
		records = append(records, r)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}


	var safe []model.Record
	lastID := -1

	for _, r := range records {
		if r.ID <= lastID {
			return nil, errors.New("records are out of order")
		}
		if !VerifyHash(r) {
			return nil, errors.New("record data is corrupted (hash mismatch)")
		}
		lastID = r.ID
		safe = append(safe, r)
	}
expectedID:= 0
for _, r := range records{
	expectedID++
	if r.ID != expectedID{
		return nil, errors.New("missing record")
	}

}
	return safe, nil
}
