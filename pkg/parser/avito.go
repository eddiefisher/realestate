package parser

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
)

// ParseAvito ...
func ParseAvito(r Realestate, c *colly.Collector) (lands Lands) {
	c.OnHTML(".snippet-horizontal.item.item_table .item__line", func(e *colly.HTMLElement) {
		link := fmt.Sprintf("%s%s", r.URLPrefix, e.ChildAttr("a.snippet-link", "href"))
		lands = append(lands, Land{
			UID:    base64.StdEncoding.EncodeToString([]byte(link)),
			Name:   e.ChildText(".snippet-link"),
			Link:   link,
			Info:   fmt.Sprintf("%s, %s", e.ChildText(".item-address__string"), e.ChildText(".fld_gaz")),
			Price:  e.ChildText(".snippet-price"),
			Date:   e.ChildText(".snippet-date-info"),
			Images: images(getAllImages(e)),
			Source: "avito",
		})
	})
	c.Visit(r.URL)

	if len(lands) == 0 {
		logrus.Errorf("%s - no lands", r.Name)
	}

	return lands
}

func getAllImages(e *colly.HTMLElement) []string {
	var images []string
	for _, i := range extractSRC(e.ChildAttrs(".photo-wrapper img", "srcset")) {
		images = append(images, i)
	}
	for _, i := range extractSRC(e.ChildAttrs(".item-slider-image img", "srcset")) {
		images = append(images, i)
	}
	for _, i := range extractSRC(e.ChildAttrs(".item-slider-image img", "data-srcset")) {
		images = append(images, i)
	}
	return images
}

func images(ix []string) Images {
	images := Images{}

	if len(ix) == 0 {
		return Images{}
	}

	for _, i := range ix {
		images = append(images, Image{URL: i})
	}
	return images
}

func extractSRC(ix []string) (result []string) {
	for _, i := range ix {
		for _, s := range strings.Split(i, ",") {
			if strings.Contains(s, "1.5x") {
				s = strings.Replace(s, " 1.5x", "", 1)
				result = append(result, strings.Trim(s, " "))
			}
		}
	}
	return result
}
