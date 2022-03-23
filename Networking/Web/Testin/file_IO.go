package main

import (
	"errors"
	"html/template"
	"log"
	"os"
)

const (
	html  int = 0
	css   int = 1
	js    int = 2
	other int = 3
)

//checks all the existing files in the wiki_pages dir
func getExistingPages() []template.HTML {
	files, err := os.ReadDir("wiki_pages")
	if err != nil {
		log.Fatal(err)
	}

	var returnStrSlice []template.HTML
	for _, file := range files {
		tempFileName := file.Name()[:len(file.Name())-4]
		returnStrSlice = append(returnStrSlice, template.HTML("<li><a href=\"/view/"+tempFileName+"\">"+tempFileName+"</a></li>"))
	}
	return returnStrSlice
}

// >> used in the sub-functions (load_html, load_js...)
// note that file_ext will only be used if using file_type of other
func load_file(file_type int, title string, file_ext string) ([]byte, error) {
	file_dir := ""
	file_extention := ""
	switch file_type {
	case html:
		file_dir = "html/"
		file_extention = ".html"
	case css:
		file_dir = "css/"
		file_extention = ".css"
	case js:
		file_dir = "js/"
		file_extention = ".js"
	case other:
		file_dir = "file/"
		file_extention = file_ext
	default:
		return nil, errors.New("file type unknown")
	}

	file_name := file_dir + title + file_extention

	// >> reading file data & error checking
	data, err := os.ReadFile(file_name)
	if err != nil {
		return nil, err
	}

	return data, nil
}
