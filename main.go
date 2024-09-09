package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"os"
	"strconv"
	"strings"
)

type Data struct {
	Title    string  `json:"title"`
	Category string  `json:"category"`
	Price    float64 `json:"price"`
	Link     string  `json:"link"`
	Sizes    []Size  `json:"sizes"`

	SizeTypeX []string   `json:"size_types_x"`
	SizeTypeY []string   `json:"size_types_y"`
	SizeInfo  [][]string `json:"size_info"`

	Reviews       []Review `json:"reviews"`
	OverallRating float64  `json:"overall_rating"`
	ImageUrls     []string `json:"image_urls"`
}

type Size struct {
	Name         string `json:"size-name"`
	Availability bool   `json:"availability"`
}

type Review struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	Date  string `json:"date"`
	Star  int    `json:"star"`
}

func main() {
	dataArr := make([]Data, 0)
	c := colly.NewCollector()
	var data Data

	c.OnHTML("#__next > div > div.page.css-1umyepy > div.contentsWrapper > main > div > div > div.articleOverview.test-articleOverview > div.articlePurchaseBox.css-gxzada > div.articleInformation.css-itvqo3 > div.articleNameHeader.css-t1z1wj > h1", func(e *colly.HTMLElement) {
		data.Title = e.Text
	})

	c.OnHTML("#__next > div > div.page.css-1umyepy > div.contentsWrapper > main > div > div > div.articleOverview.test-articleOverview > div.articlePurchaseBox.css-gxzada > div.articleInformation.css-itvqo3 > div.articleNameHeader.css-t1z1wj > a", func(element *colly.HTMLElement) {
		data.Category = element.Text
	})

	//c.OnHTML("#__next > div > div.page.css-1umyepy > div.contentsWrapper > main > div > div > div.articleOverview.test-articleOverview > div.articlePurchaseBox.css-gxzada > div.articleInformation.css-itvqo3 > div.articlePrice.test-articlePrice.css-1apqb46 > p.price-text.test-price-text.mod-flat > span")
	c.OnHTML("#__next > div > div.page.css-1umyepy > div.contentsWrapper > main > div > div > div.articleOverview.test-articleOverview > div.articlePurchaseBox.css-gxzada > div.articleInformation.css-itvqo3 > div.articlePrice.test-articlePrice.css-1apqb46 > p.price-text.test-price-text.mod-flat > span", func(element *colly.HTMLElement) {
		price, err := strconv.ParseFloat(strings.ReplaceAll(element.Text, ",", ""), 64)
		if err != nil {
			log.Fatal(err)
		}
		data.Price = price
	})

	c.OnHTML("#__next > div > div.page.css-1umyepy > div.contentsWrapper > main > div > div > div.articleOverview.test-articleOverview > div.articlePurchaseBox.css-gxzada > div.addToCartForm.css-1127cv2 > div.inputSelects.clearfix > div.test-sizeSelector.css-539bvd > ul", func(element *colly.HTMLElement) {
		sizes := make([]Size, 0)
		element.ForEach("button", func(i int, element *colly.HTMLElement) {
			sizes = append(sizes, Size{
				Name:         element.Text,
				Availability: !strings.Contains(element.Attr("class"), "disable"),
			})
		})
		data.Sizes = sizes
	})

	c.OnHTML("#__next > div > div.page.css-1umyepy > div.contentsWrapper > main > div > div > div.articleOverview.test-articleOverview > div.articleImageWrapper.clearfix.css-cdlca7 > div > div > div",
		func(element *colly.HTMLElement) {
			element.ForEach(".test-img", func(i int, element *colly.HTMLElement) {
				src := element.Attr("src")
				data.ImageUrls = append(data.ImageUrls, src)
			})
		})

	c.OnHTML(".sizeChartTable tbody", func(e *colly.HTMLElement) {
		sizeNames := make([]string, 0)
		sizeValues := make([][]string, 0)
		e.ForEach(".sizeChartTRow", func(i int, element *colly.HTMLElement) {
			sizeValues = append(sizeValues, make([]string, 0))
			element.ForEach(".sizeChartTCell", func(j int, element2 *colly.HTMLElement) {
				if i == 0 {
					sizeNames = append(sizeNames, element2.Text)
				} else {
					sizeValues[i-1] = append(sizeValues[i-1], element2.Text)
				}
			})
		})
		data.SizeTypeX = sizeNames
		data.SizeInfo = sizeValues
	})
	c.OnHTML(".sizeChart thead", func(element *colly.HTMLElement) {
		sizeTypeY := make([]string, 0)
		element.ForEach(".sizeChartTRow", func(i int, element2 *colly.HTMLElement) {
			sizeTypeY = append(sizeTypeY, element2.Text)
		})
		data.SizeTypeY = sizeTypeY
	})

	c.OnHTML("div.articleDisplay.test-articleDisplay", func(element *colly.HTMLElement) {
		data = Data{}
		element.ForEach("div > div > a", func(i int, element *colly.HTMLElement) {
			if i < 8 {
				return
			}
			if i > 10 {
				return
			}
			link := "https://shop.adidas.jp" + element.Attr("href")
			data.Link = link
			err := c.Visit(link)
			if err != nil {
				log.Fatal(err)
			}
			dataArr = append(dataArr, data)
		})
	})

	c.OnHTML("#BVRRRatingOverall_ > div.BVRRRatingNormalOutOf > span.BVRRNumber.BVRRRatingNumber > font > font", func(element *colly.HTMLElement) {
		overallRatingStr := element.Text
		overallRating, _ := strconv.ParseFloat(overallRatingStr, 0)
		data.OverallRating = overallRating
	})

	c.OnHTML("#BVRRDisplayContentID", func(element *colly.HTMLElement) {
		reviews := make([]Review, 0)
		reviewStar := make([]int, 0)
		reviewTitles := make([]string, 0)
		dates := make([]string, 0)
		reviewTexts := make([]string, 0)

		review := Review{}
		element.ForEach("#BVRRRatingOverall_Review_Display > div.BVRRRatingNormalImage > img", func(i int, element *colly.HTMLElement) {
			r, _ := strconv.Atoi(element.Attr("alt")[:1])
			reviewStar = append(reviewStar, r)
		})
		element.ForEach("#BVSubmissionPopupContainer > div.BVRRReviewDisplayStyle5Header > div.BVRRReviewTitleContainer", func(i int, element *colly.HTMLElement) {
			reviewTitles = append(reviewTitles, element.Text)
			review.Title = element.Text
		})
		element.ForEach(".BVRRReviewDate", func(i int, element2 *colly.HTMLElement) {
			dates = append(dates, element2.Text)
		})
		element.ForEach(".BVRRReviewText", func(i int, element2 *colly.HTMLElement) {
			reviewTexts = append(reviewTexts, element2.Text)
		})
		for i, _ := range reviewTitles {
			reviews = append(reviews, Review{
				Title: reviewTitles[i],
				Star:  reviewStar[i],
				Date:  dates[i],
				Text:  reviewTexts[i],
			})
		}
		data.Reviews = reviews
	})

	err := c.Visit("https://shop.adidas.jp/item/?gender=mens&category=wear&order=1&page=1")

	if err != nil {
		fmt.Println(err)
	}
	f, err := os.Create("adidas-data.csv")
	if err != nil {
		log.Fatal(err)
	}
	w := csv.NewWriter(f)
	rowsVals := make([][]string, 0)
	for _, d := range dataArr {
		rowVals := make([]string, 0)
		rowVals = append(rowVals, d.Title, d.Category, d.Link,
			strconv.FormatFloat(d.OverallRating, 'g', -1, 64))
		toStr := func(i interface{}) string {
			d, err := json.Marshal(i)
			if err != nil {
				log.Fatal(err)
			}
			if d == nil {
				return ""
			} else {
				return string(d)
			}
		}
		rowVals = append(rowVals, toStr(d.Reviews), toStr(d.Sizes), toStr(d.ImageUrls))
		rowsVals = append(rowsVals, rowVals)
	}
	err = w.WriteAll(rowsVals)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dataArr)
}
