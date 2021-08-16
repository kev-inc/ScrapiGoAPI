package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/gocolly/colly"
)

func GetMain(w http.ResponseWriter, r *http.Request) {
	type Item struct {
		Name	string  `json:"name"`
		Price 	string 	`json:"price"`
		Img    	string 	`json:"img"`
		Href	string	`json:"href"`
	}
	items := []Item{}
	category := r.URL.Query().Get("category")

	col := colly.NewCollector()
	col.OnHTML("ul.products", func(e *colly.HTMLElement) {
		e.ForEach("li > div.entry-product", func(_ int, el *colly.HTMLElement) {
			item := Item{
				Name: el.ChildText("div.entry-wrap > header > h3"),
				Price: el.ChildText("div.entry-wrap > header > span"),
				Img: el.ChildAttr("div.entry-featured > a > img", "src"),
				Href: el.ChildAttr("div.entry-featured > a", "href"),
			}
			items = append(items, item)
		})
	})
	col.OnScraped(func(r *colly.Response) {
		j, err := json.MarshalIndent(items, "", "\t")
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Fprintf(w, string(j))
	})
	if category == "" {
		col.Visit("https://www.spinworkx.com/main/")
	} else {
		col.Visit("https://www.spinworkx.com/main/product-category/" + category)
	}
}
