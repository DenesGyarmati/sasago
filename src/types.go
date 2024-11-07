package src

import (
	"encoding/xml"
)

type Atom struct {
	ID   int     `json:"id"`
	Name string  `json:"name"`
	SASA float64 `json:"sasa"`
	// Radius float64 `json:"radius"`
	Symbol string `json:"symbol"`
}

type Res struct {
	Res   string  `json:"res"`
	ResI  int     `json:"res_index"`
	SASA  float64 `json:"sasa"`
	Atoms []Atom  `json:"atoms"`
}

type Chain struct {
	Chain string  `json:"chain"`
	SASA  float64 `json:"sasa"`
	Res   []Res   `json:"residues"`
}

type Tree struct {
	DownloadURL string  `json:"download"`
	Title       string  `json:"title"`
	SASA        float64 `json:"sasa"`
	Chains      []Chain `json:"chains"`
}

type Node struct {
	ID         int        `json:"id"`
	AtomName   string     `json:"atom_name"`
	ResName    string     `json:"res_name"`
	ChainLabel string     `json:"chain_label"`
	ResIndex   int        `json:"res_index"`
	Coord      [3]float64 `json:"coord"`
	Symbol     string     `json:"symbol"`
	Radius     float64    `json:"radius"`
	Links      []int      `json:"links"`
	SASA       float64    `json:"sasa"`
}

func NewNode(line string) *Node {
	id, _ := getInt(line[6:11])
	resIndex, _ := getInt(line[22:26])
	coord := [3]float64{
		parseFloat(line[31:38]),
		parseFloat(line[39:46]),
		parseFloat(line[47:54]),
	}

	symbol := trimSymbol(line[76:78])
	radius := -1.0
	if symbol == "H" {
		radius = 1.2
	}

	return &Node{
		ID:         id,
		AtomName:   trimSymbol(line[13:16]),
		ResName:    line[17:20],
		ChainLabel: line[21:22],
		ResIndex:   resIndex,
		Coord:      coord,
		Symbol:     symbol,
		Radius:     radius,
		Links:      []int{},
		SASA:       0,
	}
}

type Cell struct {
	indexList []int
}

type XMLAtom struct {
	XMLName xml.Name `xml:"atom"`
	ID      int      `xml:"id,attr"`
	Name    string   `xml:"name,attr"`
	SASA    float64  `xml:"sasa,attr"`
	Symbol  string   `xml:"symbol,attr"`
}

type XMLRes struct {
	XMLName xml.Name  `xml:"residue"`
	Res     string    `xml:"name,attr"`
	ResI    int       `xml:"index,attr"`
	SASA    float64   `xml:"sasa,attr"`
	Atoms   []XMLAtom `xml:"atom"`
}

type XMLChain struct {
	XMLName xml.Name `xml:"chain"`
	Chain   string   `xml:"name,attr"`
	SASA    float64  `xml:"sasa,attr"`
	Res     []XMLRes `xml:"residue"`
}

type XMLTree struct {
	XMLName xml.Name   `xml:"file"`
	Title   string     `xml:"name,attr"`
	SASA    float64    `xml:"sasa,attr"`
	Chains  []XMLChain `xml:"chain"`
}
