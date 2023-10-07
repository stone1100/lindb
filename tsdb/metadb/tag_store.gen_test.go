// Code generated by tmpl; DO NOT EDIT.
// https://github.com/benbjohnson/tmpl
//
// Source: int_map_test.tmpl

package metadb

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/roaring"
)

// hack test
func _assertTagStoreData(t *testing.T, keys []uint32, m *TagStore) {
	for _, key := range keys {
		found, highIdx := m.keys.ContainsAndRankForHigh(key)
		assert.True(t, found)
		lowIdx := m.keys.RankForLow(key, highIdx-1)
		assert.True(t, found)
		assert.NotNil(t, m.values[highIdx-1][lowIdx-1])
	}
}

func TestTagStore_Put(t *testing.T) {
	m := NewTagStore()
	m.Put(1, newTagEntry(0))
	m.Put(8, newTagEntry(0))
	m.Put(3, newTagEntry(0))
	m.Put(5, newTagEntry(0))
	m.Put(6, newTagEntry(0))
	m.Put(7, newTagEntry(0))
	m.Put(4, newTagEntry(0))
	m.Put(2, newTagEntry(0))
	// test insert new high
	m.Put(2000000, newTagEntry(0))
	m.Put(2000001, newTagEntry(0))
	// test insert new high
	m.Put(200000, newTagEntry(0))

	_assertTagStoreData(t, []uint32{1, 2, 3, 4, 5, 6, 7, 8, 200000, 2000000, 2000001}, m)
	assert.Equal(t, 11, m.Size())
	assert.Len(t, m.Values(), 3)

	err := m.WalkEntry(func(key uint32, value TagEntry) error {
		return fmt.Errorf("err")
	})
	assert.Error(t, err)

	keys := []uint32{1, 2, 3, 4, 5, 6, 7, 8, 200000, 2000000, 2000001}
	idx := 0
	err = m.WalkEntry(func(key uint32, value TagEntry) error {
		assert.Equal(t, keys[idx], key)
		idx++
		return nil
	})
	assert.NoError(t, err)
}

func TestTagStore_Get(t *testing.T) {
	m := NewTagStore()
	_, ok := m.Get(uint32(10))
	assert.False(t, ok)
	m.Put(1, newTagEntry(0))
	m.Put(8, newTagEntry(0))
	_, ok = m.Get(1)
	assert.True(t, ok)
	_, ok = m.Get(2)
	assert.False(t, ok)
	_, ok = m.Get(0)
	assert.False(t, ok)
	_, ok = m.Get(9)
	assert.False(t, ok)
	_, ok = m.Get(999999)
	assert.False(t, ok)
}

func TestTagStore_Keys(t *testing.T) {
	m := NewTagStore()
	m.Put(1, newTagEntry(0))
	m.Put(8, newTagEntry(0))
	assert.Equal(t, roaring.BitmapOf(1, 8), m.Keys())
}

func TestTagStore_tryOptimize(t *testing.T) {
	m := NewTagStore()
	for i := 0; i < 100; i++ {
		m.Put(uint32(i), newTagEntry(0))
	}
}
