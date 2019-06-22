package kv

import (
	"fmt"

	"github.com/eleme/lindb/kv/table"
	"github.com/eleme/lindb/kv/version"
)

// Flusher flushs data into kv store, for big data maybe split many sstable
type Flusher interface {
	// Add puts k/v pair
	Add(key uint32, value []byte) error
	// Commit flushs data and commits metadata
	Commit() error
}

// storeFlusher family level store flusher
type storeFlusher struct {
	family  *Family
	builder table.Builder
	editLog *version.EditLog
}

func newStoreFlusher(family *Family) Flusher {
	return &storeFlusher{
		family:  family,
		editLog: version.NewEditLog(family.option.ID),
	}
}

// Add puts k/v pair
func (sf *storeFlusher) Add(key uint32, value []byte) error {
	if sf.builder == nil {
		builder, err := sf.family.newTableBuilder()
		if err != nil {
			return fmt.Errorf("create table build error:%s", err)
		}
		sf.builder = builder
	}
	//TODO add file size limit
	return sf.builder.Add(key, value)
}

// Commit flushs data and commits metadata
func (sf *storeFlusher) Commit() error {
	builder := sf.builder
	if builder != nil {
		if err := builder.Close(); err != nil {
			return fmt.Errorf("close table builder error when flush commit, error:%s", err)
		}

		fileMeta := version.NewFileMeta(builder.FileNumber(), builder.MinKey(), builder.MaxKey(), builder.Size())
		sf.editLog.Add(version.CreateNewFile(0, fileMeta))
	}

	if flag := sf.family.commitEditLog(sf.editLog); !flag {
		return fmt.Errorf("commit edit log failure")
	}
	return nil
}
