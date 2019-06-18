package version

import (
	"fmt"
	"testing"

	"github.com/eleme/lindb/pkg/util"
	"github.com/stretchr/testify/assert"
)

var vsTestPath = "../../test_data/test_vs"

func Test_Recover(t *testing.T) {
	initTest()
	defer destory()

	var vs = NewVersionSet(vsTestPath)

	err := vs.Recover()

	assert.Nil(t, err, "recover edit log error")
}

func Test_Family(t *testing.T) {
	initTest()
	defer destory()

	var vs = NewVersionSet(vsTestPath)

	familyVersion := vs.CreateFamilyVersion("family")
	assert.NotNil(t, familyVersion, "get nil family version")

	familyVersion2 := vs.GetFamilyVersion("family")
	assert.NotNil(t, familyVersion2, "get nil family version2")

	assert.Equal(t, familyVersion, familyVersion2, "get diff family version")
}

func initTest() {
	if err := util.MkDirIfNotExist(vsTestPath); err != nil {
		fmt.Println("create test path error")
	}
}

func destory() {
	if err := util.RemoveDir(vsTestPath); err != nil {
		fmt.Println("delete test path error")
	}
}
