package main

import (
    "time"
    url "net/url"
)

type DataPoint struct {
    strings map[string]*string
    floats map[string]*float64
    times map[string]*time.Time
}

func (p DataPoint) getString(key string) (str *string) {
    str = p.strings[key]
    return
}

func (p DataPoint) getFloat(key string) (num *float64, stdDev *float64) {
    num = p.floats[key]
    stdDev = p.floats[key + "Error"]
    return
}

func (p DataPoint) getTime(key string) (time *time.Time) {
    time = p.times[key]
    return
}

// type DataPoint {
//     Time                        Time
//     Summary                     string
//     Icon                        string
//     Sunrise                     *Time
//     Sunset                      *Time
//     PrecipIntensity             float64
//     PrecipIntensityMax          *float64
//     PrecipIntensityMaxTime      *Time
//     PrecipProbability           float64
//     PrecipType                  string
//     PrecipAccumulation          *float64
//     Temperature                 *float64
//     TemperatureMin              *float64
//     TemperatureMinTime          *Time
//     TemperatureMax              *float64
//     TemperatureMaxTime          *Time
//     ApparentTemperature         *float64
//     ApparentTemperatureMin      *float64
//     ApparentTemperatureMinTime  *Time
//     ApparentTemperatureMax      *float64
//     ApparentTemperatureMaxTime  *Time
//     DewPoint                    float64
//     WindSpeed                   float64
//     WindBearing                 *int
//     CloudCover                  float64
//     Humidity                    float64
//     Pressure                    float64
//     Visibility                  float64
//     Ozone                       float64
// }

type DataBlock struct {
    Summary string
    Icon    string
    Data    []DataPoint
}

type Alert struct {
    Title       string
    Expires     time.Time
    Description string
    Url         url.URL

}

type Flags struct {
    DarkSkyUnavailable  *bool
    DarkSkyStations     *[]string
    DatapointStations   *[]string
    IsdStations         *[]string
    LampStations        *[]string
    MetarStation        *[]string
    MetnoLicense        *string
    Sources             *[]string
    Units               *string
}

type Forecast struct {
    Latitude    float64
    Longitude   float64
    Timezone    string
    TzOffset    int
    Currently   DataPoint
    Minutely    DataBlock
    Hourly      DataBlock
    Daily       DataBlock
    Alerts      []Alert
    Flags       Flags
}