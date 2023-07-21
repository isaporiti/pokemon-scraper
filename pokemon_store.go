package pokemon_store

import (
	"fmt"
	"log"

	"github.com/gocolly/colly"
)

type pokemon struct {
	url, image, name, price string
}

type pokemonScraper struct {
	pokemonFound  []pokemon
	pagesToScrape Pages
	pagesVisited  Pages
	iteration     int
}

func NewScraper() *pokemonScraper {
	s := pokemonScraper{}
	s.pagesVisited = append(s.pagesVisited, "1")
	return &s
}

func (s *pokemonScraper) registerPage(page string) {
	if s.pagesToScrape.contain(page) {
		return
	}
	if s.pagesVisited.contain(page) {
		return
	}
	s.pagesToScrape = append(s.pagesToScrape, page)
}

func (s *pokemonScraper) savePokemon(p pokemon) {
	s.pokemonFound = append(s.pokemonFound, p)
}

func (s *pokemonScraper) nextPage() (url string) {
	if len(s.pagesToScrape) == 0 {
		return ""
	}
	pageToScrape := s.pagesToScrape[s.iteration]
	url = getPageUrl(pageToScrape)
	s.pagesToScrape = s.pagesToScrape[1:]
	s.pagesVisited = append(s.pagesVisited, pageToScrape)
	s.iteration++
	return
}

func (s *pokemonScraper) pokemon() []pokemon {
	return s.pokemonFound
}

func Scrape() {
	scraper := NewScraper()
	c := colly.NewCollector()

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong: ", err)
	})

	c.OnHTML("a.page-numbers", func(e *colly.HTMLElement) {
		page := e.Text
		scraper.registerPage(page)
	})

	c.OnHTML("li.product", func(e *colly.HTMLElement) {
		p := pokemon{
			url:   e.ChildAttr("a", "href"),
			image: e.ChildAttr("img", "src"),
			name:  e.ChildText("h2"),
			price: e.ChildText("span.price"),
		}
		scraper.savePokemon(p)
	})

	c.OnScraped(func(r *colly.Response) {
		url := scraper.nextPage()
		c.Visit(url)
	})

	firstPage := getPageUrl("1")
	c.Visit(firstPage)
	fmt.Print(len(scraper.pokemon()))
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
