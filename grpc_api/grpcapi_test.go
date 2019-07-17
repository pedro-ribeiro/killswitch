package grpcapi

import (
	"killswitch/features"
	"testing"
)

func TestGetFeaturesWithNoFeatures(t *testing.T) {

	//FIXME: this is probably the right way
	// server := NewGrpcServer(store)

	//setup
	values := make(map[string]features.Feature)
	store := &alwaysSuccessStore{values}

	stream, err := getFeaturesStream(store)

	if err != nil {
		t.Errorf("did not expect an error %s", err)
	}

	for i := range stream {
		t.Errorf("did not expect any features: %v", i)
	}
}

func TestGetFeaturesWithNoUpdates(t *testing.T) {
	//setup
	values := make(map[string]features.Feature)
	values["key1"] = features.Feature{
		Key:         "key1",
		Description: "desc1",
		IsActive:    false,
	}
	values["key2"] = features.Feature{
		Key:         "key2",
		Description: "desc2",
		IsActive:    true,
	}
	store := &alwaysSuccessStore{values}

	stream, err := getFeaturesStream(store)

	if err != nil {
		t.Errorf("did not expect an error %s", err)
	}

	cur := 0

	want := make([]string, 2)
	want[0] = "key1"
	want[1] = "key2"

	for i := range stream {
		if i.Key != want[cur] {
			t.Errorf("expected a different key: %s", want[cur])
		}
		cur++
	}

	if cur != 2 {
		t.Errorf("expected a different number of features. 2 vs %d", cur)
	}
}

type alwaysSuccessStore struct {
	storage map[string]features.Feature
}

func (s *alwaysSuccessStore) GetFeatureByKey(key string) (features.Feature, error) {
	return features.Feature{}, nil
}
func (s *alwaysSuccessStore) UpsertFeature(value features.Feature) (features.Feature, error) {
	return features.Feature{}, nil
}
func (s *alwaysSuccessStore) GetAllFeatures() (map[string]features.Feature, error) {
	return s.storage, nil
}

func (s *alwaysSuccessStore) addFeature(value features.Feature) {
	s.storage[value.Key] = value
}
