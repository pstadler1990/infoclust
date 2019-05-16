package main

import (
	"fmt"
	"infoclust/json_io"
	"infoclust/stem"
)

// 1. Load article and create bow map of type map[string]int from its keywords
// 2. Load subpages file and foreach main category (outer loop) and foreach sub category (inner loop)
//	  create bow map of type map[string]int,
// 3. Store each sub map in a slice of maps
// 4. Compare article's bow map against each slice of the subpages map slice

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

	fmt.Println(articleBow)
}
