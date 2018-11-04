package database

import (
	"fmt"
	"sync"

	"github.com/elastos/Elastos.ELA.SPV/database"
	"github.com/elastos/Elastos.ELA.SPV/util"
	"github.com/elastos/Elastos.ELA.Utility/common"
)

type headers struct {
	mutex sync.Mutex
	best  *util.Header
	store map[common.Uint256]*util.Header
}

// Save a header to database
func (h *headers) Put(header *util.Header, newTip bool) error {
	h.mutex.Lock()
	hash := header.Hash()
	h.store[hash] = header
	if newTip {
		h.best = header
	}
	h.mutex.Unlock()
	return nil
}

// Get previous block of the given header
func (h *headers) GetPrevious(header *util.Header) (*util.Header, error) {
	h.mutex.Lock()
	previous, ok := h.store[header.Previous()]
	h.mutex.Unlock()
	if !ok {
		return nil, fmt.Errorf("previous header not exist")
	}
	return previous, nil
}

// Get full header with it's hash
func (h *headers) Get(hash *common.Uint256) (*util.Header, error) {
	h.mutex.Lock()
	header, ok := h.store[*hash]
	h.mutex.Unlock()
	if !ok {
		return nil, fmt.Errorf("header not exist")
	}
	return header, nil
}

// Get the header on chain tip
func (h *headers) GetBest() (*util.Header, error) {
	if h.best == nil {
		return nil, fmt.Errorf("best header not exist")
	}
	return h.best, nil
}

// Clear delete all data in database.
func (h *headers) Clear() error {
	h.best = nil
	h.mutex.Lock()
	h.store = make(map[common.Uint256]*util.Header)
	h.mutex.Unlock()
	return nil
}

// Close database.
func (h *headers) Close() error {
	return nil
}

func NewDatabase() database.ChainStore {
	return database.NewHeadersOnlyChainDB(&headers{
		store: make(map[common.Uint256]*util.Header),
	})
}
