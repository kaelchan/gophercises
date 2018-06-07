package urlshort

import (
	"errors"
	"net/http"

	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	// Implement this...
	return func(response_writer http.ResponseWriter, request *http.Request) {
		url := (*request.URL).Path
		// if hit short url, redirect
		// if not hit, fall back
		if mappedURL, ok := pathsToUrls[url]; ok {
			http.Redirect(response_writer, request, mappedURL, 301)
		} else {
			fallback.ServeHTTP(response_writer, request)
		}
	}
}

func parseYAML(yml []byte) (map[string]string, error) {
	unmarshedYAML := make([]map[string]string, 0)
	err := yaml.Unmarshal(yml, &unmarshedYAML)
	if err != nil {
		return nil, err
	}

	ret := make(map[string]string)
	for _, pathMapping := range unmarshedYAML {
		path, pathOK := pathMapping["path"]
		url, urlOK := pathMapping["url"]
		if !pathOK || !urlOK {
			err = errors.New("Error format in yaml")
			break
		}
		ret[path] = url
	}
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	// TODO: Implement this...
	pathsToUrls, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}
	return MapHandler(pathsToUrls, fallback), nil
}
