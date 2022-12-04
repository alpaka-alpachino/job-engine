package models

type ByTypes struct {
	ByTypes []ProfessionByType `json:"by_types"`
}

type ProfessionByType struct {
	ProfessionType string   `json:"type"`
	Description    string   `json:"description"`
	Professions    []string `json:"professions"`
}
