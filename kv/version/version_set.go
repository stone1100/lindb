package version

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/eleme/lindb/pkg/logger"
	"github.com/eleme/lindb/pkg/util"
	"github.com/eleme/lindb/kv/journal"
	"go.uber.org/zap"
)

// VersionSet maintains all metadata for kv store
type VersionSet struct {
	manifestFileNumber int64
	nextFileNumber     int64
	storePath          string
	familyVersions     map[string]*FamilyVersion
	versionID          int64 // unique in for increasing version id

	manifest *journal.Writer
	mutex    sync.RWMutex

	logger *zap.Logger
}

// NewVersionSet new VersionSet instance
func NewVersionSet(storePath string) *VersionSet {
	vs := &VersionSet{
		manifestFileNumber: 1, // default value for initialize store
		nextFileNumber:     2, // default value
		storePath:          storePath,
		familyVersions:     make(map[string]*FamilyVersion),
		logger:             logger.GetLogger(),
	}
	return vs
}

// Recover metadata from manifest file
func (vs *VersionSet) Recover() error {
	new, err := vs.initializeIfNeeded()
	if err != nil {
		return err
	}
	if !new {
		vs.logger.Info("recover version set data from journal file", zap.String("store", vs.storePath))
		//TODO do recover log
		// do recover logic, read journal wal record and recover it

	}
	return nil
}

// Destory closes version set, release resource, such as journal writer etc.
func (vs *VersionSet) Destory() {
	vs.mutex.Unlock()
	defer vs.mutex.Unlock()

	// close manifest journal writer if it exist
	if vs.manifest != nil {
		vs.manifest.Close()
	}
}

// NextFileNumber generates next file number
func (vs *VersionSet) NextFileNumber() int64 {
	vs.mutex.Lock()
	nf := vs.nextFileNumber
	vs.nextFileNumber++
	vs.mutex.Unlock()
	return nf
}

// Commit peresists edit logs to manifest file, then apply new version to family version
func (vs *VersionSet) Commit(family string, editLog *EditLog) error {
	// get family version based on family name
	familyVersion := vs.GetFamilyVersion(family)
	if familyVersion != nil {
		return fmt.Errorf("cannot find family version for name: %s", family)
	}

	vs.mutex.Lock()
	defer vs.mutex.Unlock()

	// add next file number init edit log for each delta edit log
	editLog.Add(NewNextFileNumber(vs.nextFileNumber))

	v, err := editLog.marshal()
	if err != nil {
		return fmt.Errorf("encode edit log error:%s", err)
	}

	if err := vs.manifest.Write(v); err != nil {
		return fmt.Errorf("write edit log error:%s", err)
	}

	newVersion := familyVersion.GetCurrent().cloneVersion()

	// apply delta edit to new version
	editLog.apply(newVersion)

	// Install the new version for family level version edit log
	familyVersion.appendVersion(newVersion)

	//TODO add detail edit log data
	vs.logger.Info("log and apply new version edit", zap.String("store", vs.storePath))
	return nil
}

// CreateFamilyVersion creates family version using family name,
// if family version exist, return exist one
func (vs *VersionSet) CreateFamilyVersion(family string) *FamilyVersion {
	var familyVersion = vs.GetFamilyVersion(family)
	if familyVersion != nil {
		vs.logger.Warn("family version exist, use it.", zap.String("store", vs.storePath), zap.String("family", family))
		return familyVersion
	}
	familyVersion = newFamilyVersion(vs)
	vs.mutex.Lock()
	vs.familyVersions[family] = familyVersion
	vs.mutex.Unlock()
	return familyVersion
}

// GetFamilyVersion returns family version if exist, else return nil
func (vs *VersionSet) GetFamilyVersion(family string) *FamilyVersion {
	vs.mutex.RLock()
	defer vs.mutex.RUnlock()
	familyVersion, ok := vs.familyVersions[family]
	if ok {
		return familyVersion
	}
	return nil
}

// initializeIfNeeded, initialize if version file not exists
// return true version set data ont exist, else has old data
func (vs *VersionSet) initializeIfNeeded() (bool, error) {
	if !util.Exist(filepath.Join(vs.storePath, current())) {
		vs.logger.Info("version set's current file not exist, initialze it", zap.String("store", vs.storePath))
		manifestFileName := manifestFileName(vs.manifestFileNumber) // manifest file name

		if err := vs.setCurrent(manifestFileName); err != nil {
			return true, err
		}

		manifest := filepath.Join(vs.storePath, manifestFileName) // manifest file path

		writer, err := journal.NewWriter(manifest)
		if err != nil {
			return true, err
		}

		vs.manifest = writer
		return true, nil
	}
	return false, nil
}

// newVersionID generates new version id
func (vs *VersionSet) newVersionID() int64 {
	vs.mutex.Lock()
	versionID := vs.versionID
	vs.versionID++
	vs.mutex.Unlock()
	return versionID
}

func (vs *VersionSet) setCurrent(manifestFile string) error {
	current := filepath.Join(vs.storePath, current()) // current file path
	tmp := fmt.Sprintf("%s.%s", current, tmpSuffix)
	// write manifest file name into current file
	if err := ioutil.WriteFile(tmp, []byte(manifestFile), 0666); err != nil {
		return fmt.Errorf("write manifest file name into current tmp file error:%s", err)
	}
	if err := os.Rename(tmp, current); err != nil {
		return fmt.Errorf("rename current tmp file name to current error:%s", err)
	}
	return nil
}
