package templating

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// asConstValue formats val as a Golang constant.
func asConstValue(val any) string {
	switch t := val.(type) {
	case time.Duration:
		return fmt.Sprintf("time.Duration(%d)", t)
	case float64:
		return fmt.Sprintf("float64(%v)", t)
	case *url.URL:
		return fmt.Sprintf(`envPkg.GetURL("%s")`, t)
	case []*url.URL:
		return fmt.Sprintf(`envPkg.GetURLSlice([]string{%s})`, formatURLSlice(t))
	default:
		return fmt.Sprintf("%#v", val)
	}
}

// snakeToCamelLookupTable is a map of words to their respective value in camel case.
// It is required as simple capitalization is not enough for some common terms.
var snakeToCamelLookupTable = map[string]string{
	"dsn":  "DSN",
	"sql":  "SQL",
	"db":   "DB",
	"sasl": "SASL",
	"http": "HTTP",
	"grpc": "GRPC",
	"api":  "API",
	"id":   "ID",
	"ttl":  "TTL",
}

// snakeToCamel converts snake case string to camel case.
func snakeToCamel(s string) string {
	b := strings.Builder{}
	for part := range strings.SplitSeq(s, "_") {
		if repl, ok := snakeToCamelLookupTable[part]; ok {
			b.WriteString(repl)
		} else {
			b.WriteString(cases.Title(language.Und).String(part))
		}
	}

	return b.String()
}

// formatVarListInComment generates the following string for use in method comments:
//
//	//   - {prefix}_names[0]
//	//   - {prefix}_names[1]
//	...
//	//   - names[0]
//	//   - names[1]
//	...
func formatEnvListInComment(names []string) string {
	b := strings.Builder{}
	l := len(names)

	fullEnvList := make([]string, 2*l)
	for i, name := range names {
		fullEnvList[i] = "{envPrefix}_" + name
		fullEnvList[l+i] = name
	}

	for i, elem := range fullEnvList {
		b.WriteString("//   - ")
		b.WriteString(elem)
		if i < len(fullEnvList)-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

// canBeConst indicates if val can be declared as const.
func canBeConst(val any) bool {
	switch val.(type) {
	case []string:
		return false
	case *url.URL:
		return false
	case []*url.URL:
		return false
	default:
		return true
	}
}

// prettyValue returns pretty representation of object.
func prettyValue(val any) string {
	switch t := val.(type) {
	case time.Duration:
		return t.String()
	case *url.URL:
		return t.String()
	case []*url.URL:
		return formatURLSlice(t)
	default:
		return fmt.Sprintf("%#v", t)
	}
}

func formatURLSlice(urls []*url.URL) string {
	strs := make([]string, len(urls))
	for i, u := range urls {
		strs[i] = fmt.Sprintf("%q", u.String())
	}

	return strings.Join(strs, ", ")
}

// code wraps value in Markdown inline code backticks.
func code(v any) string {
	return "`" + fmt.Sprint(v) + "`"
}

// formatEnvNames formats env names as comma separated Markdown inline code
// entries joined with <br> so they fit into a single table cell.
func formatEnvNames(names []string) string {
	strs := make([]string, len(names))
	for i, name := range names {
		strs[i] = code(name)
	}

	return strings.Join(strs, "<br>")
}

// formatRules formats validation rules as Markdown table cell content.
// Strictness is reflected in the prefix; empty rules produce a dash.
func formatRules(rules []string, critical bool) string {
	if len(rules) == 0 {
		return "—"
	}

	prefix := "**Non-strict:** "
	if critical {
		prefix = "**Strict:** "
	}

	strs := make([]string, len(rules))
	for i, rule := range rules {
		strs[i] = code(rule)
	}

	return prefix + strings.Join(strs, ", ")
}

// escapeMarkdownCell escapes characters that would break a Markdown table cell.
func escapeMarkdownCell(s string) string {
	s = strings.ReplaceAll(s, "|", `\|`)

	return strings.ReplaceAll(s, "\n", "<br>")
}
