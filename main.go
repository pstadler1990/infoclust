package main

import (
	"fmt"
	"infoclust/json_io"
)

func main() {
	//slice1 := []interface{}{"1", "2", "3", "4", "7"}
	//slice2 := []interface{}{"1", "4", "5", "7", "9"}
	//
	//fmt.Printf("%v\n", jaccard.Distance(slice1, slice2))
	//
	//translate.Translate("in_extracted_keywords.txt", "in_bow_translated.json", "out.json")

	mArticle, err := json_io.ReadJSON("test_article.json")

	if err != nil {
		return
	}

	fmt.Println(mArticle)
}
