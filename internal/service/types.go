package service

import (
	"github.com/xuri/excelize/v2"
)

func GetTypeToCodesMapping() (map[string][]string, error) {
	f, err := excelize.OpenFile(mappingPath)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	rows, err := f.GetRows("Тест-Проф'ія-Сп'ність")
	if err != nil {
		return nil, err
	}

	mapping := make(map[string][]string)

	for _, row := range rows {
		t := row[0]
		c := row[1]

		codes, ok := mapping[t]
		switch ok {
		case true:
			codes = append(codes, c)
		default:
			codes = []string{c}
		}

		mapping[t] = codes
	}

	for t, codes := range mapping {
		m := make(map[string]struct{})

		for _, code := range codes {
			m[code] = struct{}{}
		}

		codes := make([]string, 0, len(m))
		for code := range m {
			codes = append(codes, code)
		}

		mapping[t] = codes
	}

	return mapping, nil
}
