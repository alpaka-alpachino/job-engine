package models

type Normalizer struct {
	Normalizer []NormalizeByTypes `json:"types"`
}

type NormalizeByTypes struct {
	Name   string  `json:"name"`
	Scores []Score `json:"scores"`
}

type Score struct {
	Raw    []int `json:"raw"`
	Normal int   `json:"normal"`
}
