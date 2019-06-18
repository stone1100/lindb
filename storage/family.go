package storage

import (
	"path/filepath"
	"sync"

	"github.com/eleme/lindb/pkg/logger"
	"github.com/eleme/lindb/pkg/util"
	"github.com/eleme/lindb/storage/table"
	meta "github.com/eleme/lindb/storage/version"

	"github.com/BurntSushi/toml"
	"go.uber.org/zap"
)

// Family implement column family for data isolation each family
type Family struct {
	store         *Store
	name          string
	option        FamilyOption
	familyVersion *meta.FamilyVersion

	mutex sync.RWMutex

	logger *zap.Logger
}

// newFamily creates family if it exist
func newFamily(store *Store, name string, option FamilyOption) (*Family, error) {
	log := logger.GetLogger()

	familyPath := filepath.Join(store.option.Path, name)

	if !util.Exist(familyPath) {
		if err := util.MkDir(familyPath); err != nil {
			return nil, err
		}

		optionFile := filepath.Join(familyPath, meta.Info())

		if err := util.EncodeToml(optionFile, option); err != nil {
			return nil, err
		}
	}

	f := &Family{
		store:         store,
		name:          name,
		option:        option,
		familyVersion: store.versions.CreateFamilyVersion(name),
		logger:        log,
	}

	log.Info("new family success", zap.String("family", name))
	return f, nil
}

// openFamily opens exist family
func openFamily(store *Store, name string) (*Family, error) {
	log := logger.GetLogger()
	optionFile := filepath.Join(store.option.Path, name, meta.Info())
	option := &FamilyOption{}

	if _, err := toml.DecodeFile(optionFile, option); err != nil {
		return nil, err
	}

	f := &Family{
		store:         store,
		name:          name,
		option:        *option,
		familyVersion: store.versions.CreateFamilyVersion(name),
		logger:        log,
	}

	log.Info("open family success", zap.String("family", name))
	return f, nil
}

// NewTableBuilder create table builder instance for storing kv data
func (f *Family) NewTableBuilder() table.Builder {
	fileNumber := f.store.versions.NextFileNumber()

	fileName := filepath.Join(f.store.option.Path, f.name, meta.Table(fileNumber))

	f.logger.Info(fileName)

	return nil
}

// CommitEditLog peresists eidt logs into mamanifest file
// returns ture commit successfully, else failure
func (f *Family) CommitEditLog(editLog *meta.EditLog) bool {
	if err := f.store.versions.Commit(f.name, editLog); err != nil {
		f.logger.Error("commit edit log error:", zap.String("family", f.name), zap.Error(err))
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
