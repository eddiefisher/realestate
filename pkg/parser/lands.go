package parser

import (
	"context"

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

	_, err := collection.InsertOne(context.Background(), l)
	if err != nil {
		return err
	}

	return nil
}
