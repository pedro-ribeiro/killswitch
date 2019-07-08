package store

import (
	"killswitch/features"
	"strconv"

	"github.com/go-redis/redis"
)

type RedisStore struct {
	client *redis.Client
	name   string
}

type FeatureCrudifier interface {
	GetFeature(key string) (*features.Feature, error)
	UpsertFeature(features.Feature) (*features.Feature, error)
	deleteAll() error
}

func NewRedisStore(name string, address string) (FeatureCrudifier, error) {
	client := redis.NewClient(&redis.Options{
		Addr: address,
	})

	_, err := client.Ping().Result()

	return &RedisStore{client, name}, err
}

func (s *RedisStore) UpsertFeature(feature features.Feature) (*features.Feature, error) {
	featureKey := buildFeatureKey(s, feature.Key)

	err := s.client.HMSet(featureKey, map[string]interface{}{
		"key":         feature.Key,
		"description": feature.Description,
		"isActive":    strconv.FormatBool(feature.IsActive),
	}).Err()

	if err != nil {
		return nil, err
	}

	err = s.client.SAdd(buildFeatureIndexKey(s), featureKey).Err()

	if err != nil {
		s.client.Del(featureKey)
		return nil, err
	}

	return &feature, nil
}

func buildFeatureKey(s *RedisStore, key string) string {
	return buildFeatureKeyPrefix(s) + key
}

func buildFeatureKeyPrefix(s *RedisStore) string {
	return s.name + ":feature:"
}

func buildFeatureIndexKey(s *RedisStore) string {
	return s.name + ":feature_index"
}

func (s *RedisStore) GetFeature(key string) (*features.Feature, error) {
	result, err := s.client.HGetAll(buildFeatureKey(s, key)).Result()

	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	boolValue, err := strconv.ParseBool(result["isActive"])

	if err != nil {
		return nil, err
	}

	return &features.Feature{
		Key:         result["key"],
		Description: result["description"],
		IsActive:    bool(boolValue),
	}, nil
}

func (s *RedisStore) deleteAll() error {
	//TODO: this must only delete the RedisStore#name namespace
	return s.client.FlushAll().Err()
}
