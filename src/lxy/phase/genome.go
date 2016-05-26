/*
Copyright 2009 Thomas Jager <mail@jager.no> All rights reserved.
Use of this source code is governed by a BSD-style
license that can be found in the LICENSE file.
Fixed length Bitstring genome for problems like subset sum

Modified for use in LXY where indicated
*/

package phase

import (
	"fmt"
	"math/rand"
	ga "github.com/thoj/go-galib"
	util "sequtil"
)

type GAFixedBitstringGenome struct {
	Gene     []bool
	score    float64
	hasscore bool
	sfunc    func(ga *GAFixedBitstringGenome) float64
	data *util.Links	// lxy modification
}

func NewFixedBitstringGenome(i []bool, sfunc func(ga *GAFixedBitstringGenome) float64) *GAFixedBitstringGenome {
	g := new(GAFixedBitstringGenome)
	g.Gene = i
	g.sfunc = sfunc
	g.Reset()
	g.data = &util.Links{} // lxy modification
	return g
}

//Simple 2 point crossover
func (a *GAFixedBitstringGenome) Crossover(bi ga.GAGenome, p1, p2 int) (ga.GAGenome, ga.GAGenome) {
	ca := a.Copy().(*GAFixedBitstringGenome)
	b := bi.(*GAFixedBitstringGenome)
	cb := b.Copy().(*GAFixedBitstringGenome)
	copy(ca.Gene[p1:p2+1], b.Gene[p1:p2+1])
	copy(cb.Gene[p1:p2+1], a.Gene[p1:p2+1])
	ca.Reset()
	cb.Reset()
	return ca, cb
}

func (a *GAFixedBitstringGenome) Splice(bi ga.GAGenome, from, to, length int) {
	b := bi.(*GAFixedBitstringGenome)
	copy(a.Gene[to:length+to], b.Gene[from:length+from])
	a.Reset()
}

func (g *GAFixedBitstringGenome) Valid() bool { return true }

func (g *GAFixedBitstringGenome) Switch(x, y int) {
	g.Gene[x], g.Gene[y] = g.Gene[y], g.Gene[x]
	g.Reset()
}

func (g *GAFixedBitstringGenome) Randomize() {
	l := len(g.Gene)
	for i := 0; i < l; i++ {
		x := rand.Intn(2)
		if x == 1 {
			g.Gene[i] = true
		} else {
			g.Gene[i] = false
		}
	}
	g.Reset()
}

func (g *GAFixedBitstringGenome) Copy() ga.GAGenome {
	n := new(GAFixedBitstringGenome)
	n.Gene = make([]bool, len(g.Gene))
	copy(n.Gene, g.Gene)
	n.sfunc = g.sfunc
	n.score = g.score
	n.hasscore = g.hasscore
	n.data = g.data // lxy modification
	return n
}

func (g *GAFixedBitstringGenome) Len() int { return len(g.Gene) }

func (g *GAFixedBitstringGenome) Score() float64 {
	if !g.hasscore {
		g.score = g.sfunc(g)
		g.hasscore = true
	}
	return g.score
}

func (g *GAFixedBitstringGenome) Invert(p1, p2 int) {

	if p1 > p2 {
		p1, p2 = p2, p1
	}

	for i := p1; i < p2; i++ {
		if g.Gene[i] {
			g.Gene[i] = false
		} else {
			g.Gene[i] = true
		}
	}

}

func (g *GAFixedBitstringGenome) Reset() { g.hasscore = false }

func (g *GAFixedBitstringGenome) String() string {
	return fmt.Sprintf("%v", g.Gene)
}

type GAInvertMutator struct{}

func (m GAInvertMutator) Mutate(a ga.GAGenome) ga.GAGenome {
	n := a.Copy()
	p1 := rand.Intn(a.Len())
	p2 := rand.Intn(a.Len())
	if p1 > p2 {
		p1, p2 = p2, p1
	}

	n.Invert(p1, p2)

	return n
}
func (m GAInvertMutator) String() string { return "GAInvertMutator" }





