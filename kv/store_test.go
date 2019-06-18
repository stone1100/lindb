package kv 

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Create_Family(t *testing.T) {
	option := StoreOption{Path: "../test_data"}
	var kv, err = NewStore("test_kv", option)
	defer kv.Close()
	assert.Nil(t, err, "cannot create kv store")

	f1, err2 := kv.CreateFamily("f", FamilyOption{})
	assert.Nil(t, err2, "cannot create family")

	var f2, ok = kv.GetFamily("f")
	assert.True(t, ok, "can't get family")
	assert.Equal(t, f1, f2, "family not same for same name")

	_, ok = kv.GetFamily("f1")
	assert.False(t, ok, "get not exist family")

	_, e := NewStore("test_kv", option)
	assert.NotNil(t, e, "store re-open not allow")
}
