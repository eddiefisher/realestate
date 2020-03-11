package htmltemplates

import (
	"html/template"

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
	}
}
