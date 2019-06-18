package kv

import (
	"testing"

	"github.com/eleme/lindb/kv/version"
	"github.com/stretchr/testify/assert"
)

func Test_NewTableBuilder(t *testing.T) {
	option := StoreOption{Path: "../test_data"}
	var kv, err = NewStore("test_kv", option)
	defer kv.Close()
	assert.Nil(t, err, "cannot create kv store")

	f, err := kv.CreateFamily("f", FamilyOption{})
	assert.Nil(t, err, "cannot create family")

	f.NewTableBuilder()
}

func Test_Commit_EditLog(t *testing.T) {
	option := StoreOption{Path: "../test_data"}
	var kv, _ = NewStore("test_kv", option)
	defer kv.Close()

	editLog := version.NewEditLog()
	newFile := version.CreateNewFile(1, version.NewFileMeta(12, 1, 100, 2014))
	editLog.Add(newFile)
	editLog.Add(version.NewDeleteFile(1, 123))

}
