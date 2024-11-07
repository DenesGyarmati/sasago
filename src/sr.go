package src

import (
	"math"
	"sync"
)
var ptVal int
var sasa []float64
// alg initializes the atoms and calculates the SASA using concurrency.
func srAlg(atms [][]float64, grid [][][]Cell, pointsVal int) []float64 {
	atoms = atms
	ptVal = pointsVal
	atomsLen := len(atoms)
	sasa = make([]float64, atomsLen)

	adj := newNB(atoms, grid)
	var wg sync.WaitGroup

	for x, neighbors := range adj {
		wg.Add(1)
		go func(x int, neighbors []int) {
			defer wg.Done()
			srAtomArea(x, neighbors)
		}(x, neighbors)
	}
	wg.Wait()
	return sasa
}

// pointTransform modifies the points based on the given xyz and radius r.
func pointTransform(points [][]float64, xyz []float64, r float64) {
	for pi := range points {
		points[pi][0] = points[pi][0]*r + xyz[0]
		points[pi][1] = points[pi][1]*r + xyz[1]
		points[pi][2] = points[pi][2]*r + xyz[2]
	}
}

func srAtomArea(x int, neighbors []int) {
	nSurface := 0
	at := atoms[x]
	r1 := at[3]
	r2 := r1 * r1
	tp := testPoints(ptVal)
	pointTransform(tp, []float64{at[0], at[1], at[2]}, r1)

	lTp := len(tp)
	lA := len(neighbors)

	for j := 0; j < lTp; j++ {
		buried := false
		for n := 0; n < lA; n++ {
			dx := tp[j][0] - atoms[neighbors[n]][0]
			dy := tp[j][1] - atoms[neighbors[n]][1]
			dz := tp[j][2] - atoms[neighbors[n]][2]
			ri := atoms[neighbors[n]][3] * atoms[neighbors[n]][3]
			if dx*dx+dy*dy+dz*dz <= ri {
				buried = true
				break
			}
		}
		if !buried {
			nSurface++
		}
	}
	sasa[x] = (4.0 * math.Pi * r2 * float64(nSurface)) / float64(ptVal)
}

// testPoints generates test points for a given number of points.
func testPoints(n int) [][]float64 {
	dlong := math.Pi * (3 - math.Sqrt(5))
	dz := 2.0 / float64(n)
	longitude := 0.0
	z := 1.0 - dz/2
	testPArray := make([][]float64, 0, n)

	for i := 0; i < n; i++ {
		r := math.Sqrt(1 - z*z)
		p := []float64{
			math.Cos(longitude) * r,
			math.Sin(longitude) * r,
			z,
		}
		testPArray = append(testPArray, p)
		z -= dz
		longitude += dlong
	}
	return testPArray
}