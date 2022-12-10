package server

import (
	"fmt"
	"github.com/alpaka-alpachino/job-engine/internal/models"
	"github.com/alpaka-alpachino/job-engine/internal/service"
	"github.com/alpaka-alpachino/job-engine/internal/tests"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

func newRouter(s *service.Service, t *template.Template) (*mux.Router, error) {
	router := mux.NewRouter()
	test, err := tests.GetTest()
	if err != nil {
		return nil, err
	}

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")
		if name != "" {
			test.Name = name
		}

		err = t.ExecuteTemplate(w, "test.html", test)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}).Methods(http.MethodGet)

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err = r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var result []string
		for i := 0; i < 42; i++ {
			v := r.FormValue(fmt.Sprintf("%d", i))
			if v != "" {
				result = append(result, v)
				continue
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

		profile, err := s.GetProfile(profileMatches)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		workUAProfessions, err := s.SearchWorkUAProfessions(profile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		professionStatistic, err := s.GetProfessionStatisticByWorkUAProfessions(workUAProfessions)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		professions := s.MapProfessions(workUAProfessions, professionStatistic)

		summary := models.Summary{
			Profile:     profile,
			Professions: professions,
		}

		if err := t.ExecuteTemplate(w, "result.html", summary); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}).Methods(http.MethodPost)

	return router, nil
}
