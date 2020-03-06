package main

import (
	"github.com/eddiefisher/realestate/pkg/mongo"
	"github.com/eddiefisher/realestate/pkg/parser"
)

func main() {
	conf := parser.NewConfig("configs/realestate.toml")

	mongoClient, ctx := mongo.Connect()
	defer mongoClient.Disconnect(ctx)

	lands := parser.Run(conf)
	parser.Save(lands, mongoClient)
}
