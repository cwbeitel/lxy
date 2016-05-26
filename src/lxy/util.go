package main

import (
	"fmt"
	"os"
	"strings"
)

// postfixPath takes a path and inserts a given tag string before the last period, after the last /,
// e.g. path/to/file/somefilename.fa -> path/to/file/somefilename_new.fa
func PostfixPath(path string, tag string) string {

	// if the tag is null, return string
	if len(path) == 0 {
		return path
	}

	arr := strings.Split(path, "/")
	filenameArray := strings.Split(arr[(len(arr)-1)], ".")
	filenameArray[0] = filenameArray[0] + "_" + tag

	// set the last string in the path to the new string
	arr[(len(arr) - 1)] = strings.Join(filenameArray, ".")

	return strings.Join(arr, "/")

}

func PathExists(path string) {
	if _, err := os.Stat(path); err != nil {
		fmt.Errorf("error: cannot find the specified file path on your system: %s\n", path)
	}
}

func MkdirForFile(path string) {
	arr := strings.Split(path, "/")
	pathNew := strings.Join(arr[0:(len(arr)-1)], "/")
	os.MkdirAll(pathNew, 0777)

}
