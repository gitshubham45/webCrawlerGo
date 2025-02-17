package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type ResultFile struct {
	Domain      string    `json:"domain"`
	URLs        []string  `json:"urls"`
	CreatedTime time.Time `json:"createdTime"`
	UpdatedTime time.Time `json:"updatedTime"`
}

func SaveResultsToFile(domain string, urls []string) {
	if urls == nil {
		urls = []string{} 
	}

	resultsFolder := "results"

	err := os.MkdirAll(resultsFolder, os.ModePerm)
	if err != nil {
		log.Printf("Failed to create folder %s: %v", resultsFolder, err)
		return
	}

	fileName := filepath.Join(resultsFolder, fmt.Sprintf("%s.json", domain))

	var resultFile ResultFile
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Printf("Failed to open file %s: %v", fileName, err)
		return
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	if fileInfo.Size() > 0 {
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&resultFile); err != nil {
			log.Printf("Failed to decode existing file %s: %v", fileName, err)
			return
		}
	} else {
		resultFile = ResultFile{
			Domain:      domain,
			CreatedTime: time.Now(),
		}
	}

	for _, url := range urls {
		if !contains(resultFile.URLs, url) {
			resultFile.URLs = append(resultFile.URLs, url)
		}
	}

	resultFile.UpdatedTime = time.Now()

	file.Truncate(0) 
	file.Seek(0, 0)  
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") 
	if err := encoder.Encode(resultFile); err != nil {
		log.Printf("Failed to write results to file %s: %v", fileName, err)
		return
	}

	log.Printf("Saved/updated results for domain %s to file %s", domain, fileName)
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
