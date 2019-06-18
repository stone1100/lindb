package meta

import (
	"testing"
)

func Test_NewFile(t *testing.T) {
	newFile := &NewFile{level: 1, file: NewFileMeta(12, 1, 100, 2014)}
	bytes, err := newFile.Encode()
	if err != nil {
		t.Error(err)
		return
	}

	newFile2 := &NewFile{}
	err2 := newFile2.Decode(bytes)
	if err2 != nil {
		t.Error(err2)
		return
	}

	if newFile2.level != newFile.level ||
		newFile2.file.fileNumber != newFile.file.fileNumber ||
		newFile2.file.minKey != newFile.file.minKey ||
		newFile2.file.maxKey != newFile.file.maxKey ||
		newFile2.file.fileSize != newFile.file.fileSize {
		t.Error("file1 not equals file2")
		return
	}
}

func Test_DeleteFile(t *testing.T) {
	deleteFile := &DeleteFile{level: 1, fileNumber: 120}
	bytes, err := deleteFile.Encode()
	if err != nil {
		t.Error(err)
		return
	}

	deleteFile2 := &DeleteFile{}
	err2 := deleteFile2.Decode(bytes)
	if err2 != nil {
		t.Error(err2)
		return
	}

	if deleteFile2.level != deleteFile.level ||
		deleteFile2.fileNumber != deleteFile.fileNumber {
		t.Error("file1 not equals file2")
		return
	}
}
