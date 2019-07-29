package pkg

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func LocateFiles (searchPath string) []string {
	var files []string

	err := filepath.Walk(searchPath, func(path string, info os.FileInfo, err error) error {
		// Skip over directories, only load files
		if info.IsDir() {
			return nil
		}
		// Load the file
		file, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}
		fileString := string(file)
		// Only add files which contain @Route. Not perfect but this will give us a short list for later validation
		if strings.Contains(fileString, "@Route") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		fmt.Println(file)
	}
	return files
}