package parser

import (
	"net/url"
	"strings"

	"github.com/sirupsen/logrus"
)

// Images ...
type Images []Image

// Image ...
type Image struct {
	Name string
	Path string
	URL  string
}

// BuildFileName ...
func (i Image) BuildFileName() string {
	fileURL, err := url.Parse(i.URL)
	if err != nil {
		logrus.Error(err)
	}

	segments := strings.Split(fileURL.Path, "/")
	return segments[len(segments)-1]
}
