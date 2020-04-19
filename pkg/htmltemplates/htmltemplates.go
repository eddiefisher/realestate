package htmltemplates

import (
	"html/template"
	"path/filepath"

	"github.com/eddiefisher/realestate/pkg/imageproxy"
	"github.com/eddiefisher/realestate/pkg/parser"
)

// FuncMap ...
func FuncMap() template.FuncMap {
	return template.FuncMap{
		"getFirstImageName": func(images parser.Images) string {
			if len(images) == 0 {
				return "0.gif"
			}
			return images[0].BuildFileName()
		},
		"getFirstImageFromCloud": func(prefix string, images parser.Images) string {
			if len(images) == 0 {
				return "/static/img/0.gif"
			}
			var filename = images[0].BuildFileName()
			var extension = filepath.Ext(filename)
			var name = filename[0 : len(filename)-len(extension)]
			return imageproxy.Get(prefix, name)
		},
	}
}
