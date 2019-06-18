package meta

import (
	"fmt"
)

const sstSuffix = "sst"
const tmpSuffix = "tmp"

//const options = "OPTIONS"
const Lock = "LOCK"
const manifestPrefix = "MANIFEST-"

// FileType represent a file type.
type FileType int

// File types.
const (
	TypeManifest FileType = 1 << iota
	TypeJournal
	TypeTable
	TypeTemp
	TypeInfo

	TypeAll = TypeManifest | TypeJournal | TypeTable | TypeTemp | TypeInfo
)

// FileDesc define file type and file number
type FileDesc struct {
	FileType   string
	FileNumber int64
}

// current return current file name for saving manifest file name
func current() string {
	return "CURRENT"
}

// Table file name
func Table(fileNumber int64) string {
	return fmt.Sprintf("%06d.%s", fileNumber, sstSuffix)
}

// Info file for store/family option
func Info() string {
	return "info"
}

func manifestFileName(fileNumber int64) string {
	return fmt.Sprintf("%s%06d", manifestPrefix, fileNumber)
}
