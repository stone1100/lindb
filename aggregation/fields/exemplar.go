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

package fields

import (
	"fmt"

	"github.com/lindb/common/models"

	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
)

//go:generate mockgen -source=./field.go -destination=./field_mock.go -package=fields

type Exemplar interface {
	// Type returns field's type
	Type() field.Type
	// SetValue sets field value using series iterator.
	SetValue(fieldSeries series.Iterator)
	// GetValues returns the values which function call need by given function type.
	GetExemplars() map[int64][]*models.Exemplar
	// Reset resets field's value for reusing.
	Reset()
}

type exemplarField struct {
	fieldType field.Type
	startTime int64
	interval  int64
	capacity  int

	exemplars map[int64][]*models.Exemplar
}

func NewExemplarField(fieldType field.Type, startTime, interval int64, capacity int) Exemplar {
	return &exemplarField{
		fieldType: fieldType,
		startTime: startTime,
		interval:  interval,
		capacity:  capacity,
		exemplars: make(map[int64][]*models.Exemplar),
	}
}

// Type returns the type of dynamic field.
func (f *exemplarField) Type() field.Type {
	return f.fieldType
}

// SetValue sets the field's value by time slot
func (f *exemplarField) SetValue(fieldSeries series.Iterator) {
	if fieldSeries == nil {
		fmt.Println("nil exemplar....")
		return
	}
	for fieldSeries.HasNext() {
		startTime, it := fieldSeries.Next()
		if it == nil {
			fmt.Println("nil exemplar....1")
			continue
		}
		for it.HasNext() {
			pIt := it.Next()
			fmt.Println("nil exemplar....2")
			for pIt.HasNext() {
				slot, exemplar := pIt.NextExemplar()
				ts := (int64(slot)*f.interval + startTime)
				// FIXME:
				f.exemplars[ts] = []*models.Exemplar{exemplar}
				fmt.Println("nil exemplar....3")
			}
		}
	}
}

// GetValues returns the values which function call need by given function type and field type
func (f *exemplarField) GetExemplars() map[int64][]*models.Exemplar {
	return f.exemplars
}

func (f *exemplarField) Reset() {
}
