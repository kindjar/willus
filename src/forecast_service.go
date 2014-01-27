package main

import (
    "strconv"
    "log"
    forecastio "github.com/mlbright/forecast/v2"
)

type ForecastService struct {
    ApiKey string
    Cache *WeatherCache
    logger *log.Logger
}

func NewForecastService(key string, cache *WeatherCache, logger *log.Logger) (*ForecastService) {
    return &ForecastService{ApiKey: key, Cache: cache, logger: logger}
}

func (forecaster *ForecastService) Get(lat float64, long float64) (forecast *forecastio.Forecast) {
    if forecaster.Cache != nil {
        forecast = forecaster.Cache.Get(lat, long)
    }

    if forecast == nil {
        forecaster.logger.Println("No cache available, making remote request")

        var err error
        forecast, err = forecastio.Get(forecaster.ApiKey, 
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
    }

    return forecast
}