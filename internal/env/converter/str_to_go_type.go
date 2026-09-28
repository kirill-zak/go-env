package converter

import (
	"net/url"

	envErrors "github.com/kirill-zak/go-env/errors"
)

func AsGoType(typ, val string) (resT string, resV any, err error) {
	switch typ {
	case "duration":
		resT = "time.Duration"
		resV, err = strToDuration(val)

	case "string_slice", "strings":
		resT = "[]string"
		resV = strToSlice(val)

	case "", "string":
		resT = "string"
		resV = val

	case "float":
		resT = "float64"
		resV, err = strToFloat(val)

	case "int":
		resT = "int"
		resV, err = strToInt(val)

	case "bool":
		resT = "bool"
		resV = strToBool(val)

	case "url":
		resT = "*url.URL"
		resV, err = url.Parse(val)

	case "[]url":
		resT = "[]*url.URL"
		resV, err = strToURLSlice(val)

	default:
		return "", nil, &envErrors.InvalidTypeError{
			Name: "unknown",
			Type: typ,
		}
	}

	if err != nil {
		return resT, resV, &envErrors.InvalidValueError{
			Name:  "unknown",
			Value: val,
			Type:  typ,
		}
	}

	return resT, resV, err
}
