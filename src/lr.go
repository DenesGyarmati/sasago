package src

import (
	"math"
	"sort"
	"sync"
)

const TwoPi = 2 * math.Pi

// lrAlg calculates SASA using multiple goroutines for concurrent atom processing.
func lrAlg(atoms [][]float64, grid [][][]Cell, slicesVal int) []float64 {
	sasa := make([]float64, len(atoms))
	var wg sync.WaitGroup
	adj := newNB(atoms, grid)

	// Concurrent atom area calculation
	for atomIdx, neighbors := range adj {
		wg.Add(1)
		go func(idx int, neighbors []int) {
			defer wg.Done()
			calculateAtomArea(idx, neighbors, atoms, slicesVal, sasa)
		}(atomIdx, neighbors)
	}

	wg.Wait()
	return sasa
}

// calculateAtomArea determines the exposed surface area of an atom
func calculateAtomArea(atomIdx int, neighbors []int, atoms [][]float64, slicesVal int, sasa []float64) {
	radius := atoms[atomIdx][3]
	zCoord := atoms[atomIdx][2]
	delta := 2 * radius / float64(slicesVal)
	z := zCoord - radius - 0.5*delta
	totalArea := 0.0
	xCoord := atoms[atomIdx][0]
	yCoord := atoms[atomIdx][1]

	for i := 0; i < slicesVal; i++ {
		z += delta
		distanceZ := math.Abs(zCoord - z)
		radiusSquaredPrime := radius*radius - distanceZ*distanceZ
		if radiusSquaredPrime < 0 {
			continue
		}
		radiusPrime := math.Sqrt(radiusSquaredPrime)
		if radiusPrime <= 0 {
			continue
		}

		isBuried := false
		arcs := make([][2]float64, 0)

		for _, neighborIdx := range neighbors {
			neighborZ := atoms[neighborIdx][2]
			neighborDistZ := math.Abs(neighborZ - z)
			neighborRadius := atoms[neighborIdx][3]

			if neighborDistZ < neighborRadius {
				neighborRadiusSquaredPrime := neighborRadius*neighborRadius - neighborDistZ*neighborDistZ
				neighborRadiusPrime := math.Sqrt(neighborRadiusSquaredPrime)
				dx := xCoord - atoms[neighborIdx][0]
				dy := yCoord - atoms[neighborIdx][1]
				distance := math.Sqrt(dx*dx + dy*dy)

				if distance >= radiusPrime+neighborRadiusPrime {
					continue
				}
				if distance+radiusPrime < neighborRadiusPrime {
					isBuried = true
					break
				}
				if distance+neighborRadiusPrime < radiusPrime {
					continue
				}

				alpha := math.Acos((radiusSquaredPrime + distance*distance - neighborRadiusSquaredPrime) / (2.0 * radiusPrime * distance))
				beta := math.Atan2(dy, dx) + math.Pi
				start := beta - alpha
				end := beta + alpha

				if start < 0 {
					start += TwoPi
				}
				if end > TwoPi {
					end -= TwoPi
				}
				if end < start {
					arcs = append(arcs, [2]float64{0, end}, [2]float64{start, TwoPi})
				} else {
					arcs = append(arcs, [2]float64{start, end})
				}
			}
		}

		if !isBuried {
			totalArea += delta * radius * calculateExposedArc(arcs)
		}
	}

	sasa[atomIdx] = totalArea
}

// calculateExposedArc computes the exposed arc from a list of arc intervals.
func calculateExposedArc(arcs [][2]float64) float64 {
	if len(arcs) == 0 {
		return TwoPi
	}
	sort.Slice(arcs, func(i, j int) bool {
		return arcs[i][0] < arcs[j][0]
	})

	exposedSum := arcs[0][0]
	maxSup := arcs[0][1]

	for i := 1; i < len(arcs); i++ {
		if maxSup < arcs[i][0] {
			exposedSum += arcs[i][0] - maxSup
		}
		if arcs[i][1] > maxSup {
			maxSup = arcs[i][1]
		}
	}

	return exposedSum + TwoPi - maxSup
}