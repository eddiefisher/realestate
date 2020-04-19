package remotefile

import (
	"net/url"
	"path/filepath"
	"strings"

	"github.com/eddiefisher/realestate/pkg/imageproxy"
)

// RemoteFile ...
type RemoteFile struct {
	Name   string
	URL    string
	Prefix string
	Path   string
}

// New ...
func New(url string, prefix string) *RemoteFile {
	return &RemoteFile{
		URL:    url,
		Prefix: prefix,
		Path:   "web/downloads/images",
	}
}

// Download ...
func (r *RemoteFile) Download() error {
	// Build fileName from fullPath
	err := r.buildFileName()
	if err != nil {
		return err
	}

	r.cloudinaryUpload()
	return nil
}

// cloudinaryUpload
func (r *RemoteFile) cloudinaryUpload() {
	var extension = filepath.Ext(r.Name)
	var name = r.Name[0 : len(r.Name)-len(extension)]
	imageproxy.Upload(r.URL, r.Prefix, name)
}

func (r *RemoteFile) buildFileName() error {
	fileURL, err := url.Parse(r.URL)
	if err != nil {
		return err
	}

	segments := strings.Split(fileURL.Path, "/")
	r.Name = segments[len(segments)-1]

	return nil
}
