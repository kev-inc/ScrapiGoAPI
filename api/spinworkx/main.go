package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

func GetMain(w http.ResponseWriter, r *http.Request) {
	type Item struct {
		Name  string `json:"name"`
		Price string `json:"price"`
		Img   string `json:"img"`
		Href  string `json:"href"`
	}
	items := []Item{}
	category := r.URL.Query().Get("category")
	page := r.URL.Query().Get("page")

	col := colly.NewCollector()
	col.OnHTML("ul.products", func(e *colly.HTMLElement) {
		e.ForEach("li > div.entry-product", func(_ int, el *colly.HTMLElement) {
			href := el.ChildAttr("div.entry-featured > a", "href")
			hrefSplit := strings.Split(href, "/")
			item := Item{
				Name:  el.ChildText("div.entry-wrap > header > h3"),
				Price: el.ChildText("div.entry-wrap > header > span"),
				Img:   el.ChildAttr("div.entry-featured > a > img", "src"),
				Href:  hrefSplit[len(hrefSplit)-2],
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
	col.OnError(func(r *colly.Response, err error) {
		fmt.Fprintf(w, "Server responded with a "+strconv.Itoa(r.StatusCode))
	})
	if category == "" {
		if page == "" || page == "1" {
			col.Visit("https://www.spinworkx.com/main/")
		} else {
			col.Visit("https://www.spinworkx.com/main/page/" + page)
		}
	} else {
		if page == "" || page == "1" {
			col.Visit("https://www.spinworkx.com/main/product-category/" + category)
		} else {
			col.Visit("https://www.spinworkx.com/main/product-category/" + category + "/page/" + page)
		}
	}
}
