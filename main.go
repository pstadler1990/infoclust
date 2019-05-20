package main

import (
	"fmt"
	"infoclust/json_io"
	"infoclust/stem"
	"reflect"
)

// 1. Load article and create bow map of type map[string]int from its keywords
// 2. Load subpages file and foreach main category (outer loop) and foreach sub category (inner loop)
//	  create bow map of type map[string]int,
// 3. Store each sub map in a slice of maps
// 4. Compare article's bow map against each slice of the subpages map slice
func compare(bowA, bowB map[string]int) (error, float64) {

	return nil, 0.0
}

func main() {
	//slice1 := []interface{}{"1", "2", "3", "4", "7"}
	//slice2 := []interface{}{"1", "4", "5", "7", "9"}
	//
	//fmt.Printf("%v\n", jaccard.Distance(slice1, slice2))
	//
	//translate.Translate("in_extracted_keywords.txt", "in_bow_translated.json", "out.json")

	mArticle, err := json_io.ReadJSON("test_article.json")

	if err != nil {
		panic("Article file corrupt")
	}

	articleKeywords, ok := mArticle[0]["keywords"].(map[string]interface{})
	if !ok {
		panic("Illegal article file")
	}

	articleBow := stem.Lemmatize(articleKeywords)

	mSubpages, err := json_io.ReadJSON("test_subpages.json")

	if err != nil {
		panic("Subpages file corrupt")
	}

	for _, value := range mSubpages[0] {

		if reflect.ValueOf(value).Kind() == reflect.Map {
			/* Nested map */
			switch value.(type) {
			case map[string]interface{}:
				for _, bow := range value.(map[string]interface{}) {

					bowConverted := make(map[string]int)

					switch bow.(type) {
					case map[string]interface{}:
						for k, v := range bow.(map[string]interface{}) {
							bowConverted[k] = int(v.(float64))
						}
					}

					// TODO: Compare bow with article(s)
					if ok {
						err, dist := compare(articleBow, bowConverted)
						if err != nil {
							panic("Illegal comparison")
						}
						fmt.Println(dist)
					}

				}
			}

		}
	}
}
