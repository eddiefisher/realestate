package main

import (
	"github.com/eddiefisher/realestate/pkg/client"
	"github.com/eddiefisher/realestate/pkg/mongo"
)

func main() {
	// conf := parser.NewConfig("configs/client.toml")

	db, ctx := mongo.Connect()
	defer db.Disconnect(ctx)

	client.Start(db)
}
