package memdb

import (
	"fmt"

	"github.com/lindb/common/models"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
)

type exemplar struct {
	traceID  []byte
	spanID   []byte
	duration int64
}

type exemplarStore struct {
	writer *stream.BufferWriter
	id     field.ID

	store map[uint16]exemplar
}

func newExemplarStore(id field.ID) fStoreINTF {
	return &exemplarStore{
		id:    id,
		store: make(map[uint16]exemplar),
	}
}

func (*exemplarStore) Capacity() int {
	return 10
}

func (fs *exemplarStore) FlushFieldTo(tableFlusher metricsdata.Flusher, fieldMeta field.Meta, flushCtx *flushContext) error {
	panic("unimplemented")
}

func (fs *exemplarStore) GetFieldID() field.ID {
	return fs.id
}

func (fs *exemplarStore) Load(ctx *flow.DataLoadContext, seriesIdxFromQuery uint16, fieldIdx int, fieldType field.Type, slotRange timeutil.SlotRange) {
	ctx.DownSampling(slotRange, seriesIdxFromQuery, fieldIdx, fs)
}

func (fs *exemplarStore) GetValue(slot uint16) (float64, bool) {
	return 0, false
}

func (fs *exemplarStore) GetExemplar(slot uint16) (exemplar *models.Exemplar, ok bool) {
	e, ok := fs.store[slot]
	if !ok {
		return nil, false
	}
	return &models.Exemplar{
		TraceID:  string(e.traceID),
		SpanID:   string(e.spanID),
		Duration: e.duration,
	}, true
}

func (fs *exemplarStore) Write(fieldType field.Type, slotIndex uint16, value float64) {
	// do nothing
}

// WriteExemplar implements fStoreINTF
func (fs *exemplarStore) WriteExemplar(slotIndex uint16, traceID []byte, spanID []byte, duration int64) {
	fmt.Println("write....8")
	fs.store[slotIndex] = exemplar{
		traceID:  traceID,
		spanID:   spanID,
		duration: duration,
	}
}

// func (fs *exemplarStore) WriteExemplar(slotIndex uint16, traceID, spanID []byte, duration int64) {
// offset := fs.writer.Len()
// _, err := fs.writer.Write(traceID)
// if err != nil {
// 	fmt.Println(err)
// 	return
// }
// _, err = fs.writer.Write(spanID)
// if err != nil {
// 	fmt.Println(err)
// 	return
// }
// if fs.buf[markOffset+1] == 0 {
// 	// no data written before
// 	fs.writeFirstPoint(slotIndex, value)
// 	return
// }
//
// startTime := fs.getStart()
// if slotIndex < startTime || slotIndex > startTime+fs.timeWindow()-1 {
// 	// if current slot time out of current time window, need compress block data, start new time window
// 	fs.compact(fieldType, startTime)
//
// 	// write first point after compact
// 	fs.writeFirstPoint(slotIndex, value)
// 	return
// }
//
// // write data in current write buffer
// delta := slotIndex - startTime
// pos, markIdx, flagIdx := fs.position(delta)
// if fs.buf[markOffset+markIdx]&flagIdx != 0 {
// 	// there is same point of same time slot
// 	oldValue := math.Float64frombits(binary.LittleEndian.Uint64(fs.buf[pos:]))
// 	value = fieldType.AggType().Aggregate(oldValue, value)
// } else {
// 	// new data for time slot
// 	fs.buf[endOffset] = byte(delta)
// 	fs.buf[markOffset+markIdx] |= flagIdx // mark value exist
// }
// // finally, write value into the body of current write buffer
// binary.LittleEndian.PutUint64(fs.buf[pos:], math.Float64bits(value))
// }
