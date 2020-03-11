package parser

import (
	"encoding/base64"
	"fmt"
	"log"

	"github.com/gocolly/colly"
)

// ParseAvito ...
func ParseAvito(r Realestate, c *colly.Collector) (lands Lands) {
	c.OnHTML(".js-catalog_serp .snippet-horizontal.item.item_table", func(e *colly.HTMLElement) {
		link := fmt.Sprintf("%s%s", r.URLPrefix, e.ChildAttr("a.snippet-link", "href"))
		lands = append(lands, Land{
			UID:   base64.StdEncoding.EncodeToString([]byte(link)),
			Name:  e.ChildText(".snippet-link"),
			Link:  link,
			Info:  fmt.Sprintf("%s, %s", e.ChildText(".item-address__string"), e.ChildText(".fld_gaz")),
			Price: e.ChildText(".snippet-price"),
			Date:  e.ChildText(".snippet-date-info"),
		})
	})
	c.Visit(r.URL)

	if len(lands) == 0 {
		logrus.Errorf("%s - no lands", r.Name)
	}

	return lands
}
