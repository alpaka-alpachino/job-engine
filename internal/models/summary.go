package models

type Summary struct {
	Profile     Profile
	Professions []Professions
}

type Professions struct {
	Name            string
	Code            string
	VacanciesWorkUA int
	VUIndex         float64
}
