package client

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/eddiefisher/realestate/pkg/parser"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Pagination ...
type Pagination struct {
	Current     int
	Total       int
	Limit       int
	Offset      int
	Max         int // Max maximum pagination links
	MiddleTotal []int
}

// Layout ...
type Layout struct {
	Title      string
	Lands      parser.Lands
	Pagination Pagination
}

var mongodb *mongo.Client

// Start ...
func Start(db *mongo.Client) {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}
	mongodb = db
	http.HandleFunc("/", IndexPage)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// IndexPage ...
func IndexPage(w http.ResponseWriter, r *http.Request) {
	page := r.FormValue("page")
	pagination, err := getPage(page, w)
	if err != nil {
		return
	}
	lands, err := landsPage(pagination)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("ERROR! mongo error: %s", err)
		return
	}
	dir, _ := os.Getwd()
	tmpl, err := template.ParseFiles(
		fmt.Sprintf("%s/web/templates/layout.html", dir),
		fmt.Sprintf("%s/web/templates/lands/index.html", dir),
		fmt.Sprintf("%s/web/templates/pagination/pagination.html", dir),
	)
	if err != nil {
		log.Println(err.Error())
		return
	}
	l := Layout{
		Title:      "Realestate",
		Lands:      lands,
		Pagination: pagination,
	}
	log.Println(fmt.Sprintf("%v", l.Pagination.MiddleTotal))
	tmpl.Execute(w, l)
}

func totalPage() int {
	count, err := mongodb.Database("realestate").Collection("lands").CountDocuments(context.Background(), bson.M{})
	if err != nil {
		return 0
	}
	return int(count)
}

func getPage(page string, w http.ResponseWriter) (Pagination, error) {
	pagination := Pagination{
		Limit: 20,
		Max:   9,
	}
	if len(page) == 0 {
		page = "0"
	}
	current, err := strconv.Atoi(page)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("ERROR! -=page must be int=- page: %d, offset: %d, total: %d", pagination.Current, pagination.Offset, pagination.Total)
		return Pagination{}, err
	}

	pagination.Total = totalPage() / pagination.Limit
	pagination.Current = current
	pagination.Offset = pagination.Current * pagination.Limit
	if pagination.Total < pagination.Max {
		for i := 0; i != pagination.Max-1; i++ {
			pagination.MiddleTotal = append(pagination.MiddleTotal, i)
		}
	} else {
		if pagination.Current < pagination.Max/2 {
			for i := 0; i != pagination.Max; i++ {
				pagination.MiddleTotal = append(pagination.MiddleTotal, i)
			}
		} else {
			for i := pagination.Current - pagination.Max/2; i != pagination.Total; i++ {
				pagination.MiddleTotal = append(pagination.MiddleTotal, i)
			}
		}
	}

	if current > pagination.Total {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("ERROR! -=big page=- page: %d, offset: %d, total: %d", pagination.Current, pagination.Offset, pagination.Total)
		return Pagination{}, err
	}
	return pagination, nil
}

func landsPage(p Pagination) (parser.Lands, error) {
	collection := mongodb.Database("realestate").Collection("lands")
	ops := options.Find().SetLimit(int64(p.Limit)).SetSkip(int64(p.Offset))
	cur, err := collection.Find(context.Background(), bson.M{}, ops)
	if err != nil {
		log.Printf("Error: find error: %s", err.Error())
		return nil, err
	}
	defer cur.Close(context.Background())

	var results parser.Lands
	for cur.Next(context.Background()) {
		var elem parser.Land
		err := cur.Decode(&elem)
		if err != nil {
			log.Printf("Error: parse element: %s", err.Error())
			return nil, err
		}

		results = append(results, elem)
	}

	if err := cur.Err(); err != nil {
		log.Printf("Error cursor: %s", err.Error())
		return nil, err
	}

	return results, nil
}
