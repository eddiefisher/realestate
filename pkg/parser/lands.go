package parser

import (
	"context"
	"time"

	"github.com/eddiefisher/realestate/pkg/remotefile"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

// Lands ...
type Lands []Land

// Land info about land
type Land struct {
	UID         string
	Name        string
	Link        string
	Info        string
	Area        string
	Price       string
	Description string
	Date        string
	AddedAt     time.Time
	Images      Images
	Source      string
}

// Append ...
func (lx Lands) Append(lands Lands) Lands {
	for _, land := range lands {
		lx = append(lx, land)
	}
	return lx
}

// Save ...
func (l Land) Save(client *mongo.Client) error {
	collection := client.Database("realestate").Collection("lands")
	loc, _ := time.LoadLocation("Europe/Moscow")
	l.AddedAt = time.Now().In(loc)
	l.DownloadImage()

	_, err := collection.InsertOne(context.Background(), l)
	if err != nil {
		// _, err = collection.ReplaceOne(context.Background(), bson.M{"uid": l.UID}, l)
		// if err != nil {
		// 	return err
		// }
		return err
	}

	return nil
}

// DownloadImage ...
func (l Land) DownloadImage() {
	if len(l.Images) == 0 {
		return
	}
	for _, image := range l.Images {
		err := remotefile.New(image.URL, l.Source).Download()
		if err != nil {
			logrus.Error(err)
		}
	}
}
