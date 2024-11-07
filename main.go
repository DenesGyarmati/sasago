package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"sasa/src"
	"strconv"
	"time"
)

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", formHandler)
	http.HandleFunc("/calculate", calculateHandler)
	log.Println("Server started at http://localhost:8888")
	http.ListenAndServe(":8888", nil)
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	formTemplatePath := filepath.Join("templates", "form.html")
	serveForm(w, formTemplatePath)
}

func serveForm(w http.ResponseWriter, formTemplatePath string) {
	w.Header().Set("Content-Type", "text/html")
	tmpl, err := template.ParseFiles(formTemplatePath)
	if err != nil {
		http.Error(w, "Error loading form template", http.StatusInternalServerError)
		log.Println("Template error:", err)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Error rendering form", http.StatusInternalServerError)
		log.Println("Render error:", err)
	}
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // max memory 10MB
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	pdbName := r.FormValue("pdb_name")
	email := r.FormValue("email")
	advanced := r.FormValue("advanced")
	buf := new(bytes.Buffer)
	// outputName = pdbName

	if pdbName == "" {
		file, _, err := r.FormFile("file_upload")
		if err != nil {
			http.Error(w, "Error retrieving file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		_, err = io.Copy(buf, file)
		if err != nil {
			http.Error(w, "Error processing file", http.StatusInternalServerError)
			return
		}
	}

	var adv bool
	if advanced == "on" {
		adv = true
	}

	var (
		classifier   string
		algorithm    string
		param        int
		rad          float64
		depth        string
		aa           string
		outputFormat string
	)

	if adv {
		classifier = r.FormValue("classifier")
		algorithm = r.FormValue("algorithm")
		parameter := r.FormValue("parameter")
		param, err = strconv.Atoi(parameter)
		if err != nil {
			http.Error(w, "Invalid parameter", http.StatusBadRequest)
			return
		}
		radiusStr := r.FormValue("radius")
		rad, err = strconv.ParseFloat(radiusStr, 64)
		if err != nil {
			http.Error(w, "Invalid radius", http.StatusBadRequest)
			return
		}

		depth = r.FormValue("depth")
		aa = r.FormValue("aa")
		outputFormat = r.FormValue("format")
	} else {
		// Default values when 'adv' is false
		classifier = "NACCESS"
		algorithm = "LR"
		param = 20
		rad = 1.4
		depth = "RES"
		aa = "All"
		outputFormat = "only_preview"
	}

	// Call the SASA calculation function
	treeResult, outputFile, err := src.CalculateSASA(buf.Bytes(), pdbName, email, adv, classifier, algorithm, param, rad, depth, aa, outputFormat)
	if err != nil {
		http.Error(w, "Error during calculation", http.StatusInternalServerError)
		return
	}

	if len(outputFile) > 0 {
		var fileExtension string
		if outputFormat == "CSV" {
			fileExtension = ".csv"
		} else if outputFormat == "XML" {
			fileExtension = ".xml"
		}
		timestamp := time.Now().Format("20060102150405")
		downloadFilename := fmt.Sprintf("%s_%s%s", pdbName, timestamp, fileExtension)
		// Save the outputFile to a temporary location
		downloadPath := "/download/" + downloadFilename
		http.HandleFunc(downloadPath, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", downloadFilename))
			if outputFormat == "CSV" {
				w.Header().Set("Content-Type", "text/csv")
			} else if outputFormat == "XML" {
				w.Header().Set("Content-Type", "application/xml")
			}
			w.Write(outputFile)
		})
		// Add download path to the template data
		treeResult.DownloadURL = downloadPath
	}

	// Serve the form with the JSON data included
	formTemplatePath := filepath.Join("templates", "form.html")
	w.Header().Set("Content-Type", "text/html")
	tmpl, err := template.ParseFiles(formTemplatePath)
	if err != nil {
		http.Error(w, "Error loading form template", http.StatusInternalServerError)
		log.Println("Template error:", err)
		return
	}

	// err = tmpl.Execute(w, finalResult) // Pass the populated struct to the template
	err = tmpl.Execute(w, treeResult)
	if err != nil {
		http.Error(w, "Error rendering form", http.StatusInternalServerError)
		log.Println("Render error:", err)
	}
}
