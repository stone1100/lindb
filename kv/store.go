package kv

import (
	"path/filepath"
	"sync"

	"github.com/eleme/lindb/pkg/logger"
	"github.com/eleme/lindb/pkg/util"
	"github.com/eleme/lindb/kv/version"

	"go.uber.org/zap"
)

// Store is kv store, support column family, but differnt other LSM implements.
// current implement not include memory table write logic.
type Store struct {
	name     string
	option   StoreOption
	lock     *Lock // file lock make sure store only been open once instance
	versions *version.VersionSet
	families map[string]*Family

	mutex sync.RWMutex

	logger *zap.Logger
}

// NewStore new store instance, need recover data if store existent
func NewStore(name string, option StoreOption) (*Store, error) {
	if err := util.MkDirIfNotExist(option.Path); err != nil {
		return nil, err
	}

	// file lock, only allow open by a instance
	lock := NewLock(filepath.Join(option.Path, version.Lock))
	err := lock.Lock()
	if err != nil {
		return nil, err
	}

	log := logger.GetLogger()

	// unlock file lock if error
	defer func() {
		if err != nil {
			e := lock.Unlock()
			if e != nil {
				log.Error("unlock file error:", zap.Error(e))
			}
		}
	}()

	store := &Store{
		name:     name,
		option:   option,
		lock:     lock,
		families: make(map[string]*Family),
		logger:   log,
	}

	// init and recover version set
	vs := version.NewVersionSet(store.option.Path)
	err = vs.Recover()
	if err != nil {
		return nil, err
	}

	store.versions = vs
	return store, nil
}

// CreateFamily create/load column family.
func (s *Store) CreateFamily(familyName string, option FamilyOption) (*Family, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var family, ok = s.families[familyName]
	if !ok {
		familyPath := filepath.Join(s.option.Path, familyName)

		var err error
		if !util.Exist(familyPath) {
			// create new family
			option.Name = familyName
			family, err = newFamily(s, familyName, option)
		} else {
			// open exist family
			family, err = openFamily(s, familyName)
		}

		if err != nil {
			return nil, err
		}
		s.families[familyName] = family
	}

	return family, nil
}

// GetFamily gets family based on name, if not exist return nil
func (s *Store) GetFamily(familyName string) (*Family, bool) {
	s.mutex.Lock()
	family, ok := s.families[familyName]
	s.mutex.Unlock()
	return family, ok
}

// Close closes store, then release some resoure
func (s *Store) Close() error {
	return s.lock.Unlock()
}
