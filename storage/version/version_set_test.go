package meta

import (
	"fmt"
	"testing"

	"github.com/eleme/lindb/pkg/util"
)

var vsTestPath = "../../test_data/test_vs"

func Test_Recover(t *testing.T) {
	initPath()

	var vs = NewVersionSet(vsTestPath)

	err := vs.Recover()
	if err != nil {
		t.Error(err)
		return
	}
}

func initPath() {
	if err := util.MkDirIfNotExist(vsTestPath); err != nil {
		fmt.Println("create test path error")
	}
}
