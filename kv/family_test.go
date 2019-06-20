package kv

import (
	"testing"

	"github.com/eleme/lindb/kv/version"
	"github.com/eleme/lindb/pkg/util"

	"github.com/stretchr/testify/assert"
)

func TestNewTableBuilder(t *testing.T) {
	option := DefaultStoreOption(testKVPath)
	defer util.RemoveDir(testKVPath)

	var kv, err = NewStore("test_kv", option)
	defer kv.Close()
	assert.Nil(t, err, "cannot create kv store")

	f, err := kv.CreateFamily("f", FamilyOption{})
	assert.Nil(t, err, "cannot create family")

	f.NewTableBuilder()
}

func TestCommitEditLog(t *testing.T) {
	option := DefaultStoreOption(testKVPath)
	defer util.RemoveDir(testKVPath)

	var kv, _ = NewStore("test_kv", option)
	defer kv.Close()

	f, _ := kv.CreateFamily("f", FamilyOption{})

	editLog := version.NewEditLog(1)
	newFile := version.CreateNewFile(1, version.NewFileMeta(12, 1, 100, 2014))
	editLog.Add(newFile)
	editLog.Add(version.NewDeleteFile(1, 123))

	flag := f.CommitEditLog(editLog)
	assert.True(t, flag, "commit edit log error")
}
