package restapi

import (
	"errors"
	"killswitch/features"
	"testing"
)

func TestShouldReturnFeatureOnExistingKey(t *testing.T) {
	//setup
	store := &alwaysSuccessFeaturesStore{}
	fixture := store.getUniqueFeature()
	store.UpsertFeature(fixture)

	want := FeatureResponse{
		Key:         fixture.Key,
		Description: fixture.Description,
		IsActive:    fixture.IsActive,
	}

	var got FeatureResponse

	got, err := getFeatureFromStore(store, &FeatureGetterRequest{fixture.Key})

	if err != nil {
		t.Errorf("Unexpected error: '%s'", err)
	}

	if got != want {
		t.Errorf("Response did not match expected value: '%v'", got)
	}
}

func TestShouldReturnNotFoundErrorOnNonExistingKey(t *testing.T) {
	store := &alwaysFailureFeaturesStore{}

	_, err := getFeatureFromStore(store, &FeatureGetterRequest{"some-key"})

	if err == nil {
		t.Errorf("Expected a NotFound error")
	} else {
		_, ok := err.(RestNotFoundError)
		if !ok {
			t.Errorf("Expected a FeatureNotFound error '%s'", err)
		}
	}
}

func TestShouldReturnEmptyResponseOnEmptyStore(t *testing.T) {
	store := &alwaysSuccessFeaturesStore{}

	response, err := getAllFeaturesFromStore(store)

	if err != nil {
		t.Errorf("Did not expect an error")
	}

	if len(response.Results) != 0 {
		t.Errorf("Expected empty results")
	}
}

func TestShouldReturnTwoResults(t *testing.T) {
	//setup
	store := &alwaysSuccessFeaturesStore{}
	feat1 := features.Feature{
		Key:         "my-key1",
		Description: "description 1",
		IsActive:    false,
	}
	feat2 := features.Feature{
		Key:         "my-key2",
		Description: "description 2",
		IsActive:    true,
	}
	store.UpsertFeature(feat1) //FIXME: access the storage map directly
	store.UpsertFeature(feat2)

	response, err := getAllFeaturesFromStore(store)

	if err != nil {
		t.Errorf("Did not expect an error")
	}

	if len(response.Results) != 2 {
		t.Errorf("Expected empty results")
	}

	if val, ok := response.Results[feat1.Key]; !ok || val.Key != feat1.Key {
		t.Errorf("expected something similar to feature1")
	}

	if val, ok := response.Results[feat2.Key]; !ok || val.Key != feat2.Key {
		t.Errorf("expected something similar to feature2")
	}
}

func TestShouldReturnErrorOnFailingStore(t *testing.T) {
	store := &alwaysFailureFeaturesStore{}

	_, err := getAllFeaturesFromStore(store)

	if err == nil {
		t.Errorf("expected error to be returned")
	}
}

func TestShouldCreateEntryInStoreOnFeatureCreate(t *testing.T) {
	store := &alwaysSuccessFeaturesStore{}

	_, err := upsertFeatureInStore(store,
		&FeatureUpsertRequest{
			Description: "a description",
			Key:         "my-key",
			IsActive:    true,
		})

	if err != nil {
		t.Errorf("Did not expect an error '%s'", err)
	}

	if _, ok := store.storage["my-key"]; !ok {
		t.Errorf("Store did not save the feature!")
	}
}

type alwaysSuccessFeaturesStore struct {
	storage map[string]features.Feature
}

func (store *alwaysSuccessFeaturesStore) getUniqueFeature() features.Feature {
	return features.Feature{
		Description: "a description",
		Key:         "my-key",
		IsActive:    true,
	}
}

func (store *alwaysSuccessFeaturesStore) GetFeatureByKey(key string) (features.Feature, error) {
	return store.storage[key], nil
}

func (store *alwaysSuccessFeaturesStore) GetAllFeatures() (map[string]features.Feature, error) {
	return store.storage, nil
}

func (store *alwaysSuccessFeaturesStore) UpsertFeature(feature features.Feature) (features.Feature, error) {
	if store.storage == nil {
		storage := make(map[string]features.Feature)
		store.storage = storage
	}
	store.storage[feature.Key] = feature
	return feature, nil
}

type alwaysFailureFeaturesStore struct {
}

func (store *alwaysFailureFeaturesStore) GetFeatureByKey(key string) (features.Feature, error) {
	return features.Feature{}, features.NotFoundError{}
}

func (store *alwaysFailureFeaturesStore) GetAllFeatures() (map[string]features.Feature, error) {
	return nil, errors.New("generic error")
}

func (store *alwaysFailureFeaturesStore) UpsertFeature(feature features.Feature) (features.Feature, error) {
	return features.Feature{}, errors.New("generic error")
}
