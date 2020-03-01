package mongo

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

// Connect user: realestate; pass: i0MZ9imH7aTt3aOh
func Connect() (*mongo.Client, context.Context) {
	atlasURL := "mongodb+srv://realestate:i0MZ9imH7aTt3aOh@cluster0-fw671.mongodb.net/test?retryWrites=true&w=majority"
	client, err := mongo.NewClient(options.Client().ApplyURI(atlasURL))
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	return client, ctx
}

// Init ...
func Init(client *mongo.Client) error {
	collection := client.Database("realestate").Collection("lands")
	_, err := collection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bsonx.Doc{{"uid", bsonx.Int32(1)}},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		return err
	}
	return nil
}
