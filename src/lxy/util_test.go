package main 

import (
	"testing"
	//"fmt"
)

func TestPostfixPath(t *testing.T) {
	//t.Parallel()
	//fmt.Println("testing: util/postfixPath")

	tag := "new"

	paths := map[string]string{
		"/some/path/to/file.fa": "/some/path/to/file_new.fa",
		"some/path/to/file.fa": "some/path/to/file_new.fa",
		"/some/path/to/file.fa.2": "/some/path/to/file_new.fa.2",
		"/some/path/to/file": "/some/path/to/file_new", 
		"filename": "filename_new",
		"filename.fa": "filename_new.fa",
	}

	for pOrig, pExp := range paths {
		pObs := PostfixPath(pOrig, tag)
		if pObs != pExp {
			t.Errorf("test of postfixPath failed on case %s, observing %s and expecting %s", pOrig, pObs, pExp)
		}
	}

}

