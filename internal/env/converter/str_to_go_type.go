package converter

import (
	"net/url"

	envErrors "github.com/kirill-zak/go-env/error"
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
		err = envErrors.ErrInvalidType
	}

	if err != nil {
		err = envErrors.ErrInvalidValue
	}

	return resT, resV, err
}
