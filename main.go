package main

import (
	"fmt"
	"infoclust/jaccard"
	"infoclust/translate"
)

func main() {
	/* Preserve order of json: https://gitlab.com/c0b/go-ordered-json */

	slice1 := []interface{}{"1", "2", "3", "4", "7"}
	slice2 := []interface{}{"1", "4", "5", "7", "9"}

	fmt.Printf("%v\n", jaccard.Distance(slice1, slice2))

	translate.Translate("in_extracted_keywords.txt", "in_bow_translated.json")
}
