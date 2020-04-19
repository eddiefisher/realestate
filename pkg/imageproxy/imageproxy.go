package imageproxy

import (
	"net/http"

	"github.com/eddiefisher/realestate/pkg/cloudinary"
	"github.com/sirupsen/logrus"
)

// Upload ...
func Upload(path string, folder string, publicID string) {
	cloud := cloudinary.Create("645956387348621", "i_4lsciU8xKdmt888PLc8auhqmU", "hrgs9dqgp")
	options := map[string]string{
		"public_id": publicID,
		"folder":    folder,
	}
	_, err := cloud.Upload(path, options)
	if err != nil {
		logrus.Error(err)
	}
}

// Get ...
func Get(folder string, publicID string) string {
	url := "https://res.cloudinary.com/hrgs9dqgp/image/upload/" + folder + "/" + publicID + ".jpg"
	r, _ := http.Get(url)
	if r.StatusCode != http.StatusOK {
		return "/static/img/0.gif"
	}
	return url
}
