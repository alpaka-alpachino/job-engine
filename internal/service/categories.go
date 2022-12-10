package service

import (
	"github.com/alpaka-alpachino/job-engine/internal/models"
	"github.com/xuri/excelize/v2"
	"strconv"
	"strings"
)

const (
	tablePath   = "internal/service/data/prof.xlsx"
	mappingPath = "internal/service/data/Holland_simp.xlsx"
	sheetName   = "Дані"

	jobNameRowIndex    int = 1
	codeRowIndex       int = 2
	vacancyRowIndex    int = 3
	unemployedRowIndex int = 4

	firstJob int = 19
)

func GetProfessionsMap() (map[string][]models.ProfessionStatistic, error) {
	f, err := excelize.OpenFile(tablePath)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	m := make(map[string][]models.ProfessionStatistic)

	for i, row := range rows {
		if i < firstJob {
			continue
		}

		code := row[codeRowIndex]

		vacancies, err := strconv.Atoi(row[vacancyRowIndex])
		if err != nil {
			return nil, err
		}

		if vacancies == 0 {
			continue
		}

		name := strings.TrimSpace(strings.ToLower(row[jobNameRowIndex]))

		unemployed, err := strconv.Atoi(row[unemployedRowIndex])
		if err != nil {
			return nil, err
		}

		var vuIndex float64
		if unemployed > 0 {
			vuIndex = float64(vacancies) / float64(unemployed)
		}

		if professions, ok := m[code]; ok {
			professions = append(professions, models.ProfessionStatistic{
				Name:       name,
				Code:       code,
				Vacancies:  vacancies,
				Unemployed: unemployed,
				VUIndex:    vuIndex,
			})

			m[code] = professions
		} else {
			m[code] = []models.ProfessionStatistic{{
				Name:       name,
				Code:       code,
				Vacancies:  vacancies,
				Unemployed: unemployed,
				VUIndex:    vuIndex,
			}}
		}
	}

	return m, nil
}
