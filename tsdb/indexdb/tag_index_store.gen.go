// Code generated by tmpl; DO NOT EDIT.
// https://github.com/benbjohnson/tmpl
//
// Source: int_map.tmpl

package indexdb

import (
	"github.com/lindb/roaring"
)

// TagIndexStore represents int map using roaring bitmap
type TagIndexStore struct {
	putCount int             // insert count
	keys     *roaring.Bitmap // store all keys
	values   [][]TagIndex    // store all values by high/low key
}

// NewTagIndexStore creates a int map
func NewTagIndexStore() *TagIndexStore {
	return &TagIndexStore{
		keys: roaring.New(),
	}
}

// Get returns value by key, if exist returns it, else returns nil, false
func (m *TagIndexStore) Get(key uint32) (TagIndex, bool) {
	if len(m.values) == 0 {
		return nil, false
	}
	// get high index
	found, highIdx := m.keys.ContainsAndRankForHigh(key)
	if !found {
		return nil, false
	}
	// get low index
	found, lowIdx := m.keys.ContainsAndRankForLow(key, highIdx-1)
	if !found {
		return nil, false
	}
	return m.values[highIdx-1][lowIdx-1], true
}

// Put puts the value by key
func (m *TagIndexStore) Put(key uint32, value TagIndex) {
	defer m.tryOptimize()
	if len(m.values) == 0 {
		// if values is empty, append new low container directly
		m.values = append(m.values, []TagIndex{value})

		m.keys.Add(key)
		return
	}
	found, highIdx := m.keys.ContainsAndRankForHigh(key)
	if !found {
		// high container not exist, insert it
		stores := m.values
		// insert operation, insert high values
		stores = append(stores, nil)
		copy(stores[highIdx+1:], stores[highIdx:len(stores)-1])
		stores[highIdx] = []TagIndex{value}
		m.values = stores

		m.keys.Add(key)
		return
	}
	// high container exist
	lowIdx := m.keys.RankForLow(key, highIdx-1)
	stores := m.values[highIdx-1]
	// insert operation
	stores = append(stores, nil)
	copy(stores[lowIdx+1:], stores[lowIdx:len(stores)-1])
	stores[lowIdx] = value
	m.values[highIdx-1] = stores

	m.keys.Add(key)
}

// tryOptimize optimizes the roaring bitmap when inserted in every 100
func (m *TagIndexStore) tryOptimize() {
	m.putCount++
	if m.putCount%100 == 0 {
		m.keys.RunOptimize()
	}
}

// Keys returns the all keys
func (m *TagIndexStore) Keys() *roaring.Bitmap {
	return m.keys
}

// Values returns the all values
func (m *TagIndexStore) Values() [][]TagIndex {
	return m.values
}

// Size returns the size of keys
func (m *TagIndexStore) Size() int {
	return int(m.keys.GetCardinality())
}

// WalkEntry walks each kv entry via fn.
func (m *TagIndexStore) WalkEntry(fn func(key uint32, value TagIndex) error) error {
	values := m.values
	keys := m.keys
	highKeys := keys.GetHighKeys()
	for highIdx, highKey := range highKeys {
		hk := uint32(highKey) << 16
		lowValues := values[highIdx]
		lowContainer := keys.GetContainerAtIndex(highIdx)
		it := lowContainer.PeekableIterator()
		idx := 0
		for it.HasNext() {
			lowKey := it.Next()
			value := lowValues[idx]
			idx++
			if err := fn(uint32(lowKey&0xFFFF)|hk, value); err != nil {
				return err
			}
		}
	}
	return nil
}
