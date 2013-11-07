package main

import (
    "fmt"
    "path/filepath"
    "io/ioutil"
    "encoding/json"
    "time"
    forecastio "github.com/mlbright/forecast/v2"
)

type WeatherCache struct {
    Directory string
}

var singletonCacheFile = "weather_cache"

func NewWeatherCache(cacheDir string) (*WeatherCache) {
    return &WeatherCache{Directory: cacheDir}
}

func (cache *WeatherCache) pathForData(lat float64, long float64) (path string) {
    path = filepath.Join(cache.Directory, singletonCacheFile)
    // fmt.Println("Using cache file:", path)
    return
}

func (cache *WeatherCache) Get(lat float64, long float64) (forecast *forecastio.Forecast) {
    cachedBytes, err := ioutil.ReadFile(cache.pathForData(lat, long))
    if err == nil {
        var unmarshalledForecast forecastio.Forecast
        err = json.Unmarshal(cachedBytes, &unmarshalledForecast)
        if err == nil {
            unixTime := time.Unix(int64(unmarshalledForecast.Currently.Time), 0)
            timeAgo := time.Since(unixTime)
            if timeAgo.Hours() < 1.0 {
                fmt.Printf("Using cached data from %s\n", unixTime.Format("2006-01-02 15:04:05"))
                forecast = &unmarshalledForecast
            } else {
                fmt.Printf("Cache is stale (%s)\n", unixTime.Format("2006-01-02 15:04:05"))
            }
        } else {
            // fmt.Println("Invalid cache:", err)
        }
    } else {
        // fmt.Println("Unable to read cache:", err)
    }
    return
}

func (cache *WeatherCache) Put(lat float64, long float64, forecast *forecastio.Forecast) (err error) {
    path := cache.pathForData(lat, long)
    jsonBytes, err := json.MarshalIndent(forecast, "", "  ")
    if err == nil {
        err = ioutil.WriteFile(path, jsonBytes, 0600)
    }
    return
}