package main

import (
    "strconv"
    "log"
)

type ForecastService struct {
    ApiKey string
    Cache *WeatherCache
    logger *log.Logger
}

var forecastPendingSemaphore = make(chan int, 1)

func NewForecastService(key string, cache *WeatherCache, logger *log.Logger) (*ForecastService) {
    return &ForecastService{ApiKey: key, Cache: cache, logger: logger}
}

func (forecaster *ForecastService) Get(lat float64, long float64, allowStale bool) (forecast *Forecast, isStale bool) {
    isStale = false
    if forecaster.Cache != nil {
        forecast, isStale = forecaster.Cache.Get(lat, long, true)
    }

    if forecast == nil {
        forecaster.logger.Println("No cache available, making remote request")
        forecast = forecaster.fetchForecast(lat, long)
    } else if isStale && !allowStale {
        forecaster.logger.Println("Cache is stale and we need current data, making remote request")
        forecast = forecaster.fetchForecast(lat, long)
    } else if isStale {
        forecaster.logger.Println("Cache is stale but we don't need current data, making remote request in goroutine")
        forecastPendingSemaphore <- 1
        go func() {
            forecaster.fetchForecast(lat, long)
            <-forecastPendingSemaphore
        }()
    }

    return
}

func (forecaster *ForecastService) WaitForPendingForecasts() {
    forecastPendingSemaphore <- 1
}

func (forecaster *ForecastService) fetchForecast(lat float64, long float64) (forecast *Forecast) {
    var err error
    forecast, err = CallForecastAPI(
            forecaster.ApiKey,
            strconv.FormatFloat(lat, 'f', 6, 64),
            strconv.FormatFloat(long, 'f', 6, 64),
            "now", "us")

    if err != nil {
        forecaster.logger.Println(err)
    }

    if forecaster.Cache != nil {
        forecaster.logger.Println("Caching weather forecast")
        forecaster.Cache.Put(lat, long, forecast)
    }
    return
}
