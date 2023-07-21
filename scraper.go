package pokemon_store

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gocolly/colly"
)

type pokemon struct {
	url, image, name, price string
}

type pokemonScraper struct {
	*colly.Collector
	pokemonFound  []pokemon
	pagesToScrape Pages
	pagesVisited  Pages
	iteration     int
}

func NewScraper() *pokemonScraper {
	s := pokemonScraper{}
	s.Collector = colly.NewCollector()
	s.pagesVisited = append(s.pagesVisited, "1")
	return &s
}

func (s *pokemonScraper) registerPage(page string) {
	if _, err := strconv.Atoi(page); err != nil {
		return
	}
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
	if s.iteration >= len(s.pagesToScrape) {
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

func (s *pokemonScraper) scrape() {
	s.OnError(func(_ *colly.Response, err error) {
		log.Fatalln("Something went wrong: ", err)
	})

	s.OnHTML("a.page-numbers", func(e *colly.HTMLElement) {
		page := e.Text
		s.registerPage(page)
	})

	s.OnHTML("li.product", func(e *colly.HTMLElement) {
		p := pokemon{
			url:   e.ChildAttr("a", "href"),
			image: e.ChildAttr("img", "src"),
			name:  e.ChildText("h2"),
			price: e.ChildText("span.price"),
		}
		s.savePokemon(p)
	})

	s.OnScraped(func(r *colly.Response) {
		url := s.nextPage()
		if url == "" {
			return
		}
		s.Visit(url)
	})

	firstPage := getPageUrl("1")
	log.Println("Started scraping")
	s.Visit(firstPage)
	log.Println("Finished scraping")
}

func Scrape() error {
	scraper := NewScraper()
	scraper.scrape()
	err := saveToCsv(scraper)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func saveToCsv(scraper *pokemonScraper) error {
	log.Println("Started writing to CSV")
	file, err := os.Create("pokemon.csv")
	if err != nil {
		return fmt.Errorf("could not create file: %s", err)
	}
	defer file.Close()
	csvWriter := csv.NewWriter(file)
	defer csvWriter.Flush()

	headers := []string{
		"url",
		"image",
		"name",
		"price",
	}
	csvWriter.Write(headers)
	for _, pokemonToWrite := range scraper.pokemon() {
		row := []string{
			pokemonToWrite.url,
			pokemonToWrite.image,
			pokemonToWrite.name,
			pokemonToWrite.price,
		}
		csvWriter.Write(row)
	}
	log.Println("Finished writing to CSV")
	return nil
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
