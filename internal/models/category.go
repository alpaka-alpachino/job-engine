package models

// Category is a model with keywords to find professions
type Category struct {
	Name            string // first name-pattern
	Vacancies       int
	Unemployed      int
	UnemployedMen   int
	UnemployedWomen int
	VUIndex         float64 // vacancies quantity to unemployed (if more then vacancies quantity more)
	Profs           []Prof
}

// Prof is the model used to enterprise results for user
type Prof struct {
	Name            string
	Vacancies       int
	Unemployed      int
	UnemployedMen   int
	UnemployedWomen int
	VUIndex         float64
}
