package pokemon_store

import (
	"fmt"
	"log"

	"github.com/gocolly/colly"
)

type pokemon struct {
	url, image, name, price string
}

func Scrape() {
	var pokemonFound []pokemon
	var pagesToScrape Pages
	pagesVisited := Pages{"1"}
	iteration := 0
	c := colly.NewCollector()

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong: ", err)
	})

	c.OnHTML("a.page-numbers", func(e *colly.HTMLElement) {
		page := e.Text
		if pagesToScrape.contain(page) {
			return
		}
		if pagesVisited.contain(page) {
			return
		}
		pagesToScrape = append(pagesToScrape, page)
	})

	c.OnHTML("li.product", func(e *colly.HTMLElement) {
		p := pokemon{
			url:   e.ChildAttr("a", "href"),
			image: e.ChildAttr("img", "src"),
			name:  e.ChildText("h2"),
			price: e.ChildText("span.price"),
		}
		pokemonFound = append(pokemonFound, p)
	})

	c.OnScraped(func(r *colly.Response) {
		if len(pagesToScrape) == 0 {
			return
		}
		pageToScrape := pagesToScrape[iteration]
		pagesToScrape = pagesToScrape[1:]
		pagesVisited = append(pagesVisited, pageToScrape)
		iteration++
		url := getPageUrl(pageToScrape)
		c.Visit(url)
	})

	firstPage := getPageUrl("1")
	c.Visit(firstPage)
	fmt.Print(len(pokemonFound))
}

func getPageUrl(page string) string {
	return fmt.Sprintf("https://scrapeme.live/shop/page/%s", page)
}

type Pages []string

func (p Pages) contain(page string) bool {
	for _, iterationPage := range p {
		if iterationPage == page {
			return true
		}
	}
	return false
}
