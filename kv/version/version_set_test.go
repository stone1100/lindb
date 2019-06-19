package version

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eleme/lindb/pkg/util"
)

var vsTestPath = "../../test_data/test_vs"

func TestRecover(t *testing.T) {
	initVersionSetTestData()
	defer destoryVersionTestData()
	var vs = NewVersionSet(vsTestPath, 2)
	err := vs.Recover()

	assert.Nil(t, err, "recover edit log error")
}

func TestCreateFamily(t *testing.T) {
	initVersionSetTestData()
	defer destoryVersionTestData()

	var vs = NewVersionSet(vsTestPath, 2)

	familyVersion := vs.CreateFamilyVersion("family")
	assert.NotNil(t, familyVersion, "get nil family version")

	familyVersion2 := vs.GetFamilyVersion("family")
	assert.NotNil(t, familyVersion2, "get nil family version2")

	assert.Equal(t, familyVersion, familyVersion2, "get diff family version")
}

func initVersionSetTestData() {
	if err := util.MkDirIfNotExist(vsTestPath); err != nil {
		fmt.Println("create test path error")
	}
}

func destoryVersionTestData() {
	if err := util.RemoveDir(vsTestPath); err != nil {
		fmt.Println("delete test path error")
	}
}
