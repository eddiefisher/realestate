package adapter

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

// MongoAdapter mongo adapter to be passed to paginator constructor
type MongoAdapter struct {
	db *mongo.Collection
}

// NewMongoAdapter Mongo adapter constructor which receive the Mongo db query.
func NewMongoAdapter(db *mongo.Collection) MongoAdapter {
	return MongoAdapter{db: db}
}

// Nums returns the number of records
func (a MongoAdapter) Nums() int {
	var count int64
	count, err := a.db.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		logrus.Println(err)
	}

	return int(count)
}

// Slice stores into data argument a slice of the results.
// data must be a pointer to a slice of models.
func (a MongoAdapter) Slice(offset, length int, data interface{}) error {
	logrus.Println(offset, length)
	ops := options.Find().SetLimit(int64(length)).SetSkip(int64(offset)).SetSort(bson.M{"addedad": 1})
	cur, err := a.db.Find(context.Background(), bson.M{}, ops)
	_ = cur.All(context.Background(), data)

	return err
}
