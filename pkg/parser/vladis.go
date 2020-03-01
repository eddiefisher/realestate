package parser

import (
	"encoding/base64"
	"fmt"

	"github.com/gocolly/colly"
)

// ParseVladis ...
func ParseVladis(r Realestate, c *colly.Collector) (lands Lands) {
	parseHTMLElement := func(e *colly.HTMLElement) {
		link := fmt.Sprintf("%s%s", r.URLPrefix, e.ChildAttr(".fld_tit a", "href"))
		lands = append(lands, Land{
			UID:   base64.StdEncoding.EncodeToString([]byte(link)),
			Name:  e.ChildText(".fld_tit a"),
			Link:  link,
			Info:  fmt.Sprintf("%s, %s", e.ChildText(".fld_kmn"), e.ChildText(".fld_gaz")),
			Area:  e.ChildText(".fld_area"),
			Price: e.ChildText(".fld_price"),
			Date:  e.ChildText(".fld_date span"),
		})
	}
	c.OnHTML(".holder.obj-h table tbody tr", parseHTMLElement)
	c.Visit(r.URL)

	return lands
}
