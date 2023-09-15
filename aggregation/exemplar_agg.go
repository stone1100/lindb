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

	"github.com/lindb/common/models"

	"github.com/lindb/lindb/series"
)

type exemplarAggregator struct {
	segmentStartTime int64
	start, end       int // slot range based on query interval and time range

	exemplars map[int]*models.Exemplar
}

// NewFieldAggregator creates a field aggregator,
// time range 's start and end is index based on segment start time and interval.
// e.g. segment start time = 20190905 10:00:00, start = 10, end = 50, interval = 10 seconds,
// real query time range {20190905 10:01:40 ~ 20190905 10:08:20}
func NewExemplarAggregator(_ AggregatorSpec, segmentStartTime int64, start, end int) FieldAggregator {
	agg := &exemplarAggregator{
		segmentStartTime: segmentStartTime,
		start:            start,
		end:              end,
		exemplars:        make(map[int]*models.Exemplar),
	}
	return agg
}

// ResultSet returns the result set of field aggregator
func (a *exemplarAggregator) ResultSet() (startTime int64, it series.FieldIterator) {
	return a.segmentStartTime, newExemplarIterator(a.start, a.exemplars)
}

// Aggregate aggregates the field series into current aggregator
func (a *exemplarAggregator) Aggregate(it series.FieldIterator) {
	for it.HasNext() {
		pIt := it.Next()
		for pIt.HasNext() {
			a.AggregateExemplarBySlot(pIt.NextExemplar())
		}
	}
}

// AggregateBySlot aggregates the field series into current aggregator
func (a *exemplarAggregator) AggregateBySlot(slot int, value float64) {
	// do nothing
}

func (a *exemplarAggregator) AggregateExemplarBySlot(slot int, exemplar *models.Exemplar) {
	fmt.Println("agg exemplar slot")
	a.exemplars[slot] = exemplar
}

// reset aggregator context for reusing.
func (a *exemplarAggregator) reset() {
	a.exemplars = nil
}
