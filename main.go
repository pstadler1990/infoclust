package main

import (
	"fmt"
	"github.com/atedja/go-vector"
	mapset "github.com/deckarep/golang-set"
	"github.com/logrusorgru/aurora"
	"github.com/pstadler1990/infoclust/cosine"
	"github.com/pstadler1990/infoclust/json_io"
	"github.com/pstadler1990/infoclust/stem"
	"log"
	"os"
	"reflect"
	"runtime"
	"sync"
)

const MIN_SCORE float64 = 0.85
const LOG_FILE string = "results.log"
const IN_ARTICLES_FILE string = "out.json"
const IN_SUBPAGES_FILE string = "jsonformatter.json"
const WORKERS int = 32

var wg sync.WaitGroup

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

func calculateBowPerArticle(subpages map[string]interface{},
	in <-chan map[string]interface{},
	out chan<- map[string]mapset.Set) {
	// Goroutine-able function to calculate the cosine distance for each bow model
	// in the given article and the subpages file

	for article := range in {

		name := article["title"].(string)
		// Keeps a local set of the current results
		// m[article_name] -> [category_name_a, category_name_b,...]
		mSummarize := make(map[string]mapset.Set)
		mSummarize[name] = mapset.NewSet()

		articleKeywords, ok := article["keywords"].(map[string]interface{})
		if !ok {
			panic("Illegal article file")
		}

		articleBow := make(map[string]int)

		for k, v := range articleKeywords {
			articleBow[k] = int(v.(float64))
		}
		articleBow = stem.Lemmatize(articleBow)

		fmt.Println(aurora.Red("Article received"), article)

		// Inner loop to cross each article's bow with each of the subpages file
		for cat, value := range subpages {

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
							mSummarize[name].Add(cat)
							log.Println(sub, "from category ", cat, ": ", dist, "from", name)
						}
					}
				}
			}
		}

		fmt.Println(aurora.Blue("Finished worker"), article["title"])
		out <- mSummarize
		wg.Done()
	}
}

func main() {
	//slice1 := []interface{}{"1", "2", "3", "4", "7"}
	//slice2 := []interface{}{"1", "4", "5", "7", "9"}
	//
	//fmt.Printf("%v\n", jaccard.Distance(slice1, slice2))
	//
	//translate.Translate("in_extracted_keywords.txt", "in_bow_translated.json", "out.json")

	f, logErr := os.OpenFile(LOG_FILE, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if logErr != nil {
		panic("Could not create log file")
	}

	defer f.Close()

	log.SetOutput(f)

	// out.json contains the translated bow
	mArticle, err := json_io.ReadJSON(IN_ARTICLES_FILE)
	if err != nil {
		panic("Article file corrupt")
	}

	mSubpages, err := json_io.ReadJSON(IN_SUBPAGES_FILE)
	if err != nil {
		panic("Subpages file corrupt")
	}

	// Keeps a set of all results
	// m[article_name] -> [category_name_a, category_name_b,...]
	mSummarize := make(map[string]mapset.Set)

	// Create n channels
	runtime.GOMAXPROCS(runtime.NumCPU())
	jobs := make(chan map[string]interface{}, len(mArticle))
	out := make(chan map[string]mapset.Set, len(mArticle))

	for gr := 0; gr < WORKERS; gr++ {
		go calculateBowPerArticle(mSubpages[0], jobs, out)
	}

	fmt.Println(aurora.BgBlack(aurora.Yellow("Number of articles:")), len(mArticle))
	for _, article := range mArticle {
		jobs <- article
		wg.Add(1)
	}
	close(jobs)

	for m := 0; m < len(mArticle); m++ {
		for k, v := range <-out {
			_, ok := mSummarize[k]
			if !ok {
				// Push into common summarize object
				mSummarize[k] = v
			}
			fmt.Println(aurora.BgYellow(aurora.Black(k)), v)
		}
	}
	close(out)

	wg.Wait()

	log.Println("-- BEGIN OF DUMP --")
	log.Println(mSummarize)
	// TODO: Write these sets into a json object to disk
	fmt.Println(aurora.BgGreen("Finished calculation!"))
	log.Println("-- END OF DUMP --")
}
