package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/alpaka-alpachino/job-engine/internal/data"
	"github.com/alpaka-alpachino/job-engine/internal/models"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"os"
)

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

	Professions map[string]data.Category
}

func newRouter(t *template.Template, categories map[string]data.Category) (*mux.Router, error) {
	router := mux.NewRouter()

	// Read all necessary files into proper models
	test := models.Test{}
	file, err := os.Open("internal/data/psycho-test.json")
	if err != nil {
		log.Println("H")
	}
	r := bufio.NewReader(file)
	if err := json.NewDecoder(r).Decode(&test); err != nil {
		return nil, err
	}
	file.Close()

	normalizer := models.Normalizer{}
	file, _ = os.Open("internal/data/norm_table.json")
	r = bufio.NewReader(file)
	if err := json.NewDecoder(r).Decode(&normalizer); err != nil {
		return nil, err
	}
	file.Close()

	professions := models.ByTypes{}
	file, _ = os.Open("internal/data/professions.json")
	r = bufio.NewReader(file)
	if err := json.NewDecoder(r).Decode(&professions); err != nil {
		return nil, err
	}
	file.Close()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if name := r.FormValue("name"); name != "" {
			test.Name = name
		}

		if err := t.ExecuteTemplate(w, "test.html", test); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}).Methods(http.MethodGet)

	//////////////////////////////////////////////////////////////////////////////////////////
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			return
		}

		var result []string
		var v string
		for i := 0; i < 42; i++ {
			v = r.FormValue(fmt.Sprintf("%d", i))
			if v != "" {
				result = append(result, v)
			}
		}

		matches := make(map[string]string)
		for _, v := range test.Questions {
			for _, vv := range v.Variants {
				matches[vv.Variant] = vv.Type
			}
		}

		profileMatches := make(map[string]int)

		for _, v := range result {
			for ii, vv := range matches {
				if v == ii {
					profileMatches[vv] = profileMatches[vv] + 1
				}
			}
		}

		for k, v := range profileMatches {
			for _, vN := range normalizer.Normalizer {
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

		profile := Profile{
			Front:  front,
			Side:   side,
			Behind: behind,
		}

		profile.ComplexType = front + side + behind

		for _, v := range professions.ByTypes {
			if v.ProfessionType == profile.Front {
				profile.FrontScore = profileMatches[front]
				profile.FrontDescription = v.Description
				profile.Professions = data.GetCategoriesByNames(v.Professions, categories)
				//profile.Professions = strings.Join(v.Professions, ", ")
			} else if v.ProfessionType == profile.Side {
				profile.SideScore = profileMatches[side]
				profile.SideDescription = v.Description
			} else if v.ProfessionType == profile.Behind {
				profile.BehindScore = profileMatches[behind]
				profile.BehindDescription = v.Description
			}
		}

		if err := t.ExecuteTemplate(w, "result.html", profile); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}).Methods(http.MethodPost)

	return router, nil
}
