/*
Copyright 2009 Thomas Jager <mail@jager.no> All rights reserved.
Use of this source code is governed by a BSD-style
license that can be found in the LICENSE file.
Ordered list genome for problems where the order of Genes matter, tSP for example.

Modified for use in LXY where indicated
*/

package scaff

import (
	"fmt"
	"math/rand"
	"sort"
	ga "github.com/thoj/go-galib"
	util "sequtil"
)

type GAOrderedIntGenome struct {
	Gene     []int
	score    float64
	hasscore bool
	sfunc    func(ga *GAOrderedIntGenome) float64
	data *util.Links	// lxy modification
}

func NewOrderedIntGenome(i []int, sfunc func(ga *GAOrderedIntGenome) float64) *GAOrderedIntGenome {
	g := new(GAOrderedIntGenome)
	g.Gene = i
	g.sfunc = sfunc
	g.data = &util.Links{} // lxy modification
	return g
}

//Helper for Partially mapped crossover
func (a *GAOrderedIntGenome) pmxmap(v, p1, p2 int) (int, bool) {
	for i, c := range a.Gene {
		if c == v && (i < p1 || i > p2) {
			return i, true
		}
	}
	return 0, false
}

// Partially mapped crossover.
func (a *GAOrderedIntGenome) Crossover(bi ga.GAGenome, p1, p2 int) (ga.GAGenome, ga.GAGenome) {
	ca := a.Copy().(*GAOrderedIntGenome)
	b := bi.(*GAOrderedIntGenome)
	cb := b.Copy().(*GAOrderedIntGenome)
	copy(ca.Gene[p1:p2+1], b.Gene[p1:p2+1])
	copy(cb.Gene[p1:p2+1], a.Gene[p1:p2+1])
	//Proto child needs fixing
	//amap := new(vector.IntVector)
	//bmap := new(vector.IntVector)
	amap := make([]int, 0)
	bmap := make([]int, 0)
	for i := p1; i <= p2; i++ {
		ma, found := ca.pmxmap(ca.Gene[i], p1, p2)
		if found {
			//amap.Push(ma)
			amap = append(amap, ma)
			//if bmap.Len() > 0 {
			if len(bmap) > 0 {
				//i1 := amap.Pop()
				//i2 := bmap.Pop()
				var i1, i2 int
				i1, amap = amap[len(amap)-1], amap[:len(amap)-1]
				i2, bmap = bmap[len(bmap)-1], bmap[:len(bmap)-1]
				ca.Gene[i1], cb.Gene[i2] = cb.Gene[i2], ca.Gene[i1]
			}
		}
		mb, found := cb.pmxmap(cb.Gene[i], p1, p2)
		if found {
			//bmap.Push(mb)
			bmap = append(bmap, mb)
			//if amap.Len() > 0 {
			if len(amap) > 0 {
				//i1 := amap.Pop()
				//i2 := bmap.Pop()
				var i1, i2 int
				i1, amap = amap[len(amap)-1], amap[:len(amap)-1]
				i2, bmap = bmap[len(bmap)-1], bmap[:len(bmap)-1]
				ca.Gene[i1], cb.Gene[i2] = cb.Gene[i2], ca.Gene[i1]
			}
		}
	}
	ca.Reset()
	cb.Reset()
	return ca, cb
}

func (a *GAOrderedIntGenome) Splice(bi ga.GAGenome, from, to, length int) {
	b := bi.(*GAOrderedIntGenome)
	copy(a.Gene[to:length+to], b.Gene[from:length+from])
	a.Reset()
}

/*func (a *GAOrderedIntGenome) Invert(bi ga.GAGenome, from, to int) {

	b := bi.(*GAOrderedIntGenome)
	for i := from; i < to; i++ {
		fmt.Println(i, (to - i))
		a.Gene[i] = b.Gene[to - i]
	}

	a.Reset()

}*/

func (g *GAOrderedIntGenome) Valid() bool {
	t := g.Copy().(*GAOrderedIntGenome)
	sort.Ints(t.Gene)
	last := -9
	for _, c := range t.Gene {
		if last > -1 && c == last {
			fmt.Printf("%d - %d", c, last)
			return false
		}
		last = c
	}
	return true
}

func (g *GAOrderedIntGenome) Switch(x, y int) {
	g.Gene[x], g.Gene[y] = g.Gene[y], g.Gene[x]
	g.Reset()
}

func (g *GAOrderedIntGenome) Randomize() {
	l := len(g.Gene)
	for i := 0; i < l; i++ {
		x := rand.Intn(l)
		y := rand.Intn(l)
		g.Gene[x], g.Gene[y] = g.Gene[y], g.Gene[x]
	}
	g.Reset()
}

func (g *GAOrderedIntGenome) Copy() ga.GAGenome {
	n := new(GAOrderedIntGenome)
	n.Gene = make([]int, len(g.Gene))
	copy(n.Gene, g.Gene)
	n.sfunc = g.sfunc
	n.score = g.score
	n.hasscore = g.hasscore
	n.data = g.data // lxy modification
	return n
}

func (g *GAOrderedIntGenome) Len() int { return len(g.Gene) }

func (g *GAOrderedIntGenome) Score() float64 {
	if !g.hasscore {
		g.score = g.sfunc(g)
		g.hasscore = true
	}
	return g.score
}

func (g *GAOrderedIntGenome) Invert(p1, p2 int) {

	if p1 > p2 {
		p1, p2 = p2, p1
	}

	// Until you reach the center
	for {
		if p1 >= p2{
			break
		}
		g.Switch(p1, p2)
		p1 += 1
		p2 -= 1
	}

}

func (g *GAOrderedIntGenome) Reset() { g.hasscore = false }

func (g *GAOrderedIntGenome) String() string { return fmt.Sprintf("%v", g.Gene) }


type GAInvertMutator struct{}

func (m GAInvertMutator) Mutate(a ga.GAGenome) ga.GAGenome {
	n := a.Copy()
	p1 := rand.Intn(a.Len())
	p2 := rand.Intn(a.Len())
	if p1 > p2 {
		p1, p2 = p2, p1
	}

	n.Invert(p1, p2)

	/*
	// Until you reach the center
	for {
		if p1 >= p2{
			break
		}
		n.Switch(p1, p2)
		p1 += 1
		p2 -= 1
	}*/

	return n
}
func (m GAInvertMutator) String() string { return "GAInvertMutator" }

