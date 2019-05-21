package stem

import (
	"github.com/aaaton/golem"
	"github.com/aaaton/golem/dicts/en"
)

func Lemmatize(bow map[string]int) map[string]int {
	/* Lemmatize all words in the given bag of words (bow),
	   to increase matches */
	lemmatizer, err := golem.New(en.NewPackage())
	if err != nil {
		panic(err)
	}

	lemmatizedBow := make(map[string]int)

	for k, v := range bow {
		word := lemmatizer.Lemma(k)
		lemmatizedBow[word] = v
	}

	return lemmatizedBow
}
