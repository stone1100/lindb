package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewFile(t *testing.T) {
	newFile := &NewFile{level: 1, file: NewFileMeta(12, 1, 100, 2014)}
	bytes, err := newFile.Encode()
	assert.Nil(t, err, "new file encode error")

	newFile2 := &NewFile{}
	err2 := newFile2.Decode(bytes)
	assert.Nil(t, err2, "new file decode error")

	assert.Equal(t, newFile, newFile2, "file1 not eqals files")
}

func Test_DeleteFile(t *testing.T) {
	deleteFile := &DeleteFile{level: 1, fileNumber: 120}
	bytes, err := deleteFile.Encode()
	assert.Nil(t, err, "delete file encode error")

	deleteFile2 := &DeleteFile{}
	err2 := deleteFile2.Decode(bytes)

	assert.Nil(t, err2, "delete file decode error")
	assert.Equal(t, deleteFile, deleteFile2, "delete file1 not equals delete file2")
}
