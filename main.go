package main

import (
	"fmt"
	"github.com/atedja/go-vector"
	"github.com/deckarep/golang-set"
	"github.com/logrusorgru/aurora"
	"infoclust/cosine"
	"infoclust/json_io"
	"infoclust/stem"
	"reflect"
)

const MIN_SCORE float64 = 0.75

func compare(bowA, bowB map[string]int) (error, float64) {
	// Convert the two bow's into vectors in the form [0,1,0,1,1,...]
	// Missing words in the shorter bow must be represented as 0
	var longest, shortest map[string]int

	if len(bowA) >= len(bowB) {
		longest = bowA
		shortest = bowB
	} else {
		longest = bowB
		shortest = bowA
	}

	tmpVecA := make([]float64, 0, len(longest))
	tmpVecB := make([]float64, 0, len(longest))

	for k, countOuter := range longest {
		countInner, ok := shortest[k]
		tmpVecA = append(tmpVecA, float64(countOuter))
		if ok {
			// word k is in both maps
			tmpVecB = append(tmpVecB, float64(countInner))
		} else {
			// word k is not present in the shorter slice, so set the count to 0
			tmpVecB = append(tmpVecB, 0)
		}
	}

	outVecA := vector.NewWithValues(tmpVecA)
	outVecB := vector.NewWithValues(tmpVecB)

	return cosine.Distance(outVecA, outVecB)
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

	mSubpages, err := json_io.ReadJSON("test_subpages.json")

	if err != nil {
		panic("Subpages file corrupt")
	}

	// Keeps a set of all results
	// m[article_name] -> [category_name_a, category_name_b,...]
	mSummarize := make(map[string]mapset.Set)

	for _, article := range mArticle {

		articleKeywords, ok := article["keywords"].(map[string]interface{})
		if !ok {
			panic("Illegal article file")
		}

		articleBow := make(map[string]int)

		for k, v := range articleKeywords {
			articleBow[k] = int(v.(float64))
		}
		articleBow = stem.Lemmatize(articleBow)

		fmt.Println("Distance between ", aurora.Magenta(article["title"]), " and...")

		name := article["title"].(string)

		// Inner loop to cross each article's bow with each of the subpages file
		for cat, value := range mSubpages[0] {

			if reflect.ValueOf(value).Kind() == reflect.Map {
				/* Nested map */
				switch value.(type) {
				case map[string]interface{}:
					for sub, bow := range value.(map[string]interface{}) {

						bowConverted := make(map[string]int)

						switch bow.(type) {
						case map[string]interface{}:
							for k, v := range bow.(map[string]interface{}) {
								bowConverted[k] = int(v.(float64))
							}
						}

						bowConverted = stem.Lemmatize(bowConverted)

						err, dist := compare(articleBow, bowConverted)
						if err != nil {
							panic("Illegal comparison")
						}

						if dist >= MIN_SCORE {
							_, ok := mSummarize[name]
							if !ok {
								// Article does not exist yet, allocate new set
								mSummarize[name] = mapset.NewSet()
							}
							mSummarize[name].Add(cat)
							fmt.Println(aurora.Red(sub), "from category ", aurora.Blue(cat), ": ", dist)
						}
					}
				}
			}
		}

		// Summarize detected categories (main categories)
		fmt.Println(mSummarize[name])

		// TODO: Write these sets into a json object to disk
	}
}
