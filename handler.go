package urlshort

import (
	"net/http"

	"gopkg.in/yaml.v3"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		newURL, prs := pathsToUrls[path]
		if prs {
			http.Redirect(w, r, newURL, 301)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

// URLHolder contains unmarshaled url data from a given url-
// mapping document.
type URLHolder struct {
	Path string
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
	var redirects []URLHolder
	err := yaml.Unmarshal(yml, &redirects)
	fn := func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		for _, entry := range redirects {
			if entry.Path == path {
				http.Redirect(w, r, entry.URL, 301)
				return
			}
		}
		fallback.ServeHTTP(w, r)
	}
	return fn, err
}
