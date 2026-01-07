package io

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func CreateDirectoryIfNotExists(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return os.MkdirAll(dirPath, os.ModePerm)
	}
	return nil
}

func FileExists(absoluteFilePath string) bool {
	cleanPath, err := PathWithoutTraversal(absoluteFilePath)
	if err != nil {
		log.Printf("error checking file existence - path traversal detected: %s", absoluteFilePath)
		return false
	}
	if _, err := os.Stat(cleanPath); os.IsNotExist(err) {
		return false
	}
	return true
}

func CreateDir(directoryPath string) {
	if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
		err := os.MkdirAll(directoryPath, 0755)
		if err != nil {
			log.Printf("Error creating directory %s: %s", directoryPath, err)
		} else {
			log.Printf("Directory created: %s", directoryPath)
		}
	} else {
		log.Printf("Directory already exists: %s", directoryPath)
	}
}

func PathWithoutTraversal(inputPath string) (string, error) {
	if strings.Contains(inputPath, "..") {
		return "", fmt.Errorf("path traversal detected")
	}
	cleanPath := filepath.Clean(inputPath)
	return cleanPath, nil
}
