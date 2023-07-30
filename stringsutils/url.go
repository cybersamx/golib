package stringsutils

import (
	"fmt"
	"net/url"
	"path"
	"path/filepath"
	"strings"
)

// MaskPassword masks the password component of an url.
func MaskPassword(text string) string {
	u, err := url.Parse(text)
	if err != nil {
		return text
	}

	if u.User != nil {
		pwd, set := u.User.Password()
		if set {
			return strings.Replace(u.String(), fmt.Sprintf("%s@", pwd), "***@", 1)
		}
	}

	return u.String()
}

// MapParams substitute named parameters embedded in the path with a map comprising
// the parameter keys and values. The function with `path` of "/users/:name" and `params`
// of `{"name": "lee"}` will return "/users/lee".
func MapParams(path string, params map[string]string) string {
	var (
		found bool
		sb    strings.Builder
		param string
		start int
	)

	size := len(path)

	for i, ch := range path {
		if ch == ':' && !found {
			start = i
			found = true
		}

		if ch == '/' && found {
			param = path[start+1 : i]
			value, ok := params[param]

			if ok {
				sb.WriteString(value)
				sb.WriteRune('/')
			} else {
				sb.WriteRune(':')
				sb.WriteString(param)
			}

			found = false

			continue
		}

		if i == size-1 && found {
			param = path[start+1:]
			value, ok := params[param]

			if ok {
				sb.WriteString(value)
			} else {
				sb.WriteRune(':')
				sb.WriteString(param)
			}

			found = false

			continue
		}

		if !found {
			sb.WriteRune(ch)
		}
	}

	return sb.String()
}

// TrimFileExtension remove the file extension from the filePath.
func TrimFileExtension(filePath string) string {
	ext := filepath.Ext(filePath)
	return filePath[0 : len(filePath)-len(ext)]
}

func GetPathFromURL(siteURL string) string {
	if siteURL == "" {
		return "/"
	}

	uri, err := url.Parse(siteURL)
	if err != nil {
		return "/"
	}

	if uri.Path == "" {
		return "/"
	}

	return path.Clean(uri.Path)
}
