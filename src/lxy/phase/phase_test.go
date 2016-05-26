package phase

import (
	"os"
	"path/filepath"
	util "sequtil"
	"testing"
)

func TestPhase(t *testing.T) {

	cases := [][]string{
		{filepath.Join(cwd(t), "_testdata", "test1.links"), filepath.Join(cwd(t), "_testdata", "test1.key")},
		{filepath.Join(cwd(t), "_testdata", "test2.links"), filepath.Join(cwd(t), "_testdata", "test2.key")},
	}

	testOutPath := filepath.Join(cwd(t), "_testdata", "dummy.out")

	for _, v := range cases {

		links, _ := util.LoadLinks(v[0])
		scaffolding := Phase(&links)
		key := ReadPhasing(v[1])
		score, nscore, _ := EvalScaffolding(scaffolding, key)

		if score != 1 {
			t.Errorf("error")
		}

	}

}

func cwd(t *testing.T) string {
	cwd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	return cwd
}
