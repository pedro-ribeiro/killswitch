package features

type Feature struct {
	Key, Description string
	IsActive         bool
}

type FeatureCrudifier interface {
	GetFeature(key string) (*Feature, error)
	UpsertFeature(Feature) (*Feature, error)
}

type NotFoundError struct {
	s string
}

func (e NotFoundError) Error() string {
	return e.s
}
