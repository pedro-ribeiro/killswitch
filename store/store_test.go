package store

import (
	"fmt"
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

	if /*got == nil ||*/ got != want {
		t.Errorf("return value does not match")
	}

	//teardown
	cleanupStore(store)
}

func TestSubscriberNotifiedOnCreateUpdate(t *testing.T) {
	//setup
	store := createStore()
	feat := features.Feature{
		Key:         "key",
		Description: "a new and fancy feature",
		IsActive:    false,
	}
	store.UpsertFeature(feat)

	channel, err := store.SubscribeToUpdates()
	fmt.Printf(">> received channel %v\n", channel)

	if err != nil {
		t.Errorf("got error '%s'", err)
	}

	feat = features.Feature{
		Key:         "key",
		Description: "a new and fancy feature",
		IsActive:    true,
	}
	store.UpsertFeature(feat)

	got := <-channel

	if !got.IsActive {
		t.Errorf("return value does not match")
	}

	//teardown
	cleanupStore(store)
}

//Test for concurrency problems between getallfeatures & subscriber
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

	got, err := store.GetFeatureByKey("new-key")

	if err != nil {
		t.Errorf("got error '%s'", err)
	}

	if /*got == nil ||*/ got != want {
		t.Errorf("return value does not match")
	}

	//teardown
	cleanupStore(store)
}

func TestGetAllFeatures(t *testing.T) {
	//begin setup
	store := createStore()

	want1 := features.Feature{
		Key:         "new-key1",
		Description: "a new and fancy feature",
		IsActive:    false,
	}
	store.UpsertFeature(want1)

	want2 := features.Feature{
		Key:         "new-key2",
		Description: "an even fancier feature",
		IsActive:    true,
	}
	store.UpsertFeature(want2)
	//end setup

	got, err := store.GetAllFeatures()

	if err != nil {
		t.Errorf("got error '%s'", err)
	}

	if got == nil || len(got) != 2 {
		t.Errorf("return value is nil or wrong size")
		return
	}

	if val, ok := got[want1.Key]; !ok || val != want1 {
		t.Errorf("did not contain '%s' or want different from expected", want1.Key)
	}

	if val, ok := got[want2.Key]; !ok || val != want2 {
		t.Errorf("did not contain '%s' or want different from expected", want2.Key)
	}

	//teardown
	cleanupStore(store)
}

//TestGetAllWithNoFeatures
//TestGetAllWithFeatures
//TestDeleteFeature
//TestUpsertValidations

func TestGetNonExistingFeature(t *testing.T) {
	//setup
	store := createStore()

	_, err := store.GetFeatureByKey("invalid-key")

	if err == nil {
		t.Errorf("should've gotten NotFound error")
	} else {
		_, ok := err.(features.NotFoundError)
		if !ok {
			t.Errorf("should've gotten NotFound error: '%s'", err)
		}
	}

	// if got != nil {
	// 	t.Errorf("feature should've been nil: '%s:%s'", got.Key, got.Description)
	// }
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

func TestRedisStoreCompliesToFeatureStore(t *testing.T) {
	store, _ := NewRedisStore("test", "localhost:6379")
	var i interface{} = store
	_, ok := i.(features.FeatureStore)

	if !ok {
		t.Error("RedisStore does not comply to FeatureStore interface")
	}
}

func createStore() *RedisStore {
	store, _ := NewRedisStore("test", "localhost:6379")
	return store
}

func cleanupStore(s *RedisStore) {
	s.deleteAll()
}
