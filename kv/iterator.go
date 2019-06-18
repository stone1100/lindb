package kv

// iterator iterates over a store's key/value pairs in key order.
type iterator interface {
	// Next moves the iterator to the next key/value pair.
	// It returns false if the iterator is exhausted.
	Next() bool
	// Key returns the key of the current key/value pair
	Key() uint32
	// Value returns the value of the current key/value pair
	Value() []byte
}
