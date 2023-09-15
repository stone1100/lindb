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

package series

import (
	"fmt"
	"math"

	"github.com/lindb/common/models"

	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/series/field"
)

// binaryGroupedIterator implements GroupedIterator interface.
type binaryGroupedIterator struct {
	tags       string
	fields     map[field.Name][]byte
	fieldNames []field.Name

	it *BinaryIterator

	idx int
}

func NewGroupedIterator(tags string, fields map[field.Name][]byte) GroupedIterator {
	it := &binaryGroupedIterator{tags: tags, fields: fields}
	for fieldName := range fields {
		it.fieldNames = append(it.fieldNames, fieldName)
	}
	return it
}

func (g *binaryGroupedIterator) Tags() string {
	return g.tags
}

func (g *binaryGroupedIterator) HasNext() bool {
	if g.idx >= len(g.fieldNames) {
		return false
	}
	g.idx++
	return true
}

func (g *binaryGroupedIterator) Next() Iterator {
	fieldName := g.fieldNames[g.idx-1]
	if g.it == nil {
		g.it = NewIterator(fieldName, g.fields[fieldName])
	} else {
		g.it.Reset(fieldName, g.fields[fieldName])
	}
	return g.it
}

// BinaryIterator implements Iterator interface.
type BinaryIterator struct {
	fieldName field.Name
	fieldType field.Type
	reader    *stream.Reader
	fieldIt   FieldIterator
	data      []byte
}

func NewIterator(fieldName field.Name, data []byte) *BinaryIterator {
	it := &BinaryIterator{fieldName: fieldName, reader: stream.NewReader(data), data: data}
	it.fieldType = field.Type(it.reader.ReadByte())
	return it
}

func (b *BinaryIterator) Reset(fieldName field.Name, data []byte) {
	b.fieldName = fieldName
	b.reader.Reset(data)
	b.fieldType = field.Type(b.reader.ReadByte())
}

func (b *BinaryIterator) FieldName() field.Name {
	return b.fieldName
}

func (b *BinaryIterator) FieldType() field.Type {
	return b.fieldType
}

func (b *BinaryIterator) HasNext() bool {
	return !b.reader.Empty()
}

func (b *BinaryIterator) Next() (startTime int64, fieldIt FieldIterator) {
	startTime = b.reader.ReadVarint64()
	length := b.reader.ReadVarint32()
	if length == 0 {
		return
	}
	data := b.reader.ReadBytes(int(length))
	if b.fieldIt == nil {
		if b.fieldType == field.ExemplarField {
			b.fieldIt = NewExemplarterator(data)
		} else {
			b.fieldIt = NewFieldIterator(data)
		}
	} else {
		b.fieldIt.Reset(data)
	}
	fieldIt = b.fieldIt
	return
}

func (b *BinaryIterator) MarshalBinary() ([]byte, error) {
	return b.data, nil
}

// BinaryFieldIterator implements FieldIterator interface.
type BinaryFieldIterator struct {
	reader *stream.Reader
	pIt    PrimitiveIterator
}

// NewFieldIterator create field iterator based on binary data
func NewFieldIterator(data []byte) FieldIterator {
	it := &BinaryFieldIterator{
		reader: stream.NewReader(data),
	}
	return it
}

func (it *BinaryFieldIterator) Reset(data []byte) {
	it.reader.Reset(data)
}

func (it *BinaryFieldIterator) HasNext() bool { return !it.reader.Empty() }

func (it *BinaryFieldIterator) Next() PrimitiveIterator {
	aggType := field.AggType(it.reader.ReadByte())
	length := it.reader.ReadVarint32()
	data := it.reader.ReadBytes(int(length))

	if it.pIt == nil {
		it.pIt = NewPrimitiveIterator(aggType, encoding.NewTSDDecoder(data)) // TODO get from pool?
	} else {
		it.pIt.Reset(aggType, data)
	}
	return it.pIt
}

func (it *BinaryFieldIterator) MarshalBinary() ([]byte, error) {
	return nil, fmt.Errorf("not support")
}

// BinaryPrimitiveIterator implements PrimitiveIterator interface.
type BinaryPrimitiveIterator struct {
	aggType field.AggType
	tsd     *encoding.TSDDecoder
}

func NewPrimitiveIterator(aggType field.AggType, tsd *encoding.TSDDecoder) PrimitiveIterator {
	return &BinaryPrimitiveIterator{
		aggType: aggType,
		tsd:     tsd,
	}
}

func (pi *BinaryPrimitiveIterator) Reset(aggType field.AggType, data []byte) {
	pi.aggType = aggType
	pi.tsd.Reset(data)
}

func (pi *BinaryPrimitiveIterator) AggType() field.AggType {
	return pi.aggType
}

func (pi *BinaryPrimitiveIterator) HasNext() bool {
	if pi.tsd.Error() != nil {
		return false
	}
	for pi.tsd.Next() {
		if pi.tsd.HasValue() {
			return true
		}
	}
	return false
}

func (pi *BinaryPrimitiveIterator) Next() (timeSlot int, value float64) {
	timeSlot = int(pi.tsd.Slot())
	val := pi.tsd.Value()
	value = math.Float64frombits(val)
	return
}

func (pi *BinaryPrimitiveIterator) NextExemplar() (timeSlot int, exemplar *models.Exemplar) {
	return
}

type BinaryExemplarIterator struct {
	reader *stream.Reader
	pIt    PrimitiveIterator
}

func NewExemplarterator(data []byte) FieldIterator {
	it := &BinaryExemplarIterator{
		reader: stream.NewReader(data),
	}
	return it
}

func (it *BinaryExemplarIterator) Reset(data []byte) {
	it.reader.Reset(data)
}

func (it *BinaryExemplarIterator) HasNext() bool { return !it.reader.Empty() }

func (it *BinaryExemplarIterator) Next() PrimitiveIterator {
	aggType := field.AggType(it.reader.ReadByte())
	// only one agg result
	it.pIt = NewExemplarPrimitiveIterator(aggType, it.reader)
	return it.pIt
}

func (it *BinaryExemplarIterator) MarshalBinary() ([]byte, error) {
	return nil, fmt.Errorf("not support")
}

type BinaryExemplarPrimitiveIterator struct {
	aggType field.AggType
	reader  *stream.Reader

	count int
	index int
}

func NewExemplarPrimitiveIterator(aggType field.AggType, reader *stream.Reader) PrimitiveIterator {
	return &BinaryExemplarPrimitiveIterator{
		aggType: aggType,
		reader:  reader,
		count:   int(reader.ReadVarint32()),
	}
}

func (pi *BinaryExemplarPrimitiveIterator) Reset(aggType field.AggType, data []byte) {}

func (pi *BinaryExemplarPrimitiveIterator) AggType() field.AggType {
	return pi.aggType
}

func (pi *BinaryExemplarPrimitiveIterator) HasNext() bool {
	ok := pi.index < pi.count
	pi.index++
	return ok
}

func (pi *BinaryExemplarPrimitiveIterator) Next() (timeSlot int, value float64) {
	return
}

func (pi *BinaryExemplarPrimitiveIterator) NextExemplar() (timeSlot int, exemplar *models.Exemplar) {
	fmt.Println("next binary exemplar")
	timeSlot = int(pi.reader.ReadVarint32())
	size := pi.reader.ReadVarint32()
	exemplar = &models.Exemplar{}
	exemplar.TraceID = string(pi.reader.ReadBytes(int(size)))
	size = pi.reader.ReadVarint32()
	exemplar.SpanID = string(pi.reader.ReadBytes(int(size)))
	exemplar.Duration = pi.reader.ReadVarint64()
	return
}
