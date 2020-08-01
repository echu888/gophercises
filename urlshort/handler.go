// implementation of Gophercise Exercise #2 
//   https://github.com/gophercises/urlshort/
//   
// primary features:
//	create http handlers
//	use maps
//	parse YAML

// use package main for easy testing 
//package urlshort
package main

import (
	"fmt"
	"net/http"

	"gopkg.in/yaml.v2"
)

// RedirectHandler uses a simple HTTP redirect, along with 
// an HTTP 301 (Permanently Moved) status
func RedirectHandler(w http.ResponseWriter, r *http.Request, newUrl string) {
	http.Redirect(w, r, newUrl, http.StatusMovedPermanently)
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if fullurl, ok := pathsToUrls[r.URL.Path]; ok {
			RedirectHandler(w,r,fullurl)
		} else {
			fallback.ServeHTTP(w,r)
		}
        }
}

// Converted from YAML into golang struct, using the
// "YAML to go" service available here: 
// https://yaml.to-go.online/
type YamlEntries []struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
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

	shortcuts := make(map[string]string)

	var yamlEntries YamlEntries
	err := yaml.Unmarshal([]byte(yml), &yamlEntries)
	if err == nil {
		for _,shortcut := range yamlEntries {
			shortcuts[shortcut.Path]=shortcut.URL
		}
	}
        return MapHandler(shortcuts, fallback), err
}
