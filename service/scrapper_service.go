package service

import (
	"buzzwordBot/models"
	"fmt"

	"github.com/gocolly/colly"
)

func ScrapeHackerNewsLinks() []models.Item {
	var itemList []models.Item
	c := colly.NewCollector()

	c.OnHTML("tr.athing", func(e *colly.HTMLElement) {
		e.ForEach("td.title", func(_ int, el *colly.HTMLElement) {
			if el.ChildAttr("a", "href") == "" || el.ChildAttr("a.titlelink", "href")[0:4] == "item" {
				return
			}
			item := models.Item{
				Id:    e.Attr("id"),
				Title: el.ChildText("a.titlelink"),
				Link:  el.ChildAttr("a.titlelink", "href"),
			}
			itemList = append(itemList, item)
		})
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.Visit("https://news.ycombinator.com/front")

	return itemList
}
