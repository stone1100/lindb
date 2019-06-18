package meta

import "testing"

func Test_File_Level(t *testing.T) {
	level := newLevel()

	level.addFile(*NewFileMeta(1, 1, 10, 1024))
	level.addFile(*NewFileMeta(1, 1, 10, 1024))
	level.addFile(*NewFileMeta(1, 1, 10, 1024))

	var files = level.getFiles()

	if len(files) != 1 {
		t.Errorf("add file wrong")
		return
	}

	//add file
	level.addFile(*NewFileMeta(2, 1, 10, 1024))
	level.addFile(*NewFileMeta(20, 1, 10, 1024))

	//delete file
	level.deleteFile(2)

	files = level.getFiles()
	if len(files) != 2 {
		t.Errorf("delete file wrong")
		return
	}
}

func Test_Add_Files(t *testing.T) {
	level := newLevel()

	level.addFiles(*NewFileMeta(1, 1, 10, 1024), *NewFileMeta(2, 1, 10, 1024), *NewFileMeta(3, 1, 10, 1024))

	var files = level.getFiles()

	if len(files) != 3 {
		t.Errorf("add files wrong")
		return
	}
}
