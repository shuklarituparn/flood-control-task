package redis_impelementation

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"strconv"
	"time"
)

type SlidingWindow struct {
	redisClient *redis.Client
	keyPrefix   string
	rate        float64
	window      time.Duration
	userID      int64
}

func NewSlidingWindow(redisClient *redis.Client, keyPrefix string, userID int64, rate float64, window time.Duration) *SlidingWindow {
	return &SlidingWindow{
		redisClient: redisClient,
		keyPrefix:   keyPrefix,
		rate:        rate,
		window:      window,
		userID:      userID,
	}
}

func (sw *SlidingWindow) increment() error {
	now := time.Now().UnixNano()
	score := float64(now)
	keyPrefix := sw.keyPrefix + ":" + strconv.Itoa(int(sw.userID))
	member := fmt.Sprintf("%d", now)
	_, err := sw.redisClient.ZAdd(context.Background(), keyPrefix, redis.Z{
		Score:  score,
		Member: member,
	}).Result()
	if err != nil {
		return err
	}
	return nil
}

func (sw *SlidingWindow) removeExpired() error {
	now := time.Now().UnixNano()
	minScore := float64(now) - sw.window.Seconds()*1e9
	keyPrefix := sw.keyPrefix + ":" + strconv.Itoa(int(sw.userID))
	_, err := sw.redisClient.ZRemRangeByScore(context.Background(), keyPrefix, "0", fmt.Sprintf("%.0f", minScore)).Result()
	if err != nil {
		return err
	}
	return nil
}

func (sw *SlidingWindow) countRequests() (int64, error) {
	now := time.Now().UnixNano()
	minScore := float64(now) - sw.window.Seconds()*1e9
	keyPrefix := sw.keyPrefix + ":" + strconv.Itoa(int(sw.userID))
	count, err := sw.redisClient.ZCount(context.Background(), keyPrefix, fmt.Sprintf("%.0f", minScore), "+inf").Result()
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (sw *SlidingWindow) Check(context.Context, int64) (bool, error) {
	err := sw.increment()
	if err != nil {
		log.Println("Error incrementing sliding window:", err)
		return false, err
	}
	err = sw.removeExpired()
	if err != nil {
		log.Println("Error removing expired sliding window entries:", err)
		return false, err
	}
	count, err := sw.countRequests()
	if err != nil {
		log.Println("Error counting sliding window requests:", err)
		return false, err
	}
	allowedRequests := int64(sw.rate * sw.window.Seconds())
	return count <= allowedRequests, nil
}
