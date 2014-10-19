package "main"

import (
  "html/template"
)

const longDateFormat = "2006-01-02 15:04:05"
const shortDateFormat = "Jan 2, 2006 3:04pm"

func DateTimeAsString(timestamp int64, format string) (formatted string) {
  unixTime := time.Unix(int64(timestamp), 0)
  switch format {
    case "short": formatted = unixTime.Format(shortDateFormat)
    case "long": formatted = unixTime.Format(longDateFormat)
  }
  return
}

func CommonTemplateHelpers() (functionMap FuncMap) {
  functionMap['DateTimeAsString'] = DateTimeAsString
  return
}
