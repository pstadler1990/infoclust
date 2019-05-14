package translate

import (
	"bufio"
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

func Translate(keywordsPath, translatedPath string) {
	translatedSliceRaw, errTranslated := readJsonRaw(translatedPath)
	//keywordsSlice, errKeywords := readJsonRaw(keywordsPath)

	if errTranslated != nil /*|| errKeywords != nil*/ {
		//...
	}

	/* Selector for extracted_keywords file */
	//reKeywords := regexp.MustCompile(`"(.+)": \d+`)

	/* Selector for bow_translated file */
	reBow := regexp.MustCompile(`\s*"(.+)",?\s*`)

	var translatedSlice []string

	for _, s := range translatedSliceRaw {
		if reBow.MatchString(s) {
			translatedSlice = append(translatedSlice, reBow.ReplaceAllString(s, "$1"))
		}
	}
}


