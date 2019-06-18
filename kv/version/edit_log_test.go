package version 

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_EditLog_Codec(t *testing.T) {
	editLog := NewEditLog()
	newFile := &NewFile{level: 1, file: NewFileMeta(12, 1, 100, 2014)}
	editLog.Add(newFile)
	editLog.Add(NewDeleteFile(1, 123))

	v, err := editLog.marshal()

	assert.Nil(t, err, "marshal error")
	assert.True(t, len(v) > 0, "encode edit log error")

	editLog2 := NewEditLog()
	err2 := editLog2.unmarshal(v)
	assert.Nil(t, err2, "unmarshal error")

	assert.Equal(t, editLog, editLog2, "edit log not eqauls")
}

func Test_Apply(t *testing.T) {
	editLog := NewEditLog()
	newFile := &NewFile{level: 1, file: NewFileMeta(12, 1, 100, 2014)}
	editLog.Add(newFile)
	version := newVersion(1, nil)
	editLog.apply(version)

	assert.Equal(t, 1, len(version.GetAllFiles()), "cannot add file into version")

	//delete file
	editLog2 := NewEditLog()
	editLog2.Add(NewDeleteFile(1, 12))
	editLog2.apply(version)
	assert.Equal(t, 0, len(version.GetAllFiles()), "cannot delete file from version")
}
