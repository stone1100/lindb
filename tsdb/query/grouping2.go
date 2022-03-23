package query

import (
	"github.com/lindb/roaring"
	"strings"
)

// SeriesID2Tags represents the tag values of series id
type SeriesID2Tags struct {
	tagValues strings.Builder
	c         int
}

// NewSeriesID2Tags creates a SeriesID2Tags
func NewSeriesID2Tags() SeriesID2Tags {
	return SeriesID2Tags{c: 1}
}

// TagValues returns the tag values
//func (entry *SeriesID2Tags) TagValues() []*string {
//	return entry.tagValues
//}

func (entry *SeriesID2Tags) Inc() {
	entry.c++
}

// AddTagValue adds the tag value for series id
func (entry *SeriesID2Tags) AddTagValue(tagValue string) {
	entry.tagValues.WriteString(tagValue)
}

// TagValuesEntrySet represents the tag values and series ids mapping for a tag key
type TagValuesEntrySet struct {
	values map[string]*roaring.Bitmap
}

// NewTagValuesEntrySet creates a TagValuesEntrySet
func NewTagValuesEntrySet() *TagValuesEntrySet {
	return &TagValuesEntrySet{values: make(map[string]*roaring.Bitmap)}
}

// Values returns the tag values data
func (tes *TagValuesEntrySet) Values() map[string]*roaring.Bitmap {
	return tes.values
}

// SetTagValues sets the tag values data
func (tes *TagValuesEntrySet) SetTagValues(values map[string]*roaring.Bitmap) {
	tes.values = values
}

// AddTagValue adds tag value and series ids
func (tes *TagValuesEntrySet) AddTagValue(tagValue string, seriesIDs *roaring.Bitmap) {
	oldSeriesIDs, ok := tes.values[tagValue]
	if ok {
		oldSeriesIDs.Or(seriesIDs)
	} else {
		tes.values[tagValue] = seriesIDs
	}
}

// GroupingContext represents the context of group by query for tag keys
type GroupingContext2 struct {
	tagValueIDs []*roaring.Bitmap
	seriesIDs   [][]*roaring.Bitmap
}

func (g *GroupingContext2) GetGroupByTagValueIDs() []*roaring.Bitmap {
	panic("implement me")
}

func (g *GroupingContext2) ScanTagValueIDs(highKey uint16, container roaring.Container) []*roaring.Bitmap {
	panic("implement me")
}

// NewGroupContext creates a GroupingContext
func NewGroupContext2(tagValueIDs []*roaring.Bitmap, seriesIDs [][]*roaring.Bitmap) *GroupingContext2 {
	return &GroupingContext2{
		tagValueIDs: tagValueIDs,
		seriesIDs:   seriesIDs,
	}
}

//// SetTagValuesEntrySet sets the tag values entry set for group by tag keys
//func (g *GroupingContext2) SetTagValuesEntrySet(idx int, tagValuesEntrySet *TagValuesEntrySet) {
//	g.tagValuesEntrySets[idx] = tagValuesEntrySet
//}

//// Len returns the group by tag key's length
//func (g *GroupingContext2) Len() int {
//	return len(g.tagValuesEntrySets)
//}

// BuildGroup builds the grouped series ids by the high key of series id
// and the container includes low keys of series id
func (g *GroupingContext2) BuildGroup(highKey uint16, container roaring.Container) map[string][]uint16 {
	groupTagKeysCount := len(g.tagValueIDs)
	if groupTagKeysCount == 1 {
		return g.buildForSingleTagKey(highKey, container)
	}

	// new seriesIDs2Tags array based on range of max ~ min
	//seriesIDs2Tags := g.buildSeriesIDs2Tags(highKey, container)
	_ = g.buildSeriesIDs2Tags(highKey, container)

	// finds group tags => series IDs, and builds result
	//it := container.PeekableIterator()
	//min := container.Minimum()
	result := make(map[string][]uint16)
	//for it.HasNext() {
	//	lowKey := it.Next()
	//	idx := lowKey - min
	//	seriesID2Tags := seriesIDs2Tags[idx]
	//	if seriesID2Tags.c-1 == groupTagKeysCount {
	//		tagValuesStr := seriesID2Tags.tagValues.String()
	//		//tag.ConcatTagValues(seriesID2Tags.TagValues())
	//		values, ok := result[tagValuesStr]
	//		if !ok {
	//			result[tagValuesStr] = []uint16{lowKey}
	//		} else {
	//			result[tagValuesStr] = append(values, lowKey)
	//		}
	//	}
	//}
	return result
}

// buildSeriesIDs2Tags builds for multi group by keys
func (g *GroupingContext2) buildSeriesIDs2Tags(highKey uint16, lowSeriesIDs roaring.Container) []SeriesID2Tags {
	//groupTagKeysCount := len(g.tagValuesEntrySets)
	// new seriesIDs2Tags array based on range of min ~ max
	min := lowSeriesIDs.Minimum()
	max := lowSeriesIDs.Maximum()
	seriesIDs2Tags := make([]SeriesID2Tags, int(max-min)+1)
	//sb := make([]strings.Builder, int(max-min)+1)
	c := make([]int, int(max-min)+1)
	for i := min; i < max; i++ {
		seriesIDs2Tags[i] = NewSeriesID2Tags()
		c[i] += 1
	}

	// builds seriesIDs => tags mapping, using counting sort
	// https://en.wikipedia.org/wiki/Counting_sort
	for tagIdx, tagKV := range g.tagValueIDs {
		seriesIDs := g.seriesIDs[tagIdx]
		for tagValue, lowSeriesIDs := range tagKV.Values() {
			//tagValue := &v
			lowContainer := lowSeriesIDs.GetContainer(highKey)
			if lowContainer != nil {
				it := lowContainer.PeekableIterator()
				for it.HasNext() {
					v := it.Next()
					if v < min {
						continue
					}
					if v > max {
						break
					}
					idx := v - min // index = lowKey - min
					//if seriesIDs2Tags[idx] != nil {
					seriesIDs2Tags[idx].AddTagValue(tagValue)
					seriesIDs2Tags[idx].Inc()
					//sb[idx].WriteString(tagValue)
					//c[idx]+=1
					//}
				}
			}
		}
	}
	return seriesIDs2Tags
}

// buildForSingleTagKey builds for single group by tag key
func (g *GroupingContext2) buildForSingleTagKey(highKey uint16, container roaring.Container) map[string][]uint16 {
	result := make(map[string][]uint16)
	tagValueIDs := g.tagValueIDs[0]
	seriesIDs := g.seriesIDs[0]

	for idx, tagKV := range g.tagValueIDs {
		lowContainer := lowSeriesIDs.GetContainer(highKey)
		if lowContainer != nil {
			matchContainer := lowContainer.And(container)
			result[tagKV] = matchContainer.ToArray()
		}
	}
	return result
}
