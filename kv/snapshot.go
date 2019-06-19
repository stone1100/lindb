package kv

import (
	"sync/atomic"

	"github.com/eleme/lindb/kv/version"
)

// Snapshot current family metadata, for reading data,
// snaphost instance must close after unuse
type Snapshot struct {
	version *version.Version
	closed  int32
}

// newSnapshot new snapshot instance
func newSnapshot(version *version.Version) *Snapshot {
	return &Snapshot{
		version: version,
	}
}

// Get finds values based on given key
func (s *Snapshot) Get(key uint32) {

}

// Close releases related resources
func (s *Snapshot) Close() {
	if atomic.CompareAndSwapInt32(&s.closed, 0, 1) {
		s.version.Release()
	}
}
