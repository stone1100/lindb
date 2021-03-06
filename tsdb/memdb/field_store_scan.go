package memdb

import (
	"fmt"

	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
)

// Scan scans field store for given series id
func (fs *fieldStore) Scan(
	sCtx *series.ScanContext,
	version series.Version,
	seriesID uint32,
	fieldMeta *field.Meta,
	ts *timeSeriesStore,
) {
	queryTimeRange := &sCtx.TimeRange
	worker := sCtx.Worker
	calc := sCtx.IntervalCalc
	interval := sCtx.Interval
	for _, fsStore := range fs.sStoreNodes {
		// check family time is in query time range
		familyTime := fsStore.GetFamilyTime()
		timeRange := &timeutil.TimeRange{
			Start: familyTime,
			End:   calc.CalcFamilyEndTime(familyTime),
		}
		if !queryTimeRange.Overlap(timeRange) {
			continue
		}

		worker.Emit(&series.FieldEvent{
			SeriesID:        seriesID,
			Version:         version,
			FieldIt:         newFStoreIterator(familyTime, fieldMeta, fsStore, ts),
			Interval:        interval,
			FamilyStartTime: familyTime,
		})
	}
}

//////////////////////////////////////////////////////
// fStoreIterator implements FieldIterator
//////////////////////////////////////////////////////
type fStoreIterator struct {
	familyStartTime int64
	sStore          sStoreINTF
	fieldMeta       *field.Meta
	primitiveIt     series.PrimitiveIterator
	ts              *timeSeriesStore

	idx int
}

func newFStoreIterator(
	familyStartTime int64,
	fieldMeta *field.Meta,
	sStore sStoreINTF,
	ts *timeSeriesStore,
) *fStoreIterator {
	return &fStoreIterator{
		familyStartTime: familyStartTime,
		sStore:          sStore,
		ts:              ts,
		fieldMeta:       fieldMeta,
	}
}

func (fsi *fStoreIterator) FieldMeta() field.Meta { return *fsi.fieldMeta }

func (fsi *fStoreIterator) HasNext() bool {
	//FIXME stone for complex field type
	if fsi.idx > 0 {
		return false
	}
	fsi.idx++

	fsi.ts.sl.Lock()
	data, _, _, err := fsi.sStore.Bytes(false)
	fsi.ts.sl.Unlock()

	if err != nil {
		return false
	}
	//FIXME stone1100 set fieldID
	fsi.primitiveIt = series.NewPrimitiveIterator(1, fsi.sStore.AggType(), data)
	return true
}

func (fsi *fStoreIterator) Next() series.PrimitiveIterator {
	return fsi.primitiveIt
}

func (fsi *fStoreIterator) Bytes() ([]byte, error) {
	return nil, fmt.Errorf("fStoreIterator not support Bytes method")
}

func (fsi *fStoreIterator) SegmentStartTime() int64 {
	return fsi.familyStartTime
}
