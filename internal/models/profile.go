package models

type Profile struct {
	Front       string
	Side        string
	Behind      string
	ComplexType string

	FrontScore        int
	FrontDescription  string
	SideScore         int
	SideDescription   string
	BehindScore       int
	BehindDescription string

	Professions map[string]Category
}
