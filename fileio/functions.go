package fileio

import "regexp"

func FileLines(filepath string) ([]string, error) {
	storage := CurrStorage()

	content, err := storage.Load(filepath)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`[\r\n]+`)
	return re.Split(content, -1), nil
}
