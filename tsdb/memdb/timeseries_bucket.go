package memdb

import (
	"sync"
)

type TimeSeriesBucket interface {
	GetOrCreateFStore(seriesID uint32, createFn func() (fStoreINTF, error)) (fStoreINTF, error)
	Size() int
}

type timeSeriesBucket struct {
	stores *TimeSeriesStore // series id => field store

	lock sync.RWMutex
}

func NewTimeSeriesBucket() TimeSeriesBucket {
	tsb := &timeSeriesBucket{
		stores: NewTimeSeriesStore(),
	}
	return tsb
}

func (tsb *timeSeriesBucket) Size() int {
	tsb.lock.Lock()
	defer tsb.lock.Unlock()

	return tsb.stores.Size()
}

func (tsb *timeSeriesBucket) GetOrCreateFStore(seriesID uint32, createFn func() (fStoreINTF, error)) (fStoreINTF, error) {
	tsb.lock.Lock()
	defer tsb.lock.Unlock()

	fStore, ok := tsb.stores.Get(seriesID)
	if ok {
		return fStore, nil
	}
	// create field store
	fs, err := createFn()
	if err != nil {
		return nil, err
	}
	tsb.stores.Put(seriesID, fs)
	return fs, nil
}
