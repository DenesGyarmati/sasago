package src

import (
	"math"
	"strings"
)

func radius(rName, aName string) float64 {
	p, exists := data[rName]
	if !exists {
		return -1.0
	}
	r, exists := p[aName]
	if !exists {
		return -1.0
	}
	return r
}

func InsertInTree(tree *Tree, a *Node) {
	inChain := false
	inRes := false

	for i := range tree.Chains {
		if tree.Chains[i].Chain == a.ChainLabel {
			inChain = true
			for j := range tree.Chains[i].Res {
				if tree.Chains[i].Res[j].ResI == a.ResIndex {
					tree.Chains[i].Res[j].Atoms = append(tree.Chains[i].Res[j].Atoms, Atom{ID: a.ID, Name: a.AtomName, SASA: 0, Symbol: strings.ToLower(a.Symbol)})
					inRes = true
					break
				}
			}
			if !inRes {
				// If the residue wasn't found, add it
				newRes := Res{Res: a.ResName, ResI: a.ResIndex, SASA: a.SASA}
				newRes.Atoms = append(newRes.Atoms, Atom{ID: a.ID, Name: a.AtomName, SASA: 0, Symbol: strings.ToLower(a.Symbol)})
				tree.Chains[i].Res = append(tree.Chains[i].Res, newRes)
			}
			break
		}
	}

	if !inChain {
		newChain := Chain{Chain: a.ChainLabel, SASA: a.SASA}
		newRes := Res{Res: a.ResName, ResI: a.ResIndex, SASA: a.SASA}
		newRes.Atoms = append(newRes.Atoms, Atom{ID: a.ID, Name: a.AtomName, SASA: 0, Symbol: strings.ToLower(a.Symbol)})
		newChain.Res = append(newChain.Res, newRes)
		tree.Chains = append(tree.Chains, newChain)
	}
}

func (n *Node) CheckAtomRadius(pr float64) {
	r := radius(n.ResName, n.AtomName)
	if r == -1 {
		r = radius("ANY", n.AtomName)
	}
	if r == -1 {
		r = radius("DEFAULT", n.Symbol)
	}
	n.Radius = r + pr
}

var data map[string]map[string]float64

func SetRadiusFile(c string) {
	switch c {
	case "NACCESS":
		data = GetNaccess()
	case "PROTOR":
		data = GetProtor()
	default:
		data = GetOons()
	}
}

func Radius(rName, aName string) float64 {
	p, ok := data[rName]
	if !ok {
		return -1.0
	}
	r, ok := p[aName]
	if !ok {
		return -1.0
	}
	return r
}

func (c *Cell) InsertIndex(i int) {
	if c.indexList[0] == -1 {
		c.indexList[0] = i
	} else {
		c.indexList = append(c.indexList, i)
	}
}

func GetPDBStructure(file, classifier string, pr float64) ([][]float64, [][][]Cell, Tree) {
	SetRadiusFile(classifier)
	tree := Tree{}
	minX, minY, minZ := math.Inf(1), math.Inf(1), math.Inf(1)
	maxX, maxY, maxZ := math.Inf(-1), math.Inf(-1), math.Inf(-1)
	maxRadi := 0.0
	countAtoms := 0
	var structure [][]float64

	for _, line := range strings.Split(string(file), "\n") {
		if strings.HasPrefix(line, "TITLE") {
			tree.Title += line[9:] + " "
		}
		if strings.HasPrefix(line, "ATOM") && strings.TrimSpace(line[76:78]) != "H" {
			countAtoms++
			a := NewNode(line)
			a.CheckAtomRadius(pr)
			structure = append(structure, []float64{a.Coord[0], a.Coord[1], a.Coord[2], a.Radius})

			for _, val := range a.Coord {
				minX = math.Min(minX, val)
				maxX = math.Max(maxX, val)
				minY = math.Min(minY, val)
				maxY = math.Max(maxY, val)
				minZ = math.Min(minZ, val)
				maxZ = math.Max(maxZ, val)
			}

			if a.Radius > maxRadi {
				maxRadi = a.Radius
			}
			InsertInTree(&tree, a)
		}
	}

	dim := maxRadi * 2
	x := int((maxX-minX)/dim) + 1
	y := int((maxY-minY)/dim) + 1
	z := int((maxZ-minZ)/dim) + 1

	grid := make([][][]Cell, x)
	for i := range grid {
		grid[i] = make([][]Cell, y)
		for j := range grid[i] {
			grid[i][j] = make([]Cell, z)
		}
	}

	for i, at := range structure {
		atX := int((at[0] - minX) / dim)
		atY := int((at[1] - minY) / dim)
		atZ := int((at[2] - minZ) / dim)

		if grid[atX][atY][atZ].indexList == nil {
			grid[atX][atY][atZ] = Cell{indexList: []int{-1}}
		}
		grid[atX][atY][atZ].InsertIndex(i)
	}
	return structure, grid, tree
}

func createStructure(sasa []float64, tree Tree, depth string, aa string) Tree {
	index := 0
	for i := range tree.Chains {
		chain := &tree.Chains[i]
		for j := range chain.Res {
			residue := &chain.Res[j]
			for k := range residue.Atoms {
				atom := &residue.Atoms[k]
				if sasa[index] > 0 {
					atom.SASA = trimTrailingZeros(sasa[index])
					residue.SASA += sasa[index]
					chain.SASA += sasa[index]
					tree.SASA += sasa[index]
				}
				index++
			}
		}
	}
	if depth == "Chain" {
		for i := range tree.Chains {
			chain := &tree.Chains[i]
			chain.SASA = trimTrailingZeros(chain.SASA)
			chain.Res = []Res{}
		}
		return tree
	}
	for i := range tree.Chains {
		chain := &tree.Chains[i]
		chain.SASA = trimTrailingZeros(chain.SASA)
		for j := 0; j < len(chain.Res); {
			residue := &chain.Res[j]
			if depth != "Chain" && (aa == "All" || aa == residue.Res) {
				residue.SASA = trimTrailingZeros(residue.SASA)
				if depth == "Residue" {
					residue.Atoms = []Atom{}
				}
				j++
			} else {
				chain.Res = append(chain.Res[:j], chain.Res[j+1:]...)
			}
		}
	}
	return tree
}
