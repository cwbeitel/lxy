package util

import (
	"testing"
	"reflect"
	"path/filepath"
	"fmt"
)

func TestAddKey(t *testing.T) {

	l := NewLinks()

	l.addKey("hi")
	val, ok := l.idKey["hi"]
	if !ok {
		t.Errorf("test error: error adding key to links set")
	}

	if _, ok := l.idKeyRev[val]; !ok {
		t.Errorf("test error: the id for key added via links.addKey() was not found in links.idKeyRev")
	}

}

func TestIntIDs(t *testing.T) {

	l := NewLinks()
	l.addKey("hi")
	l.addKey("hello")
	arr := l.IntIDs()
	key := []int{0,1}
	if !reflect.DeepEqual(arr, key) {
		t.Errorf("test links: the entity ids assigned upon creation of two new entities did not match the expectation, %d", arr)
		fmt.Println(arr)
		fmt.Println(key)
	}

}

func TestLoadLinks(t *testing.T) {

	path := filepath.Join(cwd(t), "_testdata", "toy.ctg.links")
	links, err := LoadLinks(path)
	if err != nil {
		t.Error(err)
	}
	linksKey := NewLinks()
	linksKey.Set(linksKey.ID("chr1"), linksKey.ID("chr2"), 2)
	linksKey.Set(linksKey.ID("chr2"), linksKey.ID("chr3"), 1)

	if !reflect.DeepEqual(links, linksKey) {
		t.Errorf("LoadLinks(%s) yielded links object not matching key", path)
	}

}

func TestStringIDs(t *testing.T) {

	l := NewLinks()
	l.addKey("a")
	l.addKey("b")
	ids := l.StringIDs()
	if !(ids[0] == "a" && ids[1] == "b") {
		t.Errorf("links.StringIDs() returned an array of string IDs different from what was expected.")
		fmt.Println("expected: [a, b]")
		fmt.Println("observed: ", ids)

	}
}

func TestSet(t *testing.T) {

	l := NewLinks()
	l.addKey("hi")
	l.addKey("hello")
	e := l.Set(0, 1, 0.01)
	if e != nil {
		t.Errorf("%s", e)
	}

}

func TestAdd(t *testing.T) {

	l := NewLinks()
	l.addKey("a")
	l.addKey("b")
	l.Set(l.ID("a"),l.ID("b"),0.1)
	l.Add(l.ID("a"),l.ID("b"),0.1)
	val, e := l.Get(l.ID("a"),l.ID("b"),)
	if e != nil {
		t.Error(e)
	}
	if val != 0.2 {
		t.Errorf("error: the value returned from the get operation (%f) did not match the expectation (%f).", val, 0.1)
	}	


}


func TestGet(t *testing.T) {

	l := NewLinks()
	l.addKey("a")
	l.addKey("b")
	l.Set(l.ID("a"),l.ID("b"),0.1)
	val, e := l.Get(l.ID("a"),l.ID("b"),)
	if e != nil {
		t.Error(e)
	}
	if val != 0.1 {
		t.Errorf("error: the value returned from the get operation (%f) did not match the expectation (%f).", val, 0.1)
	}	

}

func TestPrint(t *testing.T) {



}

func TestWriteLinks(ot *testing.T) {


}

func TestSizeLinks(t *testing.T) {

	l := NewLinks()
	l.addKey("hi")
	l.addKey("hello")
	if l.Size() != 2 {
		t.Errorf("sequtil/links: unexpected link set size")
	}

}

func TestSubsetLinksByPrefix(t *testing.T) {


}

func TestDecode(t *testing.T) {

	l := NewLinks()
	l.addKey("a")
	l.addKey("b")
	keys, e := l.Decode([]int{0,1})
	if e != nil {
		t.Error(e)
	}
	if !(keys[0] == "a" && keys[1] == "b") {
		t.Errorf("links.Decode() returned an array of string IDs different from waht was expected.")
	}

}

func TestSubset(t *testing.T) {


}


/*
func TestTabulateVariantLinks(t *testing.T) {

	vp1 := map[int]string{20766468:"R", 20766470:"A"}
	vp2 := map[int]string{95204106:"R", 95204107:"A"}

	l := NewLinks()
	l.TabulateVariantLinks("1", vp1, vp2)
	l.TabulateVariantLinks("1", vp1, vp1)

	check := map[string]map[string]float64 {
		"1_20766468": map[string]float64 {
			"1_20766468": 0.0,
			"1_20766470": -1.0,
			"1_95204106": 1.0,
			"1_95204107": -1.0,
		},
		"1_20766470": map[string]float64 {
			"1_20766468": -1.0,
			"1_20766470": 0.0,
			"1_95204106": -1.0,
			"1_95204107": 1.0,
		}, 
		"1_95204106": map[string]float64 {
			"1_20766468": 1.0,
			"1_20766470": -1.0,
			"1_95204106": 0.0,
			"1_95204107": 0.0,
		}, 
		"1_95204107": map[string]float64 {
			"1_20766468": -1.0,
			"1_20766470": 1.0,
			"1_95204106": 0.0,
			"1_95204107": 0.0,
		},
	}

	fmt.Println(l.data)

	for k1, v1 := range check {
		for k2, v2 :=  range v1 {
			link, _ := l.Get(l.ID(k1), l.ID(k2))
			if link != v2 {
				fmt.Println(k1, k2, link, v2)
				t.Errorf("test error: incorrect inference of variant phase")
			}
		}
	}

}
*/







