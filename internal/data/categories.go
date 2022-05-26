package data

import (
	"math"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

const (
	jobNameRowIndex         int = 1
	vacancyRowIndex         int = 3
	unemployedRowIndex      int = 4
	unemployedWomenRowIndex int = 5

	firstJob int = 19
)

type Category struct { // то по каким ключевым словам мы ищем в таблице
	Name            string // перше слово-паттерн
	Vacancies       int
	Unemployed      int
	UnemployedMen   int
	UnemployedWomen int
	VUIndex         float64 // вакансии поделить на безработных если больше то вакансий больше
	Profs           []Prof
}

type Prof struct { // то что будем показывать юзеру
	Name            string
	Vacancies       int
	Unemployed      int
	UnemployedMen   int
	UnemployedWomen int
	VUIndex         float64
}

func GetProfessionsMap(f *excelize.File, sheetName string, minCount int) (map[string]Category, error) {
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil,err
	}

	m := make(map[string]Category)

	for i, row := range rows {
		if i < firstJob {
			continue
		}

		vacancies, err := strconv.Atoi(row[vacancyRowIndex])
		if err != nil {
			return nil,err
		}

		if vacancies == 0 {
			continue
		}

		name := row[jobNameRowIndex]
		categoryName := strings.TrimSpace(strings.ToLower(row[jobNameRowIndex]))
		for _, symbol := range name {
			switch {
			case symbol == ' ':
				categoryName = strings.Split(categoryName, " ")[0]
				break
			case symbol == '-':
				categoryName = strings.Split(categoryName, "-")[0]
				break
			}
		}

		unemployed, err := strconv.Atoi(row[unemployedRowIndex])
		if err != nil {
			return nil,err
		}

		unemployedWomen, err := strconv.Atoi(row[unemployedWomenRowIndex])
		if err != nil {
			return nil,err
		}

		unemployedMen := unemployed - unemployedWomen

		var vuIndex float64
		if unemployed > 0 {
			vuIndex = float64(vacancies) / float64(unemployed)
		}

		if category, ok := m[categoryName]; ok {
			c := Category{
				Name:            category.Name,
				Vacancies:       category.Vacancies + vacancies,
				Unemployed:      category.Unemployed + unemployed,
				UnemployedMen:   category.UnemployedMen + unemployedMen,
				UnemployedWomen: category.UnemployedWomen + unemployedWomen,
				Profs: append(category.Profs, Prof{
					Name:            name,
					Vacancies:       vacancies,
					Unemployed:      unemployed,
					UnemployedMen:   unemployedMen,
					UnemployedWomen: unemployedWomen,
					VUIndex:         vuIndex,
				}),
			}

			m[categoryName] = c
		} else {
			category := Category{
				Name:            categoryName,
				Vacancies:       vacancies,
				Unemployed:      unemployed,
				UnemployedMen:   unemployedMen,
				UnemployedWomen: unemployedWomen,
				Profs: []Prof{{
					Name:            name,
					Vacancies:       vacancies,
					Unemployed:      unemployed,
					UnemployedMen:   unemployedMen,
					UnemployedWomen: unemployedWomen,
					VUIndex:         vuIndex,
				}},
			}

			m[categoryName] = category
		}
	}

	for name, category := range m {
		category.VUIndex = float64(category.Vacancies) / float64(category.Unemployed)
		m[name] = category
	}

	FilterByVacancyCount(m, 5)

	return m, nil
}

func FilterByVacancyCount(data map[string]Category, min int) {
	for key, value := range data {
		if len(value.Profs) < min {
			delete(data, key)
		}
	}
}

func GetCategoriesByNames(nameList []string, data map[string]Category) map[string]Category {
	m := make(map[string]Category)

	for _, name := range nameList {
		if category, ok  := data[name]; ok {
			category.VUIndex=math.Round(category.VUIndex*100)/100
			m[name] = category
		}
	}

	return m
}