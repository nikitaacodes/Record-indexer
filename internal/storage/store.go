package storage

import model "record-indexer/internal/model"

type Store interface {
	Append(record model.Record) error
	ReadAllSafe() ([]model.Record,error)
}