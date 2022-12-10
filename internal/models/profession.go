package models

// ProfessionStatistic is a model to enterprise profession from prof.xls table
type ProfessionStatistic struct {
	Name       string
	Code       string
	Vacancies  int
	Unemployed int
	// vacancies quantity to unemployed (if more then vacancies quantity more)
	VUIndex float64
}

// ProfessionWorkUA is a model to enterprise profession from workUA table
type ProfessionWorkUA struct {
	Name      string
	Code      string
	Vacancies int
}
