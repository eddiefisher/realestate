package parser

import (
	"fmt"

	"github.com/gocolly/colly"
)

// Realestate ...
type Realestate struct {
	Name      string
	Url       string
	UrlPrefix string
}

// RealestateParse ...
type RealestateParse func(r Realestate, c *colly.Collector) Lands

// ParseList ...
func (r Realestate) ParseList(parse RealestateParse) Lands {
	c := colly.NewCollector(
		colly.CacheDir(fmt.Sprintf("./cache/%s/%s", r.Name, CachDate())),
	)
	return parse(r, c)
}
