package meta

import "testing"

func Test_EditLog(t *testing.T) {
	editLog := NewEditLog()
	newFile := &NewFile{level: 1, file: NewFileMeta(12, 1, 100, 2014)}
	editLog.Add(newFile)
	editLog.Add(NewDeleteFile(1, 123))

	v, err := editLog.bytes()
	if err != nil {
		t.Error(err)
		return
	}
	if len(v) <= 0 {
		t.Errorf("encode edit log error")
		return
	}
}
