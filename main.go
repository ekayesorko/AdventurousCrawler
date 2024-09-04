package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"strconv"
	"strings"
)

type Data struct {
	Title    string  `json:"title"`
	Category string  `json:"category"`
	Price    float64 `json:"price"`
	Sizes    []Size  `json:"sizes"`
}

type Size struct {
	Name         string `json:"size-name"`
	Availability bool   `json:"availability"`
}

func main() {
	c := colly.NewCollector()
	data := Data{}
	c.OnHTML("#__next > div > div.page.css-1umyepy > div.contentsWrapper > main > div > div > div.articleOverview.test-articleOverview > div.articlePurchaseBox.css-gxzada > div.articleInformation.css-itvqo3 > div.articleNameHeader.css-t1z1wj > h1", func(e *colly.HTMLElement) {
		fmt.Println(e.Text)
		data.Title = e.Text
	})

	c.OnHTML("#__next > div > div.page.css-1umyepy > div.contentsWrapper > main > div > div > div.articleOverview.test-articleOverview > div.articlePurchaseBox.css-gxzada > div.articleInformation.css-itvqo3 > div.articleNameHeader.css-t1z1wj > a", func(element *colly.HTMLElement) {
		fmt.Println(element.Text)
		data.Category = element.Text
	})

	//c.OnHTML("#__next > div > div.page.css-1umyepy > div.contentsWrapper > main > div > div > div.articleOverview.test-articleOverview > div.articlePurchaseBox.css-gxzada > div.articleInformation.css-itvqo3 > div.articlePrice.test-articlePrice.css-1apqb46 > p.price-text.test-price-text.mod-flat > span")
	c.OnHTML("#__next > div > div.page.css-1umyepy > div.contentsWrapper > main > div > div > div.articleOverview.test-articleOverview > div.articlePurchaseBox.css-gxzada > div.articleInformation.css-itvqo3 > div.articlePrice.test-articlePrice.css-1apqb46 > p.price-text.test-price-text.mod-flat > span", func(element *colly.HTMLElement) {
		priceText := element.Text
		priceText = strings.ReplaceAll(priceText, ",", "")
		price, err := strconv.ParseFloat(priceText, 64)
		if err != nil {
			log.Fatal(err)
		}
		data.Price = price
	})
	c.OnHTML("#__next > div > div.page.css-1umyepy > div.contentsWrapper > main > div > div > div.articleOverview.test-articleOverview > div.articlePurchaseBox.css-gxzada > div.addToCartForm.css-1127cv2 > div.inputSelects.clearfix > div.test-sizeSelector.css-539bvd > ul", func(element *colly.HTMLElement) {
		sizes := make([]Size, 0)
		element.ForEach("button", func(i int, element *colly.HTMLElement) {
			sizes = append(sizes, Size{
				Name:         element.ChildText("font > font"),
				Availability: !strings.Contains(element.Attr("class"), "disable"),
			})
		})
		fmt.Println(len(sizes))
		data.Sizes = sizes
	})
	err := c.Visit("http://localhost:5030/adidas-test.html")

	if err != nil {
		log.Fatal(err)
	}
}
