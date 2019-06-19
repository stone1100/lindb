package kv

// FamilyOption defines config items for family level
type FamilyOption struct {
	Name string `toml:"name"`
}

// StoreOption defines config item for store level
type StoreOption struct {
	Path   string `toml:"path"`
	Levels int    `toml:"levels"`
}

// DefaultStoreOption build default store option
func DefaultStoreOption(path string) StoreOption {
	return StoreOption{
		Path:   path,
		Levels: 2,
	}
}
