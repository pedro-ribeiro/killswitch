package features

type Feature struct {
	Key, Description string
	IsActive         bool
}

type FeatureStore interface {
	GetFeatureByKey(key string) (Feature, error)
	UpsertFeature(value Feature) (Feature, error)
	GetAllFeatures() (map[string]Feature, error)
}

type NotFoundError struct {
	s string
}

func (e NotFoundError) Error() string {
	return e.s
}
