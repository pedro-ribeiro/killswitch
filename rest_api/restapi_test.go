package restapi

import (
	"errors"
	"killswitch/features"
	"testing"
)

func TestShouldReturnFeatureOnExistingKey(t *testing.T) {
	store := &alwaysSuccessFeaturesStore{}
	fixture := store.getUniqueFeature()

	want := FeatureGetterResponse{
		Key:         fixture.Key,
		Description: fixture.Description,
		IsActive:    fixture.IsActive,
	}

	var got FeatureGetterResponse

	got, err := getFeatureFromStore(store, &featureGetterRequest{fixture.Key})

	if err != nil {
		t.Errorf("Unexpected error: '%s'", err)
	}

	if got != want {
		t.Errorf("Response did not match expected value: '%v'", got)
	}
}

func TestShouldReturnNotFoundErrorOnNonExistingKey(t *testing.T) {
	store := &alwaysFailureFeaturesStore{}

	_, err := getFeatureFromStore(store, &featureGetterRequest{"some-key"})

	if err == nil {
		t.Errorf("Expected a NotFound error")
	} else {
		_, ok := err.(FeatureNotFoundError)
		if !ok {
			t.Errorf("Expected a FeatureNotFound error '%s'", err)
		}
	}
}

type alwaysSuccessFeaturesStore struct {
}

func (store *alwaysSuccessFeaturesStore) getUniqueFeature() *features.Feature {
	return &features.Feature{
		Description: "a description",
		Key:         "my-key",
		IsActive:    true,
	}
}

func (store *alwaysSuccessFeaturesStore) GetFeature(key string) (*features.Feature, error) {
	return store.getUniqueFeature(), nil
}

func (store *alwaysSuccessFeaturesStore) UpsertFeature(feature features.Feature) (*features.Feature, error) {
	return &feature, nil
}

type alwaysFailureFeaturesStore struct {
}

func (store *alwaysFailureFeaturesStore) GetFeature(key string) (*features.Feature, error) {
	return nil, features.NotFoundError{}
}

func (store *alwaysFailureFeaturesStore) UpsertFeature(feature features.Feature) (*features.Feature, error) {
	return nil, errors.New("generic error")
}
