package parser

import (
	"encoding/base64"
	"fmt"

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
			Images: images(e.ChildAttrs(".item-slider-image img", "src")),
			Source: "avito",
		})
	})
	c.Visit(r.URL)

	if len(lands) == 0 {
		logrus.Errorf("%s - no lands", r.Name)
	}
	for _, land := range lands {
		logrus.Println(land.Name, land.Price)
	}

	return lands
}

func images(ix []string) Images {
	images := Images{}
	logrus.Info(images, len(images))
	if len(ix) == 0 {
		return Images{}
	}

	for _, i := range ix {
		images = append(images, Image{URL: i})
	}
	return images
}
