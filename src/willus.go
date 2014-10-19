package main

import (
    "os"
    "fmt"
    "io/ioutil"
    "log"
    "strings"
    "strconv"
    "net/http"
    "encoding/json"
    htmlTemplate "html/template"
)

const DefaultConfigPath = "config/willus.cfg"
var config *Config
var logger *log.Logger
var forecaster *ForecastService
var defaultLatitude float64
var defaultLongitude float64

func mainHandler(w http.ResponseWriter, r *http.Request) {
    myLat := valueAsFloatWithDefault(r.FormValue("lat"), defaultLatitude)
    myLong := valueAsFloatWithDefault(r.FormValue("long"), defaultLongitude)
    templateName := "main"
    page := htmlTemplate.Must(
            htmlTemplate.New(templateName).
                    Funcs(CommonTemplateHelpers()).
                    ParseGlob("templates/*.tmpl"))
    // page, err := htmlTemplate.ParseFiles("templates/index.html.tmpl")
    //                 // Funcs(CommonTemplateHelpers()).
    //                 // ParseFiles("templates/index.html.tmpl")
    forecast, isStale := forecaster.Get(myLat, myLong, true)
    if isStale {
        logger.Println("Forecast is stale.")
    }
    err := page.ExecuteTemplate(w, "index.html.tmpl", forecast)
    // err = page.Execute(w, forecast)
    if err != nil {
        logger.Println("Error executing template: ", err)
    }
}

func valueAsFloatWithDefault(value string, defaultFloat float64) (float64) {
    floatVal, err := strconv.ParseFloat(value, 64)
    if err != nil {
        return defaultFloat
    } else {
        return floatVal
    }
}

func apiKeyFromEnvOrPath(path string) (key string, err error) {
    key = apiKeyFromEnvironment()
    if key != "" {
        return key, nil
    } else {
        return apiKeyFromPath(path)
    }
}

func apiKeyFromEnvironment() (key string) {
    return os.Getenv("FORECAST_API_KEY")
}

func apiKeyFromPath(path string) (key string, err error) {
    keybytes, err := ioutil.ReadFile(path)
    if err != nil {
        return "", &WillusError{fmt.Sprintf(
                "Unable to read API key from path %s: %s", path, err.Error())}
    }
    key = string(keybytes)
    return strings.TrimSpace(key), nil
}

func loadTemplate(tmplName string) (tmpl *htmlTemplate.Template) {
    tmpl, _ = htmlTemplate.ParseFiles(tmplName)
    return
}

func setupForecaster(config Config, apiKey string, logger *log.Logger) (forecaster *ForecastService) {
    cacheDir := config.Cache.Directory
    os.MkdirAll(cacheDir, 0700)
    logger.Printf("cacheDir: %s \n", cacheDir)

    jsonBytes, _ := json.MarshalIndent(config, "", "  ")
    logger.Println(string(jsonBytes))

    cache := NewWeatherCache(config.Cache.Directory,
            config.Cache.CacheTimeoutMinutes, logger)

    return NewForecastService(apiKey, cache, logger)
}

func main() {
    logger = log.New(os.Stderr, "", log.Lshortfile | log.LstdFlags)
    config, err := loadConfig(DefaultConfigPath)
    if err != nil {
        logger.Fatal(err)
    }

    apiKeyPath := config.ApiKey.File
    key, err := apiKeyFromEnvOrPath(apiKeyPath)
    if err != nil {
        logger.Fatal(err)
    }
    logger.Println("key: ", key)

    forecaster = setupForecaster(config, key, logger)

    defaultLatitude = config.Location.Lat
    defaultLongitude = config.Location.Long
    logger.Printf("lat: %f long: %f\n", defaultLatitude, defaultLongitude)

    forecast, _ := forecaster.Get(defaultLatitude, defaultLongitude, true)

    // var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

    fmt.Println("Timezone", forecast.Timezone)
    fmt.Println("Summary", forecast.Currently.Summary)
    fmt.Println("Humidity", forecast.Currently.Humidity)
    fmt.Println("Temperature", forecast.Currently.Temperature)
    // fmt.Println(forecast.Flags.Units)
    fmt.Println("Wind Speed", forecast.Currently.WindSpeed)

    http.HandleFunc("/", mainHandler)
    http.ListenAndServe(":8080", nil)

    forecaster.WaitForPendingForecasts()
}

// func main() {
//     http.HandleFunc("/forecast/", handler)
//     http.ListenAndServe(":8080", nil)
// }
