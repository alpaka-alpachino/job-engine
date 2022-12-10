package main

import (
	"github.com/alpaka-alpachino/job-engine/config"
	"github.com/alpaka-alpachino/job-engine/internal/scraper"
	"github.com/alpaka-alpachino/job-engine/internal/server"
	"github.com/alpaka-alpachino/job-engine/internal/service"
	"github.com/alpaka-alpachino/job-engine/internal/tests"
	"go.uber.org/zap"
	"html/template"
	"log"
)

func main() {
	// Setup logger
	logger, err := zap.NewDevelopment()
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

	professionsStatistic, err := service.GetProfessionsMap()
	if err != nil {
		l.With(err).Fatal("Can't get professions")
	}

	mapping, err := service.GetTypeToCodesMapping()
	if err != nil {
		l.With(err).Fatal("Can't get mapping")
	}

	professionsWorkUA, err := scraper.ScrapeWorkUA()
	if err != nil {
		log.Fatal(err)
	}
	//2669
	professionsWorkUA = scraper.SetCodes(c.SimilarityCoefficient, professionsStatistic, professionsWorkUA)

	normalizer, err := tests.GetNormalizer()
	if err != nil {
		l.With(err).Fatal("Can't get normalizer")
	}

	engineService, err := service.NewService(normalizer, professionsStatistic, professionsWorkUA, mapping)
	if err != nil {
		l.With(err).Fatal("Can't create service instance")
	}

	// Initialize server
	s, err := server.NewServer(c, engineService, t)
	if err != nil {
		l.With(err).Fatal("Can't setup server")
	}

	l.Fatal(s.RunServer())
}
