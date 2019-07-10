package store

import (
	"killswitch/features"
	"testing"
)

func TestCreateFeature(t *testing.T) {
	//setup
	store := createStore()

	got, err := store.UpsertFeature(
		features.Feature{
			Key:         "new-key",
			Description: "a new and fancy feature",
		},
	)

	want := features.Feature{
		Key:         "new-key",
		Description: "a new and fancy feature",
		IsActive:    false,
	}

	if err != nil {
		t.Errorf("got error '%s'", err)
	}

	if got == nil || *got != want {
		t.Errorf("return value does not match")
	}

	//teardown
	cleanupStore(store)
}

//TestUpdateFeature

func TestGetExistingFeature(t *testing.T) {
	//begin setup
	store := createStore()

	want := features.Feature{
		Key:         "new-key",
		Description: "a new and fancy feature",
		IsActive:    false,
	}

	store.UpsertFeature(want)
	//end setup

	got, err := store.GetFeature("new-key")

	if err != nil {
		t.Errorf("got error '%s'", err)
	}

	if got == nil || *got != want {
		t.Errorf("return value does not match")
	}

	//teardown
	cleanupStore(store)
}

//TestGetAllWithNoFeatures
//TestGetAllWithFeatures
//TestDeleteFeature

func TestGetNonExistingFeature(t *testing.T) {
	//setup
	store := createStore()

	got, err := store.GetFeature("invalid-key")

	if err == nil {
		t.Errorf("should've gotten NotFound error")
	} else {
		_, ok := err.(features.NotFoundError)
		if !ok {
			t.Errorf("should've gotten NotFound error: '%s'", err)
		}
	}

	if got != nil {
		t.Errorf("feature should've been nil: '%s:%s'", got.Key, got.Description)
	}
}

func TestCreateStoreValidAddress(t *testing.T) {
	_, err := NewRedisStore("test", "localhost:6379")

	if err != nil {
		t.Errorf("got error '%s'", err)
	}
}

func TestCreateStoreInvalidAddress(t *testing.T) {
	_, err := NewRedisStore("test", "ali:12")

	if err == nil {
		t.Error("it was supposed to fail store creation")
	}
}

func TestRedisStoreCompliesToFeatureCrudifier(t *testing.T) {
	store, _ := NewRedisStore("test", "localhost:6379")
	var i interface{} = store
	_, ok := i.(features.FeatureCrudifier)

	if !ok {
		t.Error("RedisStore does not comply to FeatureCrudifier interface")
	}
}

func createStore() *RedisStore {
	store, _ := NewRedisStore("test", "localhost:6379")
	return store
}

func cleanupStore(s *RedisStore) {
	s.deleteAll()
}
