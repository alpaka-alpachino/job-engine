package tools

import (
	"fmt"
	"github.com/gocolly/colly"
	"strconv"
	"strings"
)

const workUaBaseURL = "https://www.work.ua/jobs/by-titles/?page="

func ScrapeWorkUA() (map[string]int, error) {
	collectPagesQuantity := colly.NewCollector()
	collectProfessions := colly.NewCollector()

	professions := make(map[string]int)
	var pages int
	var err error
	collectPagesQuantity.OnHTML("#center > div > div.card > nav > ul.pagination.hidden-xs > li:nth-child(6) > a", func(e *colly.HTMLElement) {
		pages, err = strconv.Atoi(e.Text)
		if err != nil {
			return
		}
	})

	err = collectPagesQuantity.Visit(fmt.Sprintf("%s1", workUaBaseURL))
	if err != nil {
		return nil, err
	}

	collectProfessions.OnHTML("#center > div > div.card > div.row > div > ul > li", func(e *colly.HTMLElement) {
		var profession string
		var quantity int

		e.ForEach("a", func(_ int, internal *colly.HTMLElement) {
			data := internal.Attr("title")
			if data != "" {
				profession = data
			}
		})

		e.ForEach("span", func(_ int, internal *colly.HTMLElement) {
			data := strings.TrimSpace(internal.Text)
			if data != "" {
				quantity, err = strconv.Atoi(data)
				if err != nil {
					return
				}
			}
		})

		professions[profession] = quantity
	})

	for i := 1; i <= pages; i++ {
		collectProfessions.Visit(fmt.Sprintf("%s%d", workUaBaseURL, i))
	}

	return professions, nil
}
