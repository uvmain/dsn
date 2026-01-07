package io

import (
	"dsn/core/config"
	"log"
)

func CreateDirs() {

	dirs := []struct {
		path string
		msg  string
	}{
		{config.UploadsDirectory, "Uploads folder already exists"},
		{config.DatabaseDirectory, "Database folder already exists"},
	}

	for _, d := range dirs {
		if FileExists(d.path) {
			log.Println(d.msg)
		} else {
			CreateDir(d.path)
		}
	}
}
