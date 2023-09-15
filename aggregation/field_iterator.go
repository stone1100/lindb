// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package aggregation

import (
	"fmt"
	"math"

	"github.com/lindb/common/models"

	"github.com/lindb/lindb/pkg/bit"
	"github.com/lindb/lindb/pkg/collections"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
)

// for testing
var (
	toBytesFn = toBytes
)

// fieldIterator implements series.FieldIterator interface.
type fieldIterator struct {
	startSlot int
	aggTypes  []field.AggType

	fieldSeriesList []*collections.FloatArray

	length int
	idx    int
}

// newFieldIterator creates a field iterator.
func newFieldIterator(
	startSlot int,
	aggTypes []field.AggType,
	fieldSeriesList []*collections.FloatArray,
) series.FieldIterator {
	fmt.Println("**********************")
	fmt.Println(aggTypes)
	fmt.Println(fieldSeriesList)
	fmt.Println("**********************")
	return &fieldIterator{
		startSlot:       startSlot,
		aggTypes:        aggTypes,
		fieldSeriesList: fieldSeriesList,
		length:          len(fieldSeriesList),
	}
}

// HasNext returns if the iteration has more fields.
func (it *fieldIterator) HasNext() bool {
	return it.idx < it.length
}

// Next returns the data point in the iteration.
func (it *fieldIterator) Next() series.PrimitiveIterator {
	if it.idx >= it.length {
		return nil
	}
	primitiveIt := newPrimitiveIterator(it.startSlot, it.aggTypes[it.idx], it.fieldSeriesList[it.idx])
	it.idx++
	return primitiveIt
}

// MarshalBinary marshals the data.
func (it *fieldIterator) MarshalBinary() ([]byte, error) {
	if it.length == 0 {
		fmt.Println("no data......")
		return nil, nil
	}
	// need reset idx
	it.idx = 0
	writer := stream.NewBufferWriter(nil)
	var encoder *encoding.TSDEncoder
	defer encoding.ReleaseTSDEncoder(encoder)

	for it.HasNext() {
		primitiveIt := it.Next()
		if encoder == nil {
			encoder = encoding.TSDEncodeFunc(uint16(it.startSlot))
		} else {
			encoder.RestWithStartTime(uint16(it.startSlot))
		}
		idx := it.startSlot // start with start slot
		for primitiveIt.HasNext() {
			slot, value := primitiveIt.Next()
			for slot > idx {
				encoder.AppendTime(bit.Zero)
				idx++
			}
			encoder.AppendTime(bit.One)
			encoder.AppendValue(math.Float64bits(value))
			idx++
		}
		data, err := toBytesFn(encoder)
		if err != nil {
			return nil, err
		}
		writer.PutByte(byte(primitiveIt.AggType()))
		writer.PutVarint32(int32(len(data)))
		writer.PutBytes(data)
	}
	fmt.Println("write data..........")
	return writer.Bytes()
}

func (it *fieldIterator) Reset(data []byte) {}

func toBytes(e *encoding.TSDEncoder) ([]byte, error) {
	return e.Bytes()
}

// primitiveIterator represents primitive iterator using array.
type primitiveIterator struct {
	start   int
	aggType field.AggType
	it      *collections.FloatArrayIterator
}

// newPrimitiveIterator create primitive iterator using array.
func newPrimitiveIterator(start int, aggType field.AggType, values *collections.FloatArray) series.PrimitiveIterator {
	it := &primitiveIterator{
		start:   start,
		aggType: aggType,
	}
	if values != nil {
		it.it = values.NewIterator()
	}
	return it
}

// AggType returns the primitive field's agg type.
func (it *primitiveIterator) AggType() field.AggType {
	return it.aggType
}

// HasNext returns if the iteration has more data points.
func (it *primitiveIterator) HasNext() bool {
	if it.it == nil {
		return false
	}
	return it.it.HasNext()
}

// Next returns the data point in the iteration.
func (it *primitiveIterator) Next() (timeSlot int, value float64) {
	if it.it == nil {
		return -1, 0
	}
	timeSlot, value = it.it.Next()
	if timeSlot == -1 {
		return
	}
	timeSlot += it.start
	return
}

func (it *primitiveIterator) NextExemplar() (int, *models.Exemplar) {
	return -1, nil
}

func (it *primitiveIterator) Reset(aggType field.AggType, data []byte) {
}

type exemplarIterator struct {
	startSlot int

	exemplars map[int]*models.Exemplar

	length int
	idx    int
}

func newExemplarIterator(
	startSlot int,
	exemplars map[int]*models.Exemplar,
) series.FieldIterator {
	return &exemplarIterator{
		startSlot: startSlot,
		exemplars: exemplars,
		length:    1, // exemplar field only one agg type
	}
}

func (it *exemplarIterator) HasNext() bool {
	return it.idx < it.length
}

func (it *exemplarIterator) Next() series.PrimitiveIterator {
	if it.idx >= it.length {
		return nil
	}
	primitiveIt := newExemplarPrimitiveIterator(it.startSlot, it.exemplars)
	it.idx++
	return primitiveIt
}

func (it *exemplarIterator) MarshalBinary() ([]byte, error) {
	if it.length == 0 {
		return nil, nil
	}
	// need reset idx
	it.idx = 0
	writer := stream.NewBufferWriter(nil)
	writer.PutByte(byte(field.Exemplar))
	writer.PutVarint32(int32(len(it.exemplars)))
	for it.HasNext() {
		primitiveIt := it.Next()
		for primitiveIt.HasNext() {
			slot, exemplar := primitiveIt.NextExemplar()
			writer.PutVarint32(int32(slot))
			writer.PutVarint32(int32(len(exemplar.TraceID)))
			writer.PutBytes([]byte(exemplar.TraceID))
			writer.PutVarint32(int32(len(exemplar.SpanID)))
			writer.PutBytes([]byte(exemplar.SpanID))
			writer.PutVarint64(exemplar.Duration)
		}
	}
	return writer.Bytes()
}

func (it *exemplarIterator) Reset(data []byte) {}

type exemplarPrimitiveIterator struct {
	start int

	exemplars map[int]*models.Exemplar
	keys      []int
	vals      []*models.Exemplar
	index     int
}

func newExemplarPrimitiveIterator(start int, exemplars map[int]*models.Exemplar) series.PrimitiveIterator {
	keys := make([]int, 0, len(exemplars))
	vals := make([]*models.Exemplar, 0, len(exemplars))
	for k, v := range exemplars {
		keys = append(keys, k)
		vals = append(vals, v)
	}
	return &exemplarPrimitiveIterator{
		start:     start,
		exemplars: exemplars,
		keys:      keys,
		vals:      vals,
	}
}

func (it *exemplarPrimitiveIterator) AggType() field.AggType {
	return field.Exemplar
}

func (it *exemplarPrimitiveIterator) HasNext() bool {
	return it.index < len(it.keys)
}

func (it *exemplarPrimitiveIterator) Next() (timeSlot int, value float64) {
	return -1, 0
}

func (it *exemplarPrimitiveIterator) NextExemplar() (timeSlot int, exemplar *models.Exemplar) {
	exemplar = it.vals[it.index]
	timeSlot = it.keys[it.index]
	it.index++
	return
}

func (it *exemplarPrimitiveIterator) Reset(aggType field.AggType, data []byte) {
}
