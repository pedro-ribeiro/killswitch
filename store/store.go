package store

import (
	"fmt"
	"killswitch/features"
	"strconv"
	"sync"

	"github.com/go-redis/redis"
)

type RedisStore struct {
	client        *redis.Client
	name          string
	updateChannel chan features.Feature
	subscribers   []chan features.Feature
	mux           sync.Mutex
}

func NewRedisStore(name string, address string) (*RedisStore, error) {
	client := redis.NewClient(&redis.Options{
		Addr: address,
	})

	updateChannel := make(chan features.Feature)

	_, err := client.Ping().Result()

	store := RedisStore{
		client:        client,
		name:          name,
		updateChannel: updateChannel,
		subscribers:   make([]chan features.Feature, 0)}

	go store.notifySubscribers()

	return &store, err
}

func (s *RedisStore) SubscribeToUpdates() (chan features.Feature, error) {
	channel := make(chan features.Feature)

	s.mux.Lock()
	defer s.mux.Unlock()

	s.subscribers = append(s.subscribers, channel)

	fmt.Printf("new subscriber!\n")

	return channel, nil
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

	select {
	case s.updateChannel <- feature:
	default:
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

func (s *RedisStore) notifySubscribers() {
	fmt.Printf("notify initialized!\n")
	defer s.mux.Unlock()
	for f := range s.updateChannel {
		// for {
		// 	var f features.Feature

		// 	select {
		// 	case f = <-s.updateChannel:
		// 	default:
		// 	}

		fmt.Printf("new feature: %v\n", f)
		s.mux.Lock()

		fmt.Printf("lock acquired.\n")
		for _, c := range s.subscribers {
			fmt.Printf("notifying one subscriber\n")
			select {
			case c <- f:
				fmt.Println("sent to subscriber")
				continue
			default:
				fmt.Println("didnt send")
			}
		}

		s.mux.Unlock()
	}
}
