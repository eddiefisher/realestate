package parser

import (
	"fmt"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

// CachDate ...
func CachDate() string {
	t := time.Now()
	bod := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	date := fmt.Sprintf("%s%s%s", strconv.Itoa(t.Year()), t.Month().String(), strconv.Itoa(t.Day()))
	if t.Before(bod.Add(time.Hour * 12)) {
		return fmt.Sprintf("%s0", date)
	}
	return fmt.Sprintf("%s1", date)
}

// Run ...
func Run(conf Config) Lands {
	var lands Lands

	lands = lands.Append(Realestate{Name: "avito", URL: conf.Avito, URLPrefix: "https://www.avito.ru"}.ParseList(ParseAvito))
	lands = lands.Append(Realestate{Name: "vladis", URL: conf.Vladis, URLPrefix: "https://vladis.ru"}.ParseList(ParseVladis))
	return lands
}

// Save ...
func Save(lands Lands, client *mongo.Client) error {
	for _, land := range lands {
		err := land.Save(client)
		if err != nil {
			logrus.Println(err)
		}
	}
	return nil
}
