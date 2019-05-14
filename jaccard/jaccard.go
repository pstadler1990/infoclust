package jaccard

import (
	mapset "github.com/deckarep/golang-set"
)

func Distance(slice1, slice2 []interface{}) float32 {
	/* Convert two given slices into a set and return their jaccard coefficient */
	set1 := mapset.NewSetFromSlice(slice1)
	set2 := mapset.NewSetFromSlice(slice2)

	/* Jaccard distance is defined as the size of the intersection divided by the size of the union */
	return float32(set1.Intersect(set2).Cardinality()) / float32(set1.Union(set2).Cardinality())
}