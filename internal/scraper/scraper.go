package scraper

import (
	"fmt"
	"github.com/alpaka-alpachino/job-engine/internal/models"
	"github.com/gocolly/colly"
	"strconv"
	"strings"
)

const workUaBaseURL = "https://www.work.ua/jobs/by-titles/?page="

func ScrapeWorkUA() ([]models.ProfessionWorkUA, error) {
	var professions []models.ProfessionWorkUA
	var pages int
	var err error

	collectPagesQuantity := colly.NewCollector()
	collectProfessions := colly.NewCollector()

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

		professions = append(professions, models.ProfessionWorkUA{
			Name:      strings.TrimSpace(strings.ToLower(profession)),
			Vacancies: quantity,
		})
	})

	for i := 1; i <= pages; i++ {
		collectProfessions.Visit(fmt.Sprintf("%s%d", workUaBaseURL, i))
	}

	return professions, nil
}

func SetCodes(coefficient float64, professionsStatistic map[string][]models.ProfessionStatistic, professionsWorkUA []models.ProfessionWorkUA) []models.ProfessionWorkUA {
	var undefined, undefinedProcessed, result []models.ProfessionWorkUA

	for _, professionWorkUA := range professionsWorkUA {
		var defined bool
		for _, v := range professionsStatistic {
			for _, professionStattistic := range v {
				if professionWorkUA.Name == professionStattistic.Name {
					result = append(result, models.ProfessionWorkUA{
						Name:      professionWorkUA.Name,
						Code:      professionStattistic.Code,
						Vacancies: professionWorkUA.Vacancies,
					})
					defined = true
				}
			}
		}
		if !defined {
			undefined = append(undefined, professionWorkUA)
		}

	}

	for _, vUndefined := range undefined {
		var defined bool
		for _, v := range professionsStatistic {
			for _, professionStattistic := range v {
				b := []byte(vUndefined.Name)
				a := []byte(professionStattistic.Name)
				sim := JaroWrinkler(b, a)
				if sim > coefficient {
					result = append(result, models.ProfessionWorkUA{
						Name:      vUndefined.Name,
						Code:      professionStattistic.Code,
						Vacancies: vUndefined.Vacancies,
					})
					defined = true
					continue
				}
			}
		}
		if !defined {
			undefinedProcessed = append(undefinedProcessed, vUndefined)
		}

	}

	fmt.Println("UNDEFINED", undefined)

	return result
}

func JaroWrinkler(a, b []byte) float64 {

	lenPrefix := len(commonPrefix(a, b))
	if lenPrefix > 4 {
		lenPrefix = 4
	}
	// Return similarity.
	similarity := JaroSim(a, b)
	return similarity + (0.1 * float64(lenPrefix) * (1.0 - similarity))

}

func JaroSim(str1, str2 []byte) float64 {
	if len(str1) == 0 && len(str2) == 0 {
		return 1
	}
	if len(str1) == 0 || len(str2) == 0 {
		return 0
	}
	match_distance := len(str1)
	if len(str2) > match_distance {
		match_distance = len(str2)
	}
	match_distance = match_distance/2 - 1
	str1_matches := make([]bool, len(str1))
	str2_matches := make([]bool, len(str2))
	matches := 0.
	transpositions := 0.
	for i := range str1 {
		start := i - match_distance
		if start < 0 {
			start = 0
		}
		end := i + match_distance + 1
		if end > len(str2) {
			end = len(str2)
		}
		for k := start; k < end; k++ {
			if str2_matches[k] {
				continue
			}
			if str1[i] != str2[k] {
				continue
			}
			str1_matches[i] = true
			str2_matches[k] = true
			matches++
			break
		}
	}
	if matches == 0 {
		return 0
	}
	k := 0
	for i := range str1 {
		if !str1_matches[i] {
			continue
		}
		for !str2_matches[k] {
			k++
		}
		if str1[i] != str2[k] {
			transpositions++
		}
		k++
	}
	transpositions /= 2
	return (matches/float64(len(str1)) +
		matches/float64(len(str2)) +
		(matches-transpositions)/matches) / 3
}

func commonPrefix(first, second []byte) []byte {
	if len(first) > len(second) {
		first, second = second, first
	}

	var commonLen int
	sRunes := []byte(second)
	for i, r := range first {
		if r != sRunes[i] {
			break
		}
		commonLen++
	}
	return sRunes[0:commonLen]
}
