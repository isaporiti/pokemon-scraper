package pokemon_store

import (
	"log"

	"github.com/gocolly/colly"
)

type pokemon struct {
	url, image, name, price string
}

func Scrape() {
	var pokemonFound []pokemon
	c := colly.NewCollector()

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong: ", err)
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

	c.Visit("https://scrapeme.live/shop/")
}
