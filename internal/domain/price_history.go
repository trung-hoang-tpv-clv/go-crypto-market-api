package domain

import (
	"time"
)

type PriceHistory struct {
	Symbol    string    `json:"symbol"`
	High      float64   `json:"high"`
	Low       float64   `json:"low"`
	Open      float64   `json:"open"`
	Close     float64   `json:"close"`
	Time      int64     `json:"time"`
	Change    float64   `json:"change"`
	StartDate time.Time `json:"start_date" gorm:"type:date"`
	EndDate   time.Time `json:"end_date" gorm:"type:date"`
}
