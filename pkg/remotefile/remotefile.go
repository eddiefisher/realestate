package remotefile

import (
	"io"
	"net/http"
	"net/url"
	"os"
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

	// Create blank file
	file, err := r.createFile()
	if err != nil {
		return err
	}

	// Put content on file
	err = r.putFile(file, r.httpClient())
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

func (r *RemoteFile) putFile(file *os.File, client *http.Client) error {
	resp, err := client.Get(r.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)
	defer file.Close()
	if err != nil {
		return err
	}
	return nil
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

func (r *RemoteFile) httpClient() *http.Client {
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	return &client
}

func (r *RemoteFile) createFile() (*os.File, error) {
	path := r.Path + "/" + r.Prefix
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}
	file, err := os.Create(path + "/" + r.Name)
	if err != nil {
		return nil, err
	}
	return file, nil
}
