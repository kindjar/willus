package main

import (
    "log"
    "path/filepath"
    "io/ioutil"
    "encoding/json"
    "time"
    forecastio "github.com/mlbright/forecast/v2"
)

type WeatherCache struct {
    Directory string
    CacheTimeoutMinutes float64
    logger *log.Logger
}

const singletonCacheFile = "weather_cache"
const timeFormat = "2006-01-02 15:04:05"
var cacheSemaphore = make(chan int, 1)

func NewWeatherCache(cacheDir string, cacheTimeout float64, logger *log.Logger) (*WeatherCache) {
    return &WeatherCache{
        Directory: cacheDir,
        CacheTimeoutMinutes: cacheTimeout,
        logger: logger,
    }
}

func (cache *WeatherCache) pathForData(lat float64, long float64) (path string) {
    path = filepath.Join(cache.Directory, singletonCacheFile)
    // cache.logger.Println("Using cache file:", path)
    return
}

func (cache *WeatherCache) Get(lat float64, long float64, allowStale bool) (forecast *forecastio.Forecast, isStale bool) {
    isStale = true
    cacheSemaphore <- 1
    cachedBytes, err := ioutil.ReadFile(cache.pathForData(lat, long))
    <-cacheSemaphore
    if err == nil {
        var unmarshalledForecast forecastio.Forecast
        err = json.Unmarshal(cachedBytes, &unmarshalledForecast)
        if err == nil {
            unixTime := time.Unix(int64(unmarshalledForecast.Currently.Time), 0)
            timeAgo := time.Since(unixTime)
            if timeAgo.Minutes() < cache.CacheTimeoutMinutes {
                cache.logger.Printf("Using cached data from %s\n", unixTime.Format(timeFormat))
                forecast = &unmarshalledForecast
                isStale = false
            } else {
                cache.logger.Printf("Cache is stale (%s)\n", unixTime.Format(timeFormat))
                if allowStale {
                    forecast = &unmarshalledForecast
                }
            }
        } else {
            cache.logger.Fatalln("Invalid cache:", err)
        }
    } else {
        cache.logger.Println("Unable to read cache:", err)
    }
    return
}

func (cache *WeatherCache) Put(lat float64, long float64, forecast *forecastio.Forecast) (err error) {
    path := cache.pathForData(lat, long)
    jsonBytes, err := json.MarshalIndent(forecast, "", "  ")
    if err == nil {
        cacheSemaphore <- 1
        err = ioutil.WriteFile(path, jsonBytes, 0600)
        <-cacheSemaphore
    }
    return
}
