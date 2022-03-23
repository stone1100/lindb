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

package query

import (
	"encoding/binary"
	"fmt"
	"github.com/lindb/lindb/pkg/timeutil"
	"strconv"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/series"
)

func TestGroupingContext_Build(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()
	scanner := series.NewMockGroupingScanner(ctrl)
	ctx := NewGroupContext([]uint32{1}, map[uint32][]series.GroupingScanner{1: {scanner}})
	scanner.EXPECT().GetSeriesAndTagValue(uint16(1)).
		Return(roaring.BitmapOf(1, 2, 3, 10).GetContainerAtIndex(0), []uint32{10, 20, 30, 10})
	result := ctx.BuildGroup(1, roaring.BitmapOf(1, 2, 6, 10).GetContainerAtIndex(0))
	assert.Len(t, result, 2)
	tagValueIDs := make([]byte, 4)
	binary.LittleEndian.PutUint32(tagValueIDs[0:], 10)
	seriesIDs := result[string(tagValueIDs)]
	assert.Equal(t, []uint16{1, 10}, seriesIDs)
	binary.LittleEndian.PutUint32(tagValueIDs[0:], 20)
	seriesIDs = result[string(tagValueIDs)]
	assert.Equal(t, []uint16{2}, seriesIDs)

	scanner.EXPECT().GetSeriesAndTagValue(uint16(2)).
		Return(roaring.BitmapOf(1, 2).GetContainerAtIndex(0), []uint32{30, 10})
	_ = ctx.BuildGroup(2, roaring.BitmapOf(1, 2).GetContainerAtIndex(0))
	// container not found
	scanner.EXPECT().GetSeriesAndTagValue(uint16(3)).Return(nil, nil)
	_ = ctx.BuildGroup(3, roaring.BitmapOf(1, 2).GetContainerAtIndex(0))
	// case: get group by tag value ids
	groupByTagValueIDs := ctx.GetGroupByTagValueIDs()
	assert.Len(t, groupByTagValueIDs, 1)
	assert.EqualValues(t, roaring.BitmapOf(10, 20, 30).ToArray(), groupByTagValueIDs[0].ToArray())
}

func TestGroupingContext_ScanTagValueIDs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()
	scanner := series.NewMockGroupingScanner(ctrl)
	ctx := NewGroupContext([]uint32{1}, map[uint32][]series.GroupingScanner{1: {scanner}})
	// case 1: get tag value ids
	scanner.EXPECT().GetSeriesAndTagValue(uint16(1)).
		Return(roaring.BitmapOf(1, 2, 3, 10).GetContainerAtIndex(0), []uint32{10, 20, 30, 10})
	result := ctx.ScanTagValueIDs(1, roaring.BitmapOf(1, 2, 6, 10).GetContainerAtIndex(0))
	assert.Equal(t, []uint32{10, 20}, result[0].ToArray())
	// case 2: empty tag value
	scanner.EXPECT().GetSeriesAndTagValue(uint16(1)).Return(nil, nil)
	result = ctx.ScanTagValueIDs(1, roaring.BitmapOf(1, 2, 6, 10).GetContainerAtIndex(0))
	assert.Equal(t, roaring.New(), result[0])
}

func TestGroupingContext_Build2(t *testing.T) {
	hosts := NewTagValuesEntrySet()
	disks := NewTagValuesEntrySet()
	partitions := NewTagValuesEntrySet()
	total := roaring.New()
	id := uint32(0)
	count := 40000
	for i := 0; i < count; i++ {
		for j := 0; j < 4; j++ {
			for k := 0; k < 20; k++ {
				total.Add(id)
				id++
				host := "host" + strconv.Itoa(i)
				disk := "/tmp" + strconv.Itoa(j)
				partition := "partition" + strconv.Itoa(k)
				h, ok := hosts.values[host]
				if !ok {
					hosts.values[host] = roaring.BitmapOf(id)
				} else {
					h.Add(id)
				}

				d, ok := disks.values[disk]
				if !ok {
					disks.values[disk] = roaring.BitmapOf(id)
				} else {
					d.Add(id)
				}

				p, ok := partitions.values[partition]
				if !ok {
					partitions.values[partition] = roaring.BitmapOf(id)
				} else {
					p.Add(id)
				}
			}
		}
	}

	// test single group by tag keys
	ctx := NewGroupContext2(1)
	ctx.SetTagValuesEntrySet(0, disks)
	assert.Equal(t, 1, ctx.Len())
	total= roaring.New()
	total.AddRange(0, uint64(1000000))
	//keys := seriesIDs.GetHighKeys()
	keys := total.GetHighKeys()
	i := 0
	s := timeutil.Now()
	for idx, key := range keys {
		container := total.GetContainerAtIndex(idx)
		i += container.GetCardinality()
		k := key
		_= ctx.BuildGroup(k, container)
		//assert.Len(t, rs, 4)
	}
	fmt.Println(timeutil.Now() - s)
	// test single group by tag keys
	s = timeutil.Now()
	ctx = NewGroupContext2(2)
	ctx.SetTagValuesEntrySet(0, disks)
	ctx.SetTagValuesEntrySet(1, partitions)
	assert.Equal(t, 2, ctx.Len())
	var wait sync.WaitGroup
	for idx, key := range keys {
		container := total.GetContainerAtIndex(idx)
		i += container.GetCardinality()
		k := key
		wait.Add(1)
		go func() {
			_ = ctx.BuildGroup(k, container)
			//assert.Len(t, rs, 800000)
			wait.Done()
		}()
	}
	wait.Wait()
	fmt.Println(timeutil.Now() - s)
	var data [][]byte
	for _, v := range hosts.values {
		d, _ := v.MarshalBinary()
		data = append(data, d)
	}
	fmt.Println(len(data))
	s = timeutil.Now()
	for _, v := range data {
		bitmap := roaring.New()
		_, _ = bitmap.FromBuffer(v)
	}
	fmt.Println(timeutil.Now() - s)
}