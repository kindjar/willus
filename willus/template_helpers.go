package willus

import (
	"fmt"
	htmlTemplate "html/template"
	"os"
	"log"
	"time"
)

const longDateTimeFormat = "2006-01-02 15:04:05"
const shortDateTimeFormat = "Jan 2, 2006 3:04pm"
const shortDateFormat = "Jan 2"
const shortTimeFormat = "3:04pm"
const iso8601DateTimeFormat = "2006-01-02T15:04:05Z0700"

var logger = log.New(os.Stderr, "", log.Lshortfile|log.LstdFlags)

func FloatAsDateTime(timestamp float64, format string) (formatted string) {
	unixTime := time.Unix(int64(timestamp), 0)
	switch format {
	case "short":
		formatted = unixTime.Format(shortDateTimeFormat)
	case "long":
		formatted = unixTime.Format(longDateTimeFormat)
	case "shortDate":
		formatted = unixTime.Format(shortDateFormat)
	case "shortTime":
		formatted = unixTime.Format(shortTimeFormat)
	case "iso8601DateTime":
		formatted = unixTime.Format(iso8601DateTimeFormat)
	default:
		formatted = unixTime.Format(longDateTimeFormat)
	}
	return
}

func RoundToInteger(number float64) (formatted string) {
	return fmt.Sprintf("%2.f", number)
}

func FloatAsPercent(percent float64) (formatted string) {
	return fmt.Sprintf("%2.f%%", percent*100)
}

func FloatAsPrecipIntensityDescription(intensity float64, precipType string) (formatted string) {
	if len(precipType) == 0 {
		return // no-op to allow precipType parameter for future
	}

	if intensity > 0.4 {
		formatted = "heavy"
	} else if intensity > 0.1 {
		formatted = "moderate"
	} else if intensity > 0.17 {
		formatted = "light"
	} else if intensity > 0.002 {
		formatted = "very light"
	} else {
		formatted = "no"
	}
	return
}

func CommonTemplateHelpers() (functionMap htmlTemplate.FuncMap) {
	functionMap = htmlTemplate.FuncMap{
		"RoundToInteger":                    RoundToInteger,
		"FloatAsPrecipIntensityDescription": FloatAsPrecipIntensityDescription,
		"FloatAsPercent":                    FloatAsPercent,
		"FloatAsDateTime":                   FloatAsDateTime,
		"Log": func(message string) string {
			logger.Println(message)
			return message
		},
	}
	return
}
