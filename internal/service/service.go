package service

import (
	"bufio"
	"encoding/json"
	"github.com/alpaka-alpachino/job-engine/internal/models"
	"math"
	"os"
)

const professionsByTypesPath = "internal/service/data/professions-by-types.json"

type Service struct {
	normalizer        *models.Normalizer
	professions       *models.ByTypes
	categories        map[string]models.Category
	professionsWorkUA map[string]int
}

func NewService(normalizer *models.Normalizer, categories map[string]models.Category, professionsWorkUA map[string]int) (*Service, error) {
	professions, err := getProfessionsByType()
	if err != nil {
		return nil, err
	}

	return &Service{
		normalizer:        normalizer,
		professions:       professions,
		categories:        categories,
		professionsWorkUA: professionsWorkUA,
	}, nil
}

func (s *Service) GetProfile(profileMatches map[string]int) (models.Profile, error) {
	for k, v := range profileMatches {
		for _, vN := range s.normalizer.Normalizer {
			if k == vN.Name {
				for _, vS := range vN.Scores {
					for _, vR := range vS.Raw {
						if v == vR {
							profileMatches[k] = vS.Normal
						}
					}
				}
			}
		}
	}

	var front, side, behind string

	for k, v := range profileMatches {
		if k != front {
			if v > profileMatches[front] {
				front = k
			} else if k != side {
				if v > profileMatches[side] {
					side = k
				} else if k != behind {
					if v > profileMatches[behind] {
						behind = k
					}
				}
			}
		}
	}

	profile := models.Profile{
		Front:  front,
		Side:   side,
		Behind: behind,
	}

	profile.ComplexType = front + side + behind

	for _, v := range s.professions.ByTypes {
		if v.ProfessionType == profile.Front {
			profile.FrontScore = profileMatches[front]
			profile.FrontDescription = v.Description
			profile.Professions = getCategoriesByNames(v.Professions, s.categories)
		} else if v.ProfessionType == profile.Side {
			profile.SideScore = profileMatches[side]
			profile.SideDescription = v.Description
		} else if v.ProfessionType == profile.Behind {
			profile.BehindScore = profileMatches[behind]
			profile.BehindDescription = v.Description
		}
	}

	return profile, nil
}

func getProfessionsByType() (*models.ByTypes, error) {
	professions := models.ByTypes{}
	file, _ := os.Open(professionsByTypesPath)
	r := bufio.NewReader(file)
	err := json.NewDecoder(r).Decode(&professions)
	if err != nil {
		return nil, err
	}
	file.Close()

	return &professions, nil
}

func getCategoriesByNames(nameList []string, data map[string]models.Category) map[string]models.Category {
	m := make(map[string]models.Category)

	// TODO include workUA results
	for _, name := range nameList {
		if category, ok := data[name]; ok {
			category.VUIndex = math.Round(category.VUIndex*100) / 100
			m[name] = category
		}
	}

	return m
}
