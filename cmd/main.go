package main

import (
	"github.com/alpaka-alpachino/job-engine/config"
	"github.com/alpaka-alpachino/job-engine/internal/data"
	"github.com/alpaka-alpachino/job-engine/internal/server"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
	"html/template"
	"log"
)

func main() {
	// Setup logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Println("Can't build logger")
	}

	defer func(logger *zap.Logger) {
		err = logger.Sync()
		if err != nil {
			log.Println("Can't flush logging buffer")
		}
	}(logger)

	l := logger.Sugar()

	// Construct project config
	c, err := config.NewEngineConfig()
	if err != nil {
		l.With(err).Fatal("Can't read envs")
	}

	// Initialise frontend templates
	t := template.Must(template.ParseFiles("template/test.html", "template/result.html"))

	f, err := excelize.OpenFile("internal/data/prof.xlsx")
	if err != nil {
		l.With(err).Fatal("Can't open xlsx professions file")
	}

	categories, err := data.GetProfessionsMap(f,"Дані", 5)
	if err != nil {
		l.With(err).Fatal("Can't get professions' categories")
	}
	// Initialize server
	s, err := server.NewServer(c, t, categories)
	if err != nil {
		l.With(err).Fatal("Can't setup server")
	}

	l.Fatal(s.RunServer())
}
