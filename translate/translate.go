package translate

import (
	"bufio"
	"errors"
	"os"
	"regexp"
)


func readJsonRaw(path string) ([]string, error) {
	lines := make([]string, 1)

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

func writeTranslatedSliceToFile(slice []string, path string) error {
	file, err := os.Create(path)

	if err != nil {
		return err
	}

	defer file.Close()

	for _, s := range slice {
		_, err := file.WriteString(s)
		if err != nil {
			return err
		}
	}

	return nil
}

func Translate(keywordsPath, translatedPath string) error {
	translatedSliceRaw, errTranslated := readJsonRaw(translatedPath)
	keywordsSliceRaw, errKeywords := readJsonRaw(keywordsPath)

	if errTranslated != nil {
		return errTranslated
	}

	if errKeywords != nil {
		return errKeywords
	}

	/* Selector for extracted_keywords file */
	reKeywords := regexp.MustCompile(`"(.+)": \d+`)
	reKeywordsValue := regexp.MustCompile(`"(.+)"`)

	/* Selector for bow_translated file */
	reBow := regexp.MustCompile(`\s*"(.+)",?\s*`)

	var translatedSlice []string

	for _, s := range translatedSliceRaw {
		if reBow.MatchString(s) {
			translatedSlice = append(translatedSlice, reBow.ReplaceAllString(s, "\"$1\""))
		}
	}

	if translatedSlice == nil {
		return errors.New("Translated slice is nil")
	}

	var keywordsSlice []string

	counter := 0
	for _, s := range keywordsSliceRaw {
		if reKeywords.MatchString(s) {
			/* Overwrite file string */
			s = reKeywordsValue.ReplaceAllString(s, translatedSlice[counter])
			counter += 1
		}
		/* Write raw string (both modified and original) */
		keywordsSlice = append(keywordsSlice, s)
	}

	return writeTranslatedSliceToFile(keywordsSlice, "out.json")
}