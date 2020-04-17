package store

import (
	"crypto/sha256"
	"encoding/binary"
	"github.com/maxim-kuderko/mosyne/entities"
	"github.com/wangjia184/sortedset"
	"hash"
	"sync"
)

type hashedKey uint64
type InMem struct {
	shards []map[hashedKey]*sortedSet

	mu []*sync.RWMutex
}

type sortedSet struct {
	*sortedset.SortedSet
	mu sync.RWMutex
}

func NewInMemStore(opt *Options) *InMem {
	locks := make([]*sync.RWMutex, 2048)
	for i := 0; i < 2048; i++ {
		locks[i] = &sync.RWMutex{}
	}
	return &InMem{
		shards: make([]map[hashedKey]*sortedSet, 2048),
		mu:     locks,
	}
}

func (i *InMem) ZGet(get entities.ZGetRequest) entities.ZGetResponse {
	k := i.hashOfKey(get.Key)
	mu := i.mu[k%2048]
	mu.RLock()
	storage := i.getSortedSet(k)
	mu.RUnlock()
	storage.mu.RLock()
	defer storage.mu.RUnlock()
	resp := storage.GetByScoreRange(sortedset.SCORE(get.ScoreMin), sortedset.SCORE(get.ScoreMax), &sortedset.GetByScoreRangeOptions{
		Limit:        0,
		ExcludeStart: false,
		ExcludeEnd:   false,
	})
	n := len(resp)
	output := entities.ZGetResponse{
		Values: make([]entities.ZGetStruct, n),
		Error:  nil,
	}
	for i := 0; i < n; i++ {
		output.Values[i] = entities.ZGetStruct{
			Value: resp[i].Value,
			Score: float64(resp[i].Score()),
		}
	}
	return output
}

var hashes = sync.Pool{New: func() interface{} {
	return sha256.New()
}}

func (i *InMem) ZSet(set entities.ZSetRequest) entities.ZSetResponse {
	k := i.hashOfKey(set.Key)
	mu := i.mu[k%2048]
	mu.Lock()
	storage := i.getSortedSet(k)
	mu.Unlock()
	storage.mu.Lock()
	defer storage.mu.Unlock()
	storage.AddOrUpdate(set.Key, sortedset.SCORE(set.Score), set.Value)
	return entities.ZSetResponse{
		Value: set.Value,
		Error: nil,
	}
}

func (i *InMem) getSortedSet(k hashedKey) *sortedSet {
	shard := i.shards[k%2048]
	if shard == nil {
		shard = map[hashedKey]*sortedSet{}
		i.shards[k%2048] = shard
	}
	storage, ok := shard[k]
	if !ok {
		storage = &sortedSet{
			SortedSet: sortedset.New(),
			mu:        sync.RWMutex{},
		}
		shard[k] = storage
	}
	return storage
}

func (i *InMem) hashOfKey(key string) hashedKey {
	h := hashes.Get().(hash.Hash)
	h.Write([]byte(key))
	b := h.Sum(nil)
	h.Reset()
	hashes.Put(h)
	return hashedKey(binary.LittleEndian.Uint64(b))
}
