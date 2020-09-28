package testhelper

import (
	"log"
	"os"
	"strings"
)

// GetRootPath => If testing the path needs adjusting...
func GetRootPath() (path string, testing bool) {
	path, theErr := os.Getwd()
	if theErr != nil {
		log.Println(theErr)
	}
	appRoot := "gocustomerapp"
	index := strings.Index(path, appRoot)
	// Take substring from index to length of string
	substring := path[index+len(appRoot) : len(path)]
	pathToDotEnv := strings.Replace(path, substring, "", 1)
	return pathToDotEnv, len(substring) > 0
}
