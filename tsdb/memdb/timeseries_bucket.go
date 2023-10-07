package memdb

import (
	"sync"
)

type TimeSeriesBucket interface {
	GetOrCreateFStore(seriesID uint32, createFn func() (fStoreINTF, error)) (fStoreINTF, error)
}

type timeSeriesBucket struct {
	stores *TimeSeriesStore // series id => field store

	lock sync.RWMutex
}

func NewTimeSeriesBucket() TimeSeriesBucket {
	return &timeSeriesBucket{
		stores: NewTimeSeriesStore(),
	}
}

func (tsb *timeSeriesBucket) GetOrCreateFStore(seriesID uint32, createFn func() (fStoreINTF, error)) (fStoreINTF, error) {
	tsb.lock.Lock()
	defer tsb.lock.Unlock()

	fStore, ok := tsb.stores.Get(seriesID)
	if ok {
		return fStore, nil
	}

	// create field store
	fStore, err := createFn()
	if err != nil {
		return nil, err
	}
	tsb.stores.Put(seriesID, fStore)
	return fStore, nil
}
