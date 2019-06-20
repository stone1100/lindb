package kv

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/eleme/lindb/kv/table"
	"github.com/eleme/lindb/kv/version"
	"github.com/eleme/lindb/pkg/logger"
	"github.com/eleme/lindb/pkg/util"

	"go.uber.org/zap"
)

// Family implement column family for data isolation each family
type Family struct {
	store         *Store
	name          string
	option        FamilyOption
	familyVersion *version.FamilyVersion

	mutex sync.RWMutex

	logger *zap.Logger
}

// newFamily creates new family or open exsit family
func newFamily(store *Store, option FamilyOption) (*Family, error) {
	log := logger.GetLogger()
	name := option.Name

	familyPath := filepath.Join(store.option.Path, name)

	if !util.Exist(familyPath) {
		if err := util.MkDir(familyPath); err != nil {
			return nil, fmt.Errorf("mk family path error:%s", err)
		}
	}

	f := &Family{
		store:         store,
		name:          name,
		option:        option,
		familyVersion: store.versions.CreateFamilyVersion(name, option.ID),
		logger:        log,
	}

	log.Info("new family success", f.logStoreField(), f.logFamilyField())
	return f, nil
}

// NewTableBuilder create table builder instance for storing kv data
func (f *Family) NewTableBuilder() table.Builder {
	fileNumber := f.store.versions.NextFileNumber()

	fileName := filepath.Join(f.store.option.Path, f.name, version.Table(fileNumber))

	f.logger.Info(fileName)

	return nil
}

// CommitEditLog peresists eidt logs into mamanifest file
// returns ture commit successfully, else failure
func (f *Family) CommitEditLog(editLog *version.EditLog) bool {
	if editLog.IsEmpty() {
		f.logger.Warn("edit log is empty", f.logStoreField(), f.logFamilyField())
		return true
	}
	if err := f.store.versions.CommitFamilyEditLog(f.name, editLog); err != nil {
		f.logger.Error("commit edit log error:", f.logStoreField(), f.logFamilyField(), zap.Error(err))
		return false
	}
	return true
}

// GetSnapshot returns current version, includes sst files
func (f *Family) GetSnapshot() *Snapshot {

	current := f.familyVersion.GetCurrent()

	return newSnapshot(current)
}

// delete obsolete family sst files
func (f *Family) deleteObsoleteFiles() {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	//make a set of all of the live files
	/*
		Set<Long> live = newHashSet();
		live.addAll(this.getTableFiles());

		/*
		 * add live rollup reference files, maybe some roll up files is not alive but rollup job need it,
		 * so those files cannot delete, because read these files when do rollup job.
		 //*/
	//live.addAll(this.kvStore.getLiveReferenceFiles());
	//
	//List<File> files = Lists.newArrayList();
	//
	//files.addAll(FileName.listFiles(path));
	//
	//for (File file : files) {
	//FileName.FileInfo fileInfo = FileName.parseFileName(file);
	//if (fileInfo != null
	//&& fileInfo.getFileType() == FileName.FileType.SST
	//&& !live.contains(fileInfo.getFileNumber())) {
	//// 1.evict file table reader from cache, if exist
	//tableCache.evict(this, fileInfo.getFileNumber());
	//// 2.delete sst file
	//if (file.delete()) {
	//LoggerUtil.info(familyInfo, "delete file type [{}] successfully, file number[{}].",
	//fileInfo.getFileType(), fileInfo.getFileNumber());
	//}
	//}
	// */
}

// logStoreField logging store info
func (f *Family) logStoreField() zap.Field {
	return zap.String("store", f.store.option.Path)
}

// logFamilyField logging family info
func (f *Family) logFamilyField() zap.Field {
	return zap.String("family", f.name)
}
