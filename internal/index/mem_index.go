package index

import "record-indexer/internal/model"

type MemIndex struct {
    m map[int]model.Record
}

func New() *MemIndex {
    return &MemIndex{m: make(map[int]model.Record)}
}

func (mi *MemIndex) Load(records []model.Record) {
    for _, r := range records {
        mi.m[r.ID] = r
    }
}

func (mi *MemIndex) Get(id int) (model.Record, bool) {
    r, ok := mi.m[id]
    return r, ok
}