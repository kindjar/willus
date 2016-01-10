package willus

import (
	gcfg "code.google.com/p/gcfg"
	"fmt"
)

type Config struct {
	Webserver struct {
		Port int64
	}
	ApiKey struct {
		File string
	}
	Location struct {
		Lat  float64
		Long float64
	}
	Templates struct {
		Directory string
	}
	Cache struct {
		Directory           string
		CacheTimeoutMinutes float64
	}
}

func LoadConfig(path string) (cfg Config, err error) {
	err = gcfg.ReadFileInto(&cfg, path)
	if err != nil {
		err = &WillusError{fmt.Sprintf(
			"Unable to read config from %s: %s", path, err.Error())}
	}
	return
}
