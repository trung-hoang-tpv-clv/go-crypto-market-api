package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"go-crypto-market-api/internal/domain"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCache interface {
	FetchPriceHistory(ctx context.Context, symbol string, from, to time.Time) ([]domain.PriceHistory, error)
	SetCachePriceHistory(ctx context.Context, symbol string, from, to time.Time, data []domain.PriceHistory) error
}

var _ RedisCache = &RedisCacheImpl{}

type RedisCacheImpl struct {
	client *redis.Client
}

func NewRedisCacheImpl(client *redis.Client) *RedisCacheImpl {
	return &RedisCacheImpl{
		client: client,
	}
}

func (r *RedisCacheImpl) FetchPriceHistory(ctx context.Context, symbol string, from, to time.Time) ([]domain.PriceHistory, error) {
	key := fmt.Sprintf("price_history:%s:%d:%d", symbol, from.Unix(), to.Unix())
	fmt.Println("Cache key", key)
	data, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var priceHistories []domain.PriceHistory
	err = json.Unmarshal([]byte(data), &priceHistories)
	if err != nil {
		return nil, err
	}

	return priceHistories, nil
}

func (r *RedisCacheImpl) SetCachePriceHistory(ctx context.Context, symbol string, from, to time.Time, data []domain.PriceHistory) error {
	if len(data) == 0 {
		return nil
	}
	key := fmt.Sprintf("price_history:%s:%d:%d", symbol, from.Unix(), to.Unix())
	fmt.Println("Cache key", key)
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Assuming a 24-hour expiration for the cache
	if err := r.client.Set(ctx, key, jsonData, 24*time.Hour).Err(); err != nil {
		return err
	}

	return nil
}
