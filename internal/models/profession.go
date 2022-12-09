package models

type ByTypes struct {
	ByTypes []ProfessionByType `json:"by_types"`
}

// ProfessionByType test types and appropriate professions categories
type ProfessionByType struct {
	ProfessionType string   `json:"type"`
	Description    string   `json:"description"`
	Professions    []string `json:"professions"`
}
