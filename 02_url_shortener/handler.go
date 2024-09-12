package main

import (
	"fmt"
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
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		longUrl, keyExists := pathsToUrls[r.URL.Path]

		if keyExists {
			http.Redirect(w, r, longUrl, http.StatusFound)
		} else {
			fallback.ServeHTTP(w, r)
		}
	})
}

type ShortUrl struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var urlList []ShortUrl
	urlMap := make(map[string]string, 0)

	err := yaml.Unmarshal(yml, &urlList)
	if err != nil {
		return nil, err
	}

	for _, url := range urlList {
		urlMap[url.Path] = url.Url
	}

	fmt.Println(urlList, urlMap)

	return MapHandler(urlMap, fallback), nil
}
