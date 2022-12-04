package service

import (
	"github.com/alpaka-alpachino/job-engine/internal/models"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

const (
	tablePath = "internal/service/data/prof.xlsx"
	sheetName = "Дані"

	jobNameRowIndex         int = 1
	vacancyRowIndex         int = 3
	unemployedRowIndex      int = 4
	unemployedWomenRowIndex int = 5

	firstJob int = 19
)

func GetProfessionsMap() (map[string]models.Category, error) {
	f, err := excelize.OpenFile(tablePath)
	if err != nil {
		return nil, err
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	m := make(map[string]models.Category)

	for i, row := range rows {
		if i < firstJob {
			continue
		}

		vacancies, err := strconv.Atoi(row[vacancyRowIndex])
		if err != nil {
			return nil, err
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
			return nil, err
		}

		unemployedWomen, err := strconv.Atoi(row[unemployedWomenRowIndex])
		if err != nil {
			return nil, err
		}

		unemployedMen := unemployed - unemployedWomen

		var vuIndex float64
		if unemployed > 0 {
			vuIndex = float64(vacancies) / float64(unemployed)
		}

		if category, ok := m[categoryName]; ok {
			c := models.Category{
				Name:            category.Name,
				Vacancies:       category.Vacancies + vacancies,
				Unemployed:      category.Unemployed + unemployed,
				UnemployedMen:   category.UnemployedMen + unemployedMen,
				UnemployedWomen: category.UnemployedWomen + unemployedWomen,
				Profs: append(category.Profs, models.Prof{
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
			category := models.Category{
				Name:            categoryName,
				Vacancies:       vacancies,
				Unemployed:      unemployed,
				UnemployedMen:   unemployedMen,
				UnemployedWomen: unemployedWomen,
				Profs: []models.Prof{{
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

	filterByVacancyCount(m, 5)

	return m, nil
}

func filterByVacancyCount(data map[string]models.Category, min int) {
	for key, value := range data {
		if len(value.Profs) < min {
			delete(data, key)
		}
	}
}
