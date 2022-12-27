package service

import (
	"github.com/alpaka-alpachino/job-engine/internal/models"
	"github.com/alpaka-alpachino/job-engine/internal/tests/data"
	"strings"
)

const (
	minVacanciesCount = 10
)

type Service struct {
	normalizer           *models.Normalizer
	professionsStatistic map[string][]models.ProfessionStatistic
	professionsWorkUA    []models.ProfessionWorkUA
	mapping              map[string][]string
}

func NewService(
	normalizer *models.Normalizer,
	professions map[string][]models.ProfessionStatistic,
	professionsWorkUA []models.ProfessionWorkUA,
	mapping map[string][]string,
) (*Service, error) {

	return &Service{
		normalizer:           normalizer,
		professionsStatistic: professions,
		professionsWorkUA:    professionsWorkUA,
		mapping:              mapping,
	}, nil
}

func (s *Service) normalizeProfileMatches(profileMatches map[string]int) {
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
}

func (s *Service) GetProfile(profileMatches map[string]int) (models.Profile, error) {
	var front, side, behind string

	s.normalizeProfileMatches(profileMatches)

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
		Front:             front,
		FrontDescription:  data.PsychoTypes[front],
		Side:              side,
		SideDescription:   data.PsychoTypes[side],
		Behind:            behind,
		BehindDescription: data.PsychoTypes[behind],
	}

	profile.ComplexType = front + side + behind

	return profile, nil
}

func (s *Service) SearchWorkUAProfessions(p models.Profile) ([]models.ProfessionWorkUA, error) {
	profs := s.searchWorkUAProfessionsByCodes(s.searchCodesByTypes([]string{p.ComplexType}, s.mapping))

	vacanciesCount := s.getVacanciesCount(profs)

	if vacanciesCount < minVacanciesCount {
		switch len(p.ComplexType) {
		case 1:
			profs = append(profs, s.searchWorkUAProfessionsByCodes(
				s.searchCodesByTypes([]string{p.Front}, s.mapping))...)
			if s.getVacanciesCount(profs) < minVacanciesCount {
				profs = s.extend(profs, p.Front)
			}
		case 2:
			profs = append(profs, s.searchWorkUAProfessionsByCodes(
				s.searchCodesByTypes([]string{p.Front + p.Side}, s.mapping))...)
			if s.getVacanciesCount(profs) < minVacanciesCount {
				profs = append(profs, s.searchWorkUAProfessionsByCodes(
					s.searchCodesByTypes([]string{
						p.Front,
						p.Side,
					}, s.mapping))...)
			}
		case 3:
			profs = append(profs, s.searchWorkUAProfessionsByCodes(
				s.searchCodesByTypes([]string{
					p.Front + p.Side,
					p.Front + p.Behind,
					p.Side + p.Behind,
				}, s.mapping))...)

			if s.getVacanciesCount(profs) < minVacanciesCount {
				profs = append(profs, s.searchWorkUAProfessionsByCodes(
					s.searchCodesByTypes([]string{
						p.Front,
						p.Side,
						p.Behind,
					}, s.mapping))...)
			}
		}

	}

	return profs, nil
}

func (s *Service) GetProfessionStatisticByWorkUAProfessions(profs []models.ProfessionWorkUA) ([]models.ProfessionStatistic, error) {
	uniqueCodes := make(map[string]struct{})
	for _, prof := range profs {
		uniqueCodes[prof.Code] = struct{}{}
	}

	ps := make([]models.ProfessionStatistic, 0)

	for code := range uniqueCodes {
		ps = append(ps, s.professionsStatistic[code]...)
	}

	return ps, nil
}

func (s *Service) searchWorkUAProfessionsByCodes(codes []string) []models.ProfessionWorkUA {
	m := make([]models.ProfessionWorkUA, 0)
	for _, code := range codes {
		for _, prof := range s.professionsWorkUA {
			if prof.Code == code {
				m = append(m, prof)
			}
		}
	}

	return m
}

func (s *Service) searchCodesByTypes(types []string, mapping map[string][]string) []string {
	uniqueCodes := make(map[string]struct{})

	for k, v := range mapping {
		for _, t := range types {
			if k == t {
				for _, v := range v {
					uniqueCodes[v] = struct{}{}
				}
			}
		}
	}

	codes := make([]string, 0, len(uniqueCodes))

	for code := range uniqueCodes {
		codes = append(codes, code)
	}

	return codes
}

func (s *Service) getVacanciesCount(p []models.ProfessionWorkUA) int {
	vacancies := 0
	for _, prof := range p {
		vacancies += prof.Vacancies
	}

	return vacancies
}

func (s *Service) MapProfessions(workUAProfessions []models.ProfessionWorkUA, professionStatistic []models.ProfessionStatistic) []models.Professions {
	var professions []models.Professions

	for _, vWorkUA := range workUAProfessions {
		for _, vStatistics := range professionStatistic {
			if vWorkUA.Code == vStatistics.Code {
				professions = append(professions, models.Professions{
					Name:            vStatistics.Name,
					Code:            vStatistics.Code,
					VacanciesWorkUA: vWorkUA.Vacancies,
					VUIndex:         vStatistics.VUIndex,
				})
			}
		}
	}

	return professions
}

func (s *Service) extend(profs []models.ProfessionWorkUA, front string) []models.ProfessionWorkUA {
	related := make([]string, 0)
	for t := range s.mapping {
		if strings.Contains(t, front) {
			related = append(related, t)
		}
	}

	profs = append(profs, s.searchWorkUAProfessionsByCodes(s.searchCodesByTypes(related, s.mapping))...)

	return profs
}
