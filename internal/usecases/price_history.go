package usecases

import (
	"context"
	"fmt"
	"go-crypto-market-api/internal/domain"
	"go-crypto-market-api/internal/infrastructure/cache"
	"go-crypto-market-api/internal/infrastructure/coingecko"
	"go-crypto-market-api/internal/infrastructure/repository"
	"go-crypto-market-api/internal/interfaces/dto"
	"time"
)

type PriceHistoryUsecase struct {
	repo     repository.PriceHistoryRepository
	cache    cache.RedisCache
	cgClient *coingecko.Client
}

func NewPriceHistoryUsecase(repo repository.PriceHistoryRepository, cache cache.RedisCache, cgClient *coingecko.Client) *PriceHistoryUsecase {
	return &PriceHistoryUsecase{
		repo:     repo,
		cache:    cache,
		cgClient: cgClient,
	}
}

func (u *PriceHistoryUsecase) GetPriceHistories(ctx context.Context, symbol string, startDate, endDate string) ([]dto.GetPriceHistory, error) {
	// Parse the start and end dates
	startTime, err := time.Parse(time.DateOnly, startDate)
	if err != nil {
		return nil, err
	}
	endTime, err := time.Parse(time.DateOnly, endDate)
	if err != nil {
		return nil, err
	}

	// Get the data from the cache
	data, err := u.cache.FetchPriceHistory(ctx, symbol, startTime, endTime)
	if err != nil {
		return nil, err
	}

	if len(data) > 0 {
		fmt.Println("debug: fetch data from cache", data)
		return transformData(data), nil
	}

	// If cache missed data then get from database
	data, err = u.repo.FetchPriceHistories(ctx, symbol, startTime, endTime)
	if err != nil {
		return nil, err
	}

	if len(data) > 0 {
		fmt.Println("debug: fetch data from database", data)
		// Save data from database to redis cache
		err = u.cache.SetCachePriceHistory(ctx, symbol, startTime, endTime, data)
		if err != nil {
			fmt.Println("Save data from database to redis cache error", err)
		}
		return transformData(data), nil
	}

	// fetch data from CoinGecko
	var cgResp *coingecko.CoinGeckoResponse
	cgResp, err = u.cgClient.FetchPriceHistories(ctx, symbol, startTime, endTime)
	if err != nil {
		return nil, err
	}

	if len(cgResp.Prices) == 0 {
		return nil, nil
	}

	high := cgResp.Prices[0][1]
	low := cgResp.Prices[0][1]
	open := cgResp.Prices[0][1]
	close := cgResp.Prices[len(cgResp.Prices)-1][1]
	change := 0.0

	// Calculate high, low, and change
	for _, entry := range cgResp.Prices {
		price := entry[1]
		if price > high {
			high = price
		}
		if price < low {
			low = price
		}
	}

	// Calculate the change from the opening price to the closing price
	if open != 0 {
		change = (close - open) / open * 100
	}

	priceHistory := domain.PriceHistory{
		Symbol:    symbol,
		High:      high,
		Low:       low,
		Open:      open,
		Close:     close,
		Time:      int64(cgResp.Prices[0][0] / 1000),
		Change:    change,
		StartDate: startTime,
		EndDate:   endTime,
	}

	// Save the fetched data in the cache and database for future requests
	err = u.cache.SetCachePriceHistory(ctx, symbol, startTime, endTime, []domain.PriceHistory{priceHistory})
	if err != nil {
		fmt.Println("Save price history to redis cache error:", err)
	}

	err = u.repo.SavePriceHistories(ctx, []domain.PriceHistory{priceHistory})
	if err != nil {
		fmt.Println("Save price history to database error:", err)
		return nil, nil
	}

	return transformData([]domain.PriceHistory{priceHistory}), nil
}

func transformData(data []domain.PriceHistory) []dto.GetPriceHistory {
	var priceHistoriesDTO []dto.GetPriceHistory

	for _, item := range data {
		priceHistoriesDTO = append(priceHistoriesDTO, dto.GetPriceHistory{
			High:   item.High,
			Low:    item.Low,
			Open:   item.Open,
			Close:  item.Close,
			Time:   item.Time,
			Change: item.Change,
		})
	}
	return priceHistoriesDTO
}
