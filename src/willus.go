package main

import (
    "os"
    "fmt"
    "io/ioutil"
    "log"
    "strings"
    "net/http"
    // forecast "github.com/mlbright/forecast/v2"
    "html/template"
)

const DefaultConfigPath = "config/willus.cfg"

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
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

func loadTemplate(tmplName string) (tmpl *template.Template) {
    tmpl, _ = template.ParseFiles(tmplName)
    return
}

func main() {
    config, err := loadConfig(DefaultConfigPath)
    if err != nil {
        log.Fatal(err)
    }

    apiKeyPath := config.ApiKey.File
    key, err := apiKeyFromEnvOrPath(apiKeyPath)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("key: ", key)

    lat := config.Location.Lat
    long := config.Location.Long
    fmt.Printf("lat: %f long: %f\n", lat, long)

    cacheDir := config.Cache.Directory
    os.MkdirAll(cacheDir, 0700)
    fmt.Printf("cacheDir: %s \n", cacheDir)
    // var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

    // forecast, err := forecast.Get(key, strconv.FormatFloat(lat, 'f', 6, 64), 
    //         strconv.FormatFloat(long, 'f', 6, 64), "now", "us")
    // if err != nil {
    //     log.Fatal(err)
    // }

    // fmt.Println("Timezone", forecast.Timezone)
    // fmt.Println("Summary", forecast.Currently.Summary)
    // fmt.Println("Humidity", forecast.Currently.Humidity)
    // fmt.Println("Temperature", forecast.Currently.Temperature)
    // // fmt.Println(forecast.Flags.Units)
    // fmt.Println("Wind Speed", forecast.Currently.WindSpeed)
}

// func main() {
//     http.HandleFunc("/forecast/", handler)
//     http.ListenAndServe(":8080", nil)
// }
