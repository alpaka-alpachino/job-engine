package tests

import (
	"bufio"
	"encoding/json"
	"github.com/alpaka-alpachino/job-engine/internal/models"
	"os"
)

const (
	psychoTestPath = "internal/tests/data/psycho-test.json"
	normTablePath  = "internal/tests/data/norm-table.json"
)

func GetTest() (*models.Test, error) {
	test := models.Test{}
	file, err := os.Open(psychoTestPath)
	if err != nil {
		return nil, err
	}
	r := bufio.NewReader(file)
	if err := json.NewDecoder(r).Decode(&test); err != nil {
		return nil, err
	}
	file.Close()

	return &test, err
}

func GetNormalizer() (*models.Normalizer, error) {
	normalizer := models.Normalizer{}

	file, _ := os.Open(normTablePath)
	r := bufio.NewReader(file)
	err := json.NewDecoder(r).Decode(&normalizer)
	if err != nil {
		return nil, err
	}
	file.Close()

	return &normalizer, err
}
