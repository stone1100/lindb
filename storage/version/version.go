package meta

import (
	"sync/atomic"
)

// Version is snapshot for current storage metadata includes levels/sst files
type Version struct {
	id int64 // unique id in kv store level
	fv *FamilyVersion

	ref int32 // current version ref count for using

	levels []level // each level sst files exclude level0
}

// newVersion new Version instance
func newVersion(id int64, fv *FamilyVersion) *Version {
	return &Version{
		id: id,
		fv: fv,
	}
}

// Release decrements version ref count,
// if ref==0, then remove current version from list of family level.
func (v *Version) Release() {
	val := atomic.AddInt32(&v.ref, -1)
	if val == 0 {
		v.fv.removeVersion(v)
	}
}

// retain increments version ref count
func (v *Version) retain() {
	atomic.AddInt32(&v.ref, 1)
}

// cloneVersion builds new version based on current version
func (v *Version) cloneVersion() *Version {
	newVersion := newVersion(v.fv.versionSet.versionID, v.fv)
	for level, value := range v.levels {
		for _, file := range value.files {
			v.addFile(level, file)
		}
	}
	return newVersion
}

// addFiles adds file meta into spec level
func (v *Version) addFiles(level int, files []FileMeta) {
	v.levels[level].addFiles(files...)
}

// addFile adds file meta into spec level
func (v *Version) addFile(level int, file FileMeta) {
	v.levels[level].addFile(file)
}

// deleteFile delete file from spec level file list
func (v *Version) deleteFile(level int, fileNumber int64) {
	v.levels[level].deleteFile(fileNumber)
}
