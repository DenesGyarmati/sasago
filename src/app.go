package src

import (
	"bytes"
	"encoding/csv"
	"encoding/xml"
	"fmt"
)

func CalculateSASA(fileContent []byte, pdbName string, email string, advanced bool, classifier string, algorithm string, parameter int, radius float64, depth string, aa string, outputFormat string) (Tree, []byte, error) {

	var sasa []float64

	if pdbName != "" {
		var err error
		fileContent, err = pdbDownload(pdbName)
		if err != nil {
			return Tree{}, fileContent, fmt.Errorf("failed to download PDB: %v", err)
		}
	}

	fileStr := string(fileContent)

	atoms, grid, tree := GetPDBStructure(fileStr, classifier, radius)
	// Call your calculation function
	if algorithm == "LR" {
		sasa = lrAlg(atoms, grid, parameter)
	} else {
		sasa = srAlg(atoms, grid, parameter)
	}

	result := createStructure(sasa, tree, depth, aa)
	var outputFile []byte
	if outputFormat == "CSV" {
		outputFile = printResCSV(result, depth)
	} else if outputFormat == "XML" {
		outputFile, _ = printResXML(result)
	}

	return result, outputFile, nil
}

func printResCSV(result Tree, depth string) []byte {
	var buffer bytes.Buffer
	writer := csv.NewWriter(&buffer)
	switch depth {
	case "All":
		writer.Write([]string{"CHAIN", "CHAIN_SASA", "RES", "RES_INDEX", "RES_SASA", "ATOM_ID", "ATOM_NAME", "ATOM_SASA", "ATOM_SYMBOL"})
		for _, chain := range result.Chains {
			for _, res := range chain.Res {
				for _, atom := range res.Atoms {
					writer.Write([]string{
						chain.Chain,
						fmt.Sprintf("%.3f", chain.SASA),
						res.Res,
						fmt.Sprintf("%d", res.ResI),
						fmt.Sprintf("%.3f", res.SASA),
						fmt.Sprintf("%d", atom.ID),
						atom.Name,
						fmt.Sprintf("%.3f", atom.SASA),
						atom.Symbol,
					})
				}
			}
		}
	case "Residue":
		writer.Write([]string{"CHAIN", "CHAIN_SASA", "RES", "RES_INDEX", "RES_SASA"})
		for _, chain := range result.Chains {
			for _, res := range chain.Res {
				writer.Write([]string{
					chain.Chain,
					fmt.Sprintf("%.3f", chain.SASA),
					res.Res,
					fmt.Sprintf("%d", res.ResI),
					fmt.Sprintf("%.3f", res.SASA),
				})
			}
		}
	default:
		writer.Write([]string{"CHAIN", "CHAIN_SASA"})
		for _, chain := range result.Chains {
			writer.Write([]string{
				chain.Chain,
				fmt.Sprintf("%.3f", chain.SASA),
			})
		}
	}

	writer.Flush()
	return buffer.Bytes()
}

func printResXML(result Tree) ([]byte, error) {
	xmlTree := XMLTree{
		Title:  trimSymbol(result.Title),
		SASA:   result.SASA,
		Chains: []XMLChain{},
	}

	for _, chain := range result.Chains {
		xmlChain := XMLChain{
			Chain: chain.Chain,
			SASA:  chain.SASA,
			Res:   []XMLRes{},
		}
		for _, res := range chain.Res {
			xmlRes := XMLRes{
				Res:   res.Res,
				ResI:  res.ResI,
				SASA:  res.SASA,
				Atoms: []XMLAtom{},
			}
			for _, atom := range res.Atoms {
				xmlAtom := XMLAtom{
					ID:     atom.ID,
					Name:   atom.Name,
					SASA:   atom.SASA,
					Symbol: atom.Symbol,
				}
				xmlRes.Atoms = append(xmlRes.Atoms, xmlAtom)
			}
			xmlChain.Res = append(xmlChain.Res, xmlRes)
		}
		xmlTree.Chains = append(xmlTree.Chains, xmlChain)
	}

	output, err := xml.MarshalIndent(xmlTree, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error marshalling XML: %w", err)
	}
	return output, nil
}
