package main

import (
	"fmt"

	"github.com/eddiefisher/realestate/pkg/mongo"
	"github.com/eddiefisher/realestate/pkg/parser"
)

func main() {
	conf := parser.NewConfig("configs/realestate.toml")

	mongoClient, ctx := mongo.Connect()
	// mongoInit(mongoClient)
	defer mongoClient.Disconnect(ctx)

	lands := parser.Run(conf)
	parser.Save(lands, mongoClient)
}

func landToS(lands parser.Lands) []string {
	r := []string{}
	for _, l := range lands {
		r = append(r, fmt.Sprintf("%s, %s, %s, %s, %s\n", l.Name, l.Info, l.Area, l.Price, l.Link))
	}
	return r
}
