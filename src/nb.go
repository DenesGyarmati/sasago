package src

import (
	"sync"
)
var contacts [][]int
var atoms [][]float64

func checkRadius(at1, at2 []float64) bool {
	dx := at1[0] - at2[0]
	dy := at1[1] - at2[1]
	dz := at1[2] - at2[2]
	rs := (at1[3] + at2[3]) * (at1[3] + at2[3])
	cent := dx*dx + dy*dy + dz*dz
	return rs > cent
}

func worker(nb []int, nba []int, wg *sync.WaitGroup) {
	defer wg.Done()
	comb := make([][2]int, 0)

	lb := len(nb)
	lba := len(nba)

	if lb > 1 {
		for cx := 0; cx < lb; cx++ {
			for cy := cx + 1; cy < lb; cy++ {
				comb = append(comb, [2]int{nb[cx], nb[cy]})
			}
		}
		if lba > 0 {
			for cx := 0; cx < lb; cx++ {
				for cy := 0; cy < lba; cy++ {
					comb = append(comb, [2]int{nb[cx], nba[cy]})
				}
			}
		}
	} else {
		for cy := 0; cy < lba; cy++ {
			comb = append(comb, [2]int{nb[0], nba[cy]})
		}
	}

	for _, c := range comb {
		conn := checkRadius(atoms[c[0]], atoms[c[1]])
		if conn {
			contacts[c[0]] = append(contacts[c[0]], c[1])
			contacts[c[1]] = append(contacts[c[1]], c[0])
		}
	}
}

func newNB(atms [][]float64, grid [][][]Cell) [][]int {
	atoms = atms
	contacts = make([][]int, len(atoms))
	for i := range contacts {
		contacts[i] = []int{}
	}

	var wg sync.WaitGroup

	lx, ly, lz := len(grid), len(grid[0]), len(grid[0][0])
	for x := 0; x < lx; x++ {
		for y := 0; y < ly; y++ {
			for z := 0; z < lz; z++ {
				if grid[x][y][z].indexList != nil {
					nb := grid[x][y][z]
					nba := []int{}

					// Collect neighboring indices (with bounds checking)
					for dx := -1; dx <= 1; dx++ {
						for dy := -1; dy <= 1; dy++ {
							for dz := -1; dz <= 1; dz++ {
								if dx == 0 && dy == 0 && dz == 0 {
									continue
								}
								nx, ny, nz := x+dx, y+dy, z+dz
								if nx >= 0 && nx < lx && ny >= 0 && ny < ly && nz >= 0 && nz < lz {
									if grid[nx][ny][nz].indexList != nil {
										nba = append(nba, grid[nx][ny][nz].indexList...)
									}
								}
							}
						}
					}

					wg.Add(1)
					go worker(nb.indexList, nba, &wg)
				}
			}
		}
	}

	wg.Wait()
	return contacts
}