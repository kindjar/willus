package willus

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// URL example:  "https://api.forecast.io/forecast/APIKEY/LATITUDE,LONGITUDE,TIME?units=ca"
const (
	BASEURL = "https://api.darksky.net/forecast"
)

type Flags struct {
	DarkSkyUnavailable string
	DarkSkyStations    []string
	DataPointStations  []string
	ISDStations        []string
	LAMPStations       []string
	METARStations      []string
	METNOLicense       string
	Sources            []string
	Units              string
}

type DataPoint struct {
	Time                       float64
	Summary                    string
	Icon                       string
	SunriseTime                float64
	SunsetTime                 float64
	MoonPhase                  float64
	NearestStormDistance       float64
	NearestStormBearing        float64
	PrecipIntensity            float64
	PrecipIntensityMax         float64
	PrecipIntensityMaxTime     float64
	PrecipProbability          float64
	PrecipType                 string
	PrecipAccumulation         float64
	Temperature                float64
	TemperatureMin             float64
	TemperatureMinTime         float64
	TemperatureMax             float64
	TemperatureMaxTime         float64
	ApparentTemperature        float64
	ApparentTemperatureMin     float64
	ApparentTemperatureMinTime float64
	ApparentTemperatureMax     float64
	ApparentTemperatureMaxTime float64
	DewPoint                   float64
	WindSpeed                  float64
	WindBearing                float64
	CloudCover                 float64
	Humidity                   float64
	Pressure                   float64
	Visibility                 float64
	Ozone                      float64
}

type DataBlock struct {
	Summary string
	Icon    string
	Data    []DataPoint
}

type Alert struct {
	Title       string
	Expires     float64
	Description string
	URI         string
}

type Forecast struct {
	Latitude  float64
	Longitude float64
	Timezone  string
	Offset    float64
	Currently DataPoint
	Minutely  DataBlock
	Hourly    DataBlock
	Daily     DataBlock
	Alerts    []Alert
	Flags     Flags
}

type Units string

const (
	CA Units = "ca"
	SI Units = "si"
)

func CallForecastAPI(key string, lat string, long string, time string, units Units) (*Forecast, error) {
	coord := lat + "," + long

	var url string
	if time == "now" {
		url = BASEURL + "/" + key + "/" + coord + "?units=" + string(units)
	} else {
		url = BASEURL + "/" + key + "/" + coord + "," + time + "?units=" + string(units)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false}, // does not seem required any longer
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var f Forecast
	err = json.Unmarshal(body, &f)
	if err != nil {
		return nil, err
	}

	return &f, nil
}
