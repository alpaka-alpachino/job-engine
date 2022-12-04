package models

type Test struct {
	Name      string     `json:"name"`
	Questions []Question `json:"questions"`
}

type Question struct {
	Variants []Variant `json:"variants"`
}

type Variant struct {
	Variant string `json:"variant"`
	Type    string `json:"type"`
}
