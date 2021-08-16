package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/gocolly/colly"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	type Images struct {
		Src       string	`json:"src"`
		Thumbnail string	`json:"thumbnail"`
	}
	type Product struct {
		Name   string   `json:"name"`
		Price  string   `json:"price"`
		Tags   []string `json:"tags"`
		Desc   string   `json:"desc"`
		Img    []Images `json:"img"`
		Colors []string `json:"colors"`
	}
	product := Product{}
	name := r.URL.Query().Get("name")

	col := colly.NewCollector()
	col.OnHTML("h1.product_title.entry-title", func(e *colly.HTMLElement) {
		product.Name = e.Text
	})
	col.OnHTML("p.price", func(e *colly.HTMLElement) {
		product.Price = e.Text
	})
	col.OnHTML("div.product_meta > span.posted_in", func(e *colly.HTMLElement) {
		product.Tags = make([]string, 0)
		e.ForEach("a", func(_ int, el *colly.HTMLElement) {
			product.Tags = append(product.Tags, el.Text)
		})
	})
	col.OnHTML("div.x-tab-content > div > p", func(e *colly.HTMLElement) {
		product.Desc = e.Text
	})
	col.OnHTML("figure.woocommerce-product-gallery__wrapper", func(e *colly.HTMLElement) {
		product.Img = make([]Images, 0)
		e.ForEach("div", func(_ int, el *colly.HTMLElement) {
			product.Img = append(product.Img, Images{
				Thumbnail: el.Attr("data-thumb"),
				Src: el.ChildAttr("a", "href"),
			})
		})
	})
	col.OnHTML("select#pa_color", func(e *colly.HTMLElement) {
		product.Colors = make([]string, 0)
		e.ForEach("option", func(_ int, el *colly.HTMLElement) {
			product.Colors = append(product.Colors, el.Text)
		})
	})
	col.OnScraped(func(r *colly.Response) {
		j, err := json.MarshalIndent(product, "", "\t")
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Fprintf(w, string(j))
	})
	col.Visit("https://www.spinworkx.com/main/product/" + name)
}
