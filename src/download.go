package src

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
)

func pdbDownload(file string) ([]byte, error) {
	url := fmt.Sprintf("https://files.rcsb.org/download/%s.pdb.gz", file)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Failed to download:", err)
		return nil, err
	}
	defer resp.Body.Close()
	reader, err := gzip.NewReader(resp.Body)
	if err != nil {
		log.Println("Failed to create gzip reader:", err)
		return nil, err
	}
	defer reader.Close()
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, reader); err != nil {
		log.Println("Failed to read decompressed content:", err)
		return nil, err
	}

	fmt.Printf("%s retrieved successfully.\n", file)
	return buf.Bytes(), nil
}
