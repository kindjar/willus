package main

import (
    "fmt"
    gcfg "code.google.com/p/gcfg"
)

type Config struct {
    ApiKey struct {
        File string
    }
    Location struct {
        Lat float64
        Long float64
    }
    Templates struct {
        Directory string
    }
}

func loadConfig(path string) (cfg Config, err error) {
    err = gcfg.ReadFileInto(&cfg, path)
    if err != nil {
        err = &WillusError{fmt.Sprintf(
                "Unable to read config from %s: %s", path, err.Error())}
    }
    return
}
