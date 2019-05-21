package cosine

import "github.com/atedja/go-vector"

func Distance(vectorA, vectorB vector.Vector) (error, float64) {
	dot, err := vector.Dot(vectorA, vectorB)
	if err != nil {
		return err, 0.0
	}
	return nil, float64(dot / (vectorA.Magnitude() * vectorB.Magnitude()))
}
