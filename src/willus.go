package main

import (
    "os"
    "fmt"
    "io/ioutil"
    "log"
    "strings"
    "strconv"
    "net/http"
    forecast "github.com/mlbright/forecast/v2"
)

const DefaultApiKeyPath = "config/secrets/api_key.txt"
const DefaultLocationPath = "config/location.txt"

type WillusError struct {
    Message string
}

func (e WillusError) Error() (error string) {
    return e.Message
}

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
    if path == "" {
        path = DefaultApiKeyPath
    }
    keybytes, err := ioutil.ReadFile(path)
    if err != nil {
        return "", &WillusError{fmt.Sprintf(
                "Unable to read API key from path %s: %s", path, err.Error())}
    }
    key = string(keybytes)
    return strings.TrimSpace(key), nil
}

func getLocationFromPath(path string) (lat float64, long float64, err error) {
    locationbytes, err := ioutil.ReadFile(path)
    if err != nil {
        return 0, 0, &WillusError{fmt.Sprintf(
                "Unable to read location from path %s: %s", path, err.Error())}
    }
    location := strings.Split(string(locationbytes), ",")
    lat, err = strconv.ParseFloat(strings.TrimSpace(location[0]), 64)
    if err != nil {
        return 0, 0, &WillusError{fmt.Sprintf(
                "Unable to read latitude from location %s (path: %s): %s", 
                        string(locationbytes), path, err.Error())}
    }
    long, err = strconv.ParseFloat(strings.TrimSpace(location[1]), 64)
    if err != nil {
        return 0, 0, &WillusError{fmt.Sprintf(
                "Unable to read longitude from location %s (path: %s): %s", 
                        string(locationbytes), path, err.Error())}
    }
    return lat, long, nil
}

func main() {
    var err error
    key, err := apiKeyFromEnvOrPath(DefaultApiKeyPath)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("key: ", key)
    lat, long, err := getLocationFromPath(DefaultLocationPath)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("lat: %f long: %f\n", lat, long)

    forecast, err := forecast.Get(key, strconv.FormatFloat(lat, 'f', 6, 64), 
            strconv.FormatFloat(long, 'f', 6, 64), "now", "us")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(forecast.Timezone)
    fmt.Println(forecast.Currently.Summary)
    fmt.Println(forecast.Currently.Humidity)
    fmt.Println(forecast.Currently.Temperature)
    fmt.Println(forecast.Flags.Units)
    fmt.Println(forecast.Currently.WindSpeed)
}

// func main() {
//     http.HandleFunc("/", handler)
//     http.ListenAndServe(":8080", nil)
// }
