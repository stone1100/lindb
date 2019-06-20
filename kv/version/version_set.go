package version

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	"go.uber.org/zap"

	"github.com/eleme/lindb/kv/journal"
	"github.com/eleme/lindb/pkg/logger"
	"github.com/eleme/lindb/pkg/util"
)

// StoreVersionSet maintains all metadata for kv store
type StoreVersionSet struct {
	manifestFileNumber int64
	nextFileNumber     int64
	storePath          string
	familyVersions     map[string]*FamilyVersion
	familyIDs          map[int]string
	versionID          int64 // unique in for increasing version id

	numOfLevels int // num of levels

	manifest *journal.Writer
	mutex    sync.RWMutex

	logger *zap.Logger
}

// NewStoreVersionSet new VersionSet instance
func NewStoreVersionSet(storePath string, numOfLevels int) *StoreVersionSet {
	return &StoreVersionSet{
		manifestFileNumber: 1, // default value for initialize store
		nextFileNumber:     2, // default value
		storePath:          storePath,
		numOfLevels:        numOfLevels,
		familyVersions:     make(map[string]*FamilyVersion),
		familyIDs:          make(map[int]string),
		logger:             logger.GetLogger(),
	}
}

// Destroy closes version set, release resource, such as journal writer etc.
func (vs *StoreVersionSet) Destroy() {
	vs.mutex.Lock()
	defer vs.mutex.Unlock()

	// close manifest journal writer if it exist
	if vs.manifest != nil {
		vs.manifest.Close()
	}
}

// NextFileNumber generates next file number
func (vs *StoreVersionSet) NextFileNumber() int64 {
	nextNumber := atomic.AddInt64(&vs.nextFileNumber, 1)
	return nextNumber - 1
}

// CommitFamilyEditLog peresists edit logs to manifest file, then apply new version to family version
func (vs *StoreVersionSet) CommitFamilyEditLog(family string, editLog *EditLog) error {
	// get family version based on family name
	familyVersion := vs.GetFamilyVersion(family)
	if familyVersion == nil {
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
	if err := vs.manifest.Sync(); err != nil {
		return fmt.Errorf("sync edit log error:%s", err)
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
func (vs *StoreVersionSet) CreateFamilyVersion(family string, familyID int) *FamilyVersion {
	var familyVersion = vs.GetFamilyVersion(family)
	if familyVersion != nil {
		vs.logger.Warn("family version exist, use it.", zap.String("store", vs.storePath), zap.String("family", family))
		return familyVersion
	}
	familyVersion = newFamilyVersion(vs)
	vs.mutex.Lock()
	vs.familyVersions[family] = familyVersion
	vs.familyIDs[familyID] = family
	vs.mutex.Unlock()
	return familyVersion
}

// GetFamilyVersion returns family version if exist, else return nil
func (vs *StoreVersionSet) GetFamilyVersion(family string) *FamilyVersion {
	vs.mutex.RLock()
	defer vs.mutex.RUnlock()
	familyVersion, ok := vs.familyVersions[family]
	if ok {
		return familyVersion
	}
	return nil
}

// Recover recover version set if exist, recover been invoked when kv store init.
// Initialize if version file not exists, else recover old data then init journal writer.
func (vs *StoreVersionSet) Recover() error {
	if !util.Exist(filepath.Join(vs.storePath, current())) {
		vs.logger.Info("version set's current file not exist, initialize it", zap.String("store", vs.storePath))
		if err := vs.initJournal(); err != nil {
			return err
		}
		return nil
	}
	vs.logger.Info("recover version set data from journal file", zap.String("store", vs.storePath))
	if err := vs.recover(); err != nil {
		return err
	}
	if err := vs.initJournal(); err != nil {
		return err
	}
	return nil
}

// recover does recover logic, read journal wal record and recover it
func (vs *StoreVersionSet) recover() error {
	manifestFileName, err := vs.readManifestFileName()
	if err != nil {
		return err
	}
	manifestPath := vs.getManifestFilePath(manifestFileName)
	reader, err := journal.NewReader(manifestPath)
	if err != nil {
		return err
	}
	// read edit log
	for {
		next, err := reader.Next()
		if err != nil {
			return fmt.Errorf("recover data from manifest file error:%s", err)
		}
		if !next {
			break
		}
		record := reader.Record()
		editLog := &EditLog{}
		unmalshalErr := editLog.unmarshal(record)
		if unmalshalErr != nil {
			return fmt.Errorf("unmarshal edit log data from manifest file error:%s", unmalshalErr)
		}

		familyID := editLog.familyID
		if familyID == StoreFamilyID {
			editLog.applyVersionSet(vs)
		} else {
			// find releted family version
			familyVersion := vs.getFamilyVersion(familyID)
			if familyVersion == nil {
				return fmt.Errorf("cannot get family version by id:%d", familyID)
			}
			// apply edit log to family current family
			editLog.apply(familyVersion.GetCurrent())
		}
	}
	return nil
}

// setNextFileNumberWithoutLock set next file number, invoker must add lock
func (vs *StoreVersionSet) setNextFileNumberWithoutLock(newNextFileNumber int64) {
	vs.manifestFileNumber = newNextFileNumber
	vs.nextFileNumber = newNextFileNumber + 1
}

// readManifestFileName reads manifest file name from current file
func (vs *StoreVersionSet) readManifestFileName() (string, error) {
	current := vs.getCurrentPath()
	v, err := ioutil.ReadFile(current)
	if err != nil {
		return "", fmt.Errorf("write manifest file name error:%s", err)
	}
	return string(v), nil
}

func (vs *StoreVersionSet) initJournal() error {
	if vs.manifest == nil {
		manifestFileName := manifestFileName(vs.manifestFileNumber) // manifest file name
		manifest := vs.getManifestFilePath(manifestFileName)
		if err := vs.setCurrent(manifestFileName); err != nil {
			return err
		}
		writer, err := journal.NewWriter(manifest)
		if err != nil {
			return err
		}
		vs.manifest = writer
	}
	return nil
}

// getFamilyVersion returns family version
func (vs *StoreVersionSet) getFamilyVersion(familyID int) *FamilyVersion {
	vs.mutex.RLock()
	defer vs.mutex.RUnlock()
	familyName, ok := vs.familyIDs[familyID]
	if !ok {
		return nil
	}
	familyVerion := vs.familyVersions[familyName]
	return familyVerion
}

// newVersionID generates new version id
func (vs *StoreVersionSet) newVersionID() int64 {
	newID := atomic.AddInt64(&vs.versionID, 1)
	return newID - 1
}

// setCurrent writes manifest file name into CURRENT file
func (vs *StoreVersionSet) setCurrent(manifestFile string) error {
	current := vs.getCurrentPath()
	tmp := fmt.Sprintf("%s.%s", current, TmpSuffix)
	// write manifest file name into current file
	if err := ioutil.WriteFile(tmp, []byte(manifestFile), 0666); err != nil {
		return fmt.Errorf("write manifest file name into current tmp file error:%s", err)
	}
	if err := os.Rename(tmp, current); err != nil {
		return fmt.Errorf("rename current tmp file name to current error:%s", err)
	}
	return nil
}

// getCurrent returns current file path
func (vs *StoreVersionSet) getCurrentPath() string {
	return filepath.Join(vs.storePath, current())
}

// getMainfiestFilePath returns manifest file path
func (vs *StoreVersionSet) getManifestFilePath(manifestFileName string) string {
	return filepath.Join(vs.storePath, manifestFileName)
}
