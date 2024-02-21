package handlers

import (
	"go-crypto-market-api/internal/config"
	"go-crypto-market-api/internal/usecases"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type MarketHandler struct {
	usecase *usecases.PriceHistoryUsecase
}

func NewMarketHandler(usecase *usecases.PriceHistoryUsecase) *MarketHandler {
	return &MarketHandler{
		usecase: usecase,
	}
}

func (h *MarketHandler) GetPriceHistories(c *gin.Context) {
	symbol := c.Query("symbol")
	startDateString := c.Query("startDate")
	endDateString := c.Query("endDate")

	if symbol == "" || startDateString == "" || endDateString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "symbol, startDate, and endDate parameters are required"})
		return
	}

	// Validate startDate
	startDate, err := time.Parse(config.DATE_ONLY_FORMAT, startDateString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "startDate must be in YYYY-MM-DD format"})
		return
	}

	// Validate endDate
	endDate, err := time.Parse(config.DATE_ONLY_FORMAT, endDateString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "endDate must be in YYYY-MM-DD format"})
		return
	}

	// Ensure startDate is not after endDate
	if startDate.After(endDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "startDate cannot be after endDate"})
		return
	}

	data, err := h.usecase.GetPriceHistories(c.Request.Context(), symbol, startDateString, endDateString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(data) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No market data found"})
		return
	}

	c.JSON(http.StatusOK, data)
}
