package version

import "sync"

// FamilyVersion maintains family level metadata
type FamilyVersion struct {
	versionSet *VersionSet

	current        *Version           // current mutable version
	activeVersions map[int64]*Version // all active versions include mutable/immutable versions

	mutex sync.RWMutex
}

// newFamilyVersion new FamilyVersion instance
func newFamilyVersion(versionSet *VersionSet) *FamilyVersion {
	fv := &FamilyVersion{
		versionSet:     versionSet,
		activeVersions: make(map[int64]*Version),
	}

	// create new version for current mutable version
	fv.createVersion()

	return fv
}

// GetCurrent returns current mutable version
func (fv *FamilyVersion) GetCurrent() *Version {
	fv.mutex.RLock()
	defer fv.mutex.Unlock()

	// inc ref count of version
	fv.current.retain()

	return fv.current
}

// removeVersion removes version from active versions
func (fv *FamilyVersion) removeVersion(v *Version) {
	fv.mutex.Lock()
	delete(fv.activeVersions, v.id)
	fv.mutex.Unlock()
}

//createVersion creates new version
func (fv *FamilyVersion) createVersion() {
	fv.mutex.Lock()
	current := newVersion(fv.versionSet.newVersionID(), fv)
	fv.activeVersions[current.id] = current
	fv.current = current
	fv.mutex.Unlock()
}

// appendVersion swap famliy's current version, then release previous version
func (fv *FamilyVersion) appendVersion(v *Version) {
	pervious := fv.current

	fv.activeVersions[v.id] = v
	fv.current = v

	if pervious != nil {
		pervious.Release()
	}
}
