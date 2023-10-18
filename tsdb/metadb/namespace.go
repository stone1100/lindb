package metadb

import (
	"io"
	"sync"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/kv/table"
)

type NamespaceReader interface {
}

type NamespaceWriter interface {
	io.Closer
	StartBucket()
	WriteNamespace(namespace []byte, namespaceID uint32)
	FlushBucket(namespaceBucket uint32)
}

type namespaceWriter struct {
	kvWriter  table.StreamWriter
	kvFlusher kv.Flusher
}

func NewNamespaceWriter(kvFlusher kv.Flusher) (NamespaceWriter, error) {
	kvWriter, err := kvFlusher.StreamWriter()
	if err != nil {
		return nil, err
	}
	return &namespaceWriter{
		kvWriter:  kvWriter,
		kvFlusher: kvFlusher,
	}, nil
}

func (f *namespaceWriter) StartBucket() {

}

func (f *namespaceWriter) WriteNamespace(namespace []byte, namespaceID uint32) {

}

func (f *namespaceWriter) FlushBucket(namespaceBucket uint32) {
}

func (f *namespaceWriter) Close() error {
	return nil
}

type NamespaceStore interface {
}

type namespaceStore struct {
	kvStore kv.Store
	family  kv.Family

	mutable   *sync.Map
	immutable *sync.Map

	seq atomic.Uint32

	lock sync.RWMutex
}

func NewNamespaceStore(kvStore kv.Store) (NamespaceStore, error) {
	family, err := kvStore.CreateFamily("ns", kv.FamilyOption{})
	if err != nil {
		return nil, err
	}
	return &namespaceStore{
		kvStore: kvStore,
		family:  family,
	}, nil
}

func (nss *namespaceStore) GetNamespaceID(namespace string) uint32 {
	return 0
}

func (nss *namespaceStore) Put(namespace []byte) {
	id := nss.seq.Inc()
	nss.mutable.Store(namespace, id)
}

func (nss *namespaceStore) Flush() error {
	nss.lock.Lock()
	nss.immutable = nss.mutable
	nss.mutable = &sync.Map{}
	nss.lock.Unlock()

	return nil
}
