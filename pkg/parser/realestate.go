package parser

import (
	"fmt"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

// Realestate ...
type Realestate struct {
	Name      string
	URL       string
	URLPrefix string
}

// RealestateParse ...
type RealestateParse func(r Realestate, c *colly.Collector) Lands

// ParseList ...
func (r Realestate) ParseList(parse RealestateParse) Lands {
	c := colly.NewCollector(
		colly.CacheDir(fmt.Sprintf("./cache/%s/%s", r.Name, CachDate())),
	)
	extensions.RandomUserAgent(c)
	return parse(r, c)
}
