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

func NewRedisStore(name string, address string) (*RedisStore, error) {
	client := redis.NewClient(&redis.Options{
		Addr: address,
	})

	_, err := client.Ping().Result()

	return &RedisStore{client, name}, err
}

func (s *RedisStore) GetAllFeatures() (map[string](features.Feature), error) {
	index, err := s.client.SMembers(buildFeatureIndexKey(s)).Result()

	if err != nil {
		return nil, err
	}

	errChannel := make(chan error)
	valChannel := make(chan *features.Feature, len(index))

	for _, i := range index {
		cur := i
		go func() {
			val, err := s.GetFeatureByKey(cur)
			if err != nil {
				errChannel <- err
			} else {
				valChannel <- &val
			}
		}()
	}

	results := make(map[string](features.Feature))

	for i := 0; i < len(index); i++ {
		select {
		case err := <-errChannel:
			return nil, err
		case val := <-valChannel:
			results[val.Key] = *val
		}
	}

	return results, nil
}

func (s *RedisStore) GetFeatureByKey(key string) (features.Feature, error) {
	result, err := s.client.HGetAll(buildFeatureKey(s, key)).Result()

	if err != nil {
		return features.Feature{}, err
	}

	if len(result) == 0 {
		return features.Feature{}, features.NotFoundError{}
	}

	value, err := redisResultToFeature(result)
	return value, err
}

func (s *RedisStore) UpsertFeature(feature features.Feature) (features.Feature, error) {
	featureKey := buildFeatureKey(s, feature.Key)

	err := s.client.HMSet(featureKey, map[string]interface{}{
		"key":         feature.Key,
		"description": feature.Description,
		"isActive":    strconv.FormatBool(feature.IsActive),
	}).Err()

	if err != nil {
		return features.Feature{}, err
	}

	err = s.client.SAdd(buildFeatureIndexKey(s), feature.Key).Err()

	if err != nil {
		s.client.Del(featureKey)
		return features.Feature{}, err
	}

	return feature, nil
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

func redisResultToFeature(result map[string]string) (features.Feature, error) {
	boolValue, err := strconv.ParseBool(result["isActive"])

	if err != nil {
		return features.Feature{}, err
	}

	return features.Feature{
		Key:         result["key"],
		Description: result["description"],
		IsActive:    boolValue,
	}, nil
}

func (s *RedisStore) deleteAll() error {
	//TODO: this must only delete the RedisStore#name namespace
	return s.client.FlushAll().Err()
}
