package repository

import (
	"context"
	"go-crypto-market-api/internal/config"
	"go-crypto-market-api/internal/domain"
	"time"

	"gorm.io/gorm"
)

type PriceHistoryRepository interface {
	FetchPriceHistories(ctx context.Context, symbol string, from, to time.Time) ([]domain.PriceHistory, error)
	SavePriceHistories(ctx context.Context, data []domain.PriceHistory) error
}

var _ PriceHistoryRepository = &PriceHistoryRepositoryImpl{}

type PriceHistoryRepositoryImpl struct {
	db *gorm.DB
}

func NewPriceHistoryRepositoryImpl(db *gorm.DB) *PriceHistoryRepositoryImpl {
	return &PriceHistoryRepositoryImpl{
		db: db,
	}
}

func (r *PriceHistoryRepositoryImpl) FetchPriceHistories(ctx context.Context, symbol string, from, to time.Time) ([]domain.PriceHistory, error) {
	var markets []domain.PriceHistory
	result := r.db.WithContext(ctx).Where("symbol = ? AND start_date >= ? AND end_date <= ?", symbol, from.Format(config.DATE_ONLY_FORMAT), to.Format(config.DATE_ONLY_FORMAT)).Find(&markets)
	if result.Error != nil {
		return nil, result.Error
	}

	return markets, nil
}

func (r *PriceHistoryRepositoryImpl) SavePriceHistories(ctx context.Context, data []domain.PriceHistory) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, market := range data {
			if err := tx.Create(&market).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
