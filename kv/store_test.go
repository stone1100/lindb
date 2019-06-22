package kv

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eleme/lindb/pkg/util"
)

var testKVPath = "../test_data"

func TestReOpen(t *testing.T) {
	option := DefaultStoreOption(testKVPath)
	defer util.RemoveDir(testKVPath)

	var kv, _ = NewStore("test_kv", option)
	assert.NotNil(t, kv, "cannot create kv store")
	_, e := NewStore("test_kv", option)
	assert.NotNil(t, e, "store re-open not allow")

	f1, _ := kv.CreateFamily("f", FamilyOption{})
	assert.NotNil(t, f1, "cannot create family")
	assert.Equal(t, 1, kv.familyID, "store family id is wrong")

	kv.Close()

	kv2, e := NewStore("test_kv", option)
	if e != nil {
		t.Error(e)
	}
	assert.NotNil(t, kv2, "cannot re-open kv store")
	f2, _ := kv.GetFamily("f")
	assert.Equal(t, f1.name, f2.name, "family diff when store reopen")
	assert.Equal(t, f1.option.ID, f2.option.ID, "family id diff")
	assert.Equal(t, 1, kv2.familyID, "store family id is wrong")
}

func TestCreateFamily(t *testing.T) {
	option := DefaultStoreOption(testKVPath)
	defer util.RemoveDir(testKVPath)

	var kv, _ = NewStore("test_kv", option)
	defer kv.Close()
	assert.NotNil(t, kv, "cannot create kv store")
	assert.Equal(t, 0, kv.familyID, "wrong famiy id")

	f1, _ := kv.CreateFamily("f", FamilyOption{})
	assert.NotNil(t, f1, "cannot create family")
	assert.Equal(t, 1, f1.option.ID, "wrong famiy id")
	assert.Equal(t, 1, kv.familyID, "wrong famiy id")

	var f2, ok = kv.GetFamily("f")
	assert.True(t, ok, "can't get family")
	assert.Equal(t, f1, f2, "family not same for same name")

	_, ok = kv.GetFamily("f1")
	assert.False(t, ok, "get not exist family")
}
