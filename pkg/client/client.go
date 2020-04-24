package client

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/eddiefisher/realestate/pkg/htmltemplates"
	"github.com/eddiefisher/realestate/pkg/middleware"
	"github.com/eddiefisher/realestate/pkg/parser"
	"github.com/sirupsen/logrus"
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
		logrus.Fatal("$PORT must be set")
	}
	mongodb = db
	rootdir := "web/downloads/images/"
	http.Handle("/static/img/", http.StripPrefix("/static/img/",
		http.FileServer(http.Dir(path.Join(rootdir)))))
	indexHandler := http.HandlerFunc(IndexPage)
	http.Handle("/", middleware.BasicAuthMiddleware(indexHandler))

	logrus.Fatal(http.ListenAndServe(":"+port, nil))
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
		logrus.Errorf("mongo error: %s", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	dir, _ := os.Getwd()
	tmpl, err := template.New("layout.html").Funcs(htmltemplates.FuncMap()).ParseFiles(
		fmt.Sprintf("%s/web/templates/layout.html", dir),
		fmt.Sprintf("%s/web/templates/lands/index.html", dir),
		fmt.Sprintf("%s/web/templates/lands/search.html", dir),
		fmt.Sprintf("%s/web/templates/lands/sort.html", dir),
		fmt.Sprintf("%s/web/templates/lands/item.html", dir),
		fmt.Sprintf("%s/web/templates/pagination/pagination.html", dir),
		fmt.Sprintf("%s/web/templates/metrics/yandex.html", dir),
	)
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	l := Layout{
		Title:      "Realestate",
		Lands:      lands,
		Pagination: pagination,
	}
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
		Limit: 50,
		Max:   9,
	}
	if len(page) == 0 {
		page = "0"
	}
	current, err := strconv.Atoi(page)
	if err != nil {
		logrus.Errorf("-=page must be int=- page: %d, offset: %d, total: %d", pagination.Current, pagination.Offset, pagination.Total)
		http.Error(w, "", http.StatusBadRequest)
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
		logrus.Errorf("-=big page=- page: %d, offset: %d, total: %d", pagination.Current, pagination.Offset, pagination.Total)
		http.Error(w, "", http.StatusBadRequest)
		return Pagination{}, err
	}
	return pagination, nil
}

func landsPage(p Pagination) (parser.Lands, error) {
	collection := mongodb.Database("realestate").Collection("lands")
	ops := options.Find().SetLimit(int64(p.Limit)).SetSkip(int64(p.Offset)).SetSort(bson.D{{"addedat", -1}})
	cur, err := collection.Find(context.Background(), bson.M{}, ops)
	if err != nil {
		logrus.Errorf("find error: %s", err.Error())
		return nil, err
	}
	defer cur.Close(context.Background())

	var results parser.Lands
	for cur.Next(context.Background()) {
		var elem parser.Land
		err := cur.Decode(&elem)
		if err != nil {
			logrus.Errorf("parse element: %s", err.Error())
			return nil, err
		}

		results = append(results, elem)
	}

	if err := cur.Err(); err != nil {
		logrus.Errorf("cursor: %s", err.Error())
		return nil, err
	}

	return results, nil
}
