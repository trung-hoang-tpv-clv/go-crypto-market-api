package routers

import (
	"go-crypto-market-api/internal/config"
	"go-crypto-market-api/internal/infrastructure/cache"
	"go-crypto-market-api/internal/infrastructure/coingecko"
	"go-crypto-market-api/internal/infrastructure/repository"
	"go-crypto-market-api/internal/interfaces/handlers"
	"go-crypto-market-api/internal/usecases"

	"github.com/gin-gonic/gin"
)

func MarketRoutes(router *gin.Engine) {
	cache := cache.NewRedisCacheImpl(config.RedisClient)
	cgClient := coingecko.NewClient()
	repo := repository.NewPriceHistoryRepositoryImpl(config.GetDB())
	priceHistoryUsecase := usecases.NewPriceHistoryUsecase(repo, cache, cgClient)

	marketDataHandler := handlers.NewMarketHandler(priceHistoryUsecase)
	router.GET("/get_histories", marketDataHandler.GetPriceHistories)
}
