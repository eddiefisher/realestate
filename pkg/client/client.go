package client

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/eddiefisher/realestate/pkg/htmltemplates"
	"github.com/eddiefisher/realestate/pkg/middleware"
	"github.com/eddiefisher/realestate/pkg/paginator"
	"github.com/eddiefisher/realestate/pkg/paginator/adapter"
	"github.com/eddiefisher/realestate/pkg/paginator/view"
	"github.com/eddiefisher/realestate/pkg/parser"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

// Pagination ...
type Pagination struct {
	Pages   []int
	Next    int
	Prev    int
	Current int
	Last    int
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
	page, _ := strconv.Atoi(r.FormValue("page"))
	lands, pagination, err := landsPage(page)
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

	view := view.New(&pagination)

	tmpl.Execute(w, Layout{
		Title: "Realestate",
		Lands: lands,
		Pagination: Pagination{
			Pages:   view.Pages(),   // [2 3 4 5 6 7 8 9 10 11]
			Next:    view.Next(),    // 8
			Prev:    view.Prev(),    // 6
			Current: view.Current(), // 7
			Last:    view.Last(),
		},
	})
}

func landsPage(page int) (parser.Lands, paginator.Paginator, error) {
	collection := mongodb.Database("realestate").Collection("lands")
	p := paginator.New(adapter.NewMongoAdapter(collection), 13)
	p.SetPage(page)

	var lands parser.Lands
	if err := p.Results(&lands); err != nil {
		logrus.Println(err)
	}

	logrus.Println(p.HasNext(), p.HasPrev(), p.HasPages(), p.PageNums())

	return lands, p, nil
}
