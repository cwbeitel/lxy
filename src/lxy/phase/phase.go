package phase

import (
	"bufio"
	"fmt"
	ga "github.com/thoj/go-galib"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	util "sequtil"

	"github.com/golang/glog"
)

var scores int

/*

gb build all; lxy phase infer --links data/GM12878/hic/split/22.links --outpath data/GM12878/hic/split/22 --key=data/GM12878/hic/split/22.key

*/

type PhaseIntermediate struct {
	iteration int
	score     float64
	solution  interface{}
}

// Phase infers a haplotype phasing from a set of variant links
func Phase(links *util.Links, outPath string, iterations int, optReportFreq float64) (map[string]bool, []PhaseIntermediate) {

	if (*links).Size() <= 0 {
		glog.Fatal("error: link object empty, cannot phase without link data")
	}

	rand.Seed(time.Now().UTC().UnixNano())

	m := ga.NewMultiMutator()
	inv := new(GAInvertMutator)
	msw := new(ga.GASwitchMutator)
	m.Add(inv)
	m.Add(msw)

	param := ga.GAParameter{
		Initializer: new(ga.GARandomInitializer),
		Selector:    ga.NewGATournamentSelector(0.7, 5),
		Breeder:     new(ga.GA2PointBreeder),
		Mutator:     m,
		PMutate:     0.7,
		PBreed:      0.7}

	gao := ga.NewGAParallel(param, 7)

	fmt.Println((*links).Size())
	genome := NewFixedBitstringGenome(make([]bool, (*links).Size()), score)

	(*genome).data = links

	gao.Init(10, genome)

	reportEvery := int(math.Floor(float64(iterations)*optReportFreq) + 1.0)
	intermedArraySize := int(math.Floor(float64(iterations) / float64(reportEvery)))
	intermediateSolutions := make([]PhaseIntermediate, intermedArraySize)

	numiter := iterations
	ct := 0
	intermed := 0
	intermedCt := 0

	for {

		ct += 1

		gao.Optimize(1)
		best := gao.Best().(*GAFixedBitstringGenome)
		//fmt.Println("best:", best.Score())

		if ct >= numiter {
			break
		}

		if intermed > reportEvery {
			//fmt.Println(best.Gene)
			fmt.Println("best:", best.Score())
			fmt.Printf("Doing iteration %d (of %d)\n", ct, numiter)
			intermed = 0
			intermedCt += 1
		}
		intermed += 1

	}

	glog.Infof("Finished optimization")
	fmt.Printf("Calls to score = %d\n", scores)
	best := gao.Best().(*GAFixedBitstringGenome)
	//fmt.Println(best)

	dec, _ := (*links).DecodePhasing(best.Gene)
	fmt.Println(dec)

	if e := writePhasing(dec, outPath); e != nil {
		glog.Errorf("error: couldn't write phasing")
	}

	return dec, intermediateSolutions

}

// readPhasing reads a phasing solution from disk and returns an array of booleans corresponding
// to the phase of variants.
//
// readPhasing will return an error if a filesystem path is provided which does not exist on
// the system. Furthermore, since readPhasing is currently designed only for diploids, an
// error will be returned if a value is found in the phasing other than 0 or 1.
func readPhasing(path string) (map[string]bool, error) {

	pf, err := os.Open(path)
	if err != nil {
		fmt.Errorf("Couldn't open input file with path %s\n", path)
	}
	defer pf.Close()

	s := bufio.NewScanner(pf)
	phasing := map[string]bool{}
	//ct := 0
	for s.Scan() {
		arr := strings.Split(s.Text(), " ")
		//fmt.Println(arr)
		if len(arr) != 2 || ((arr[1] != "0") && (arr[1] != "1")) {
			return map[string]bool{}, fmt.Errorf("Tried to read phasing file with more than two haplotypes, not currently supported.\n")
		}
		phasing[arr[0]] = (arr[1] == "1")
		//ct += 1
	}

	/*
		phasing := make(map[string]bool, ct)
		for k, v := range phasingMap {
			phasing[k] = (v == "1")
		}
	*/

	return phasing, nil

}

func writePhasing(phasing map[string]bool, path string) error {
	out, err := os.Create(path)
	if err != nil {
		fmt.Printf("Couldn't open output file (%s) for writing: %s\n", path, err)
	}
	defer out.Close()

	for k, v := range phasing {
		if v {
			out.WriteString(k + " " + strconv.Itoa(1) + "\n")
		} else {
			out.WriteString(k + " " + strconv.Itoa(0) + "\n")
		}
	}

	return err
}

// score determine the quality score of a haplotype phasing as
// represented in a GAFixedBitstringGenome object.
func score(g *GAFixedBitstringGenome) float64 {

	scores++
	total := 0.0
	for i, c := range g.Gene {

		steps := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 15, 20, 25, 30}
		//steps := []int{1, 2, 3}

		for _, v := range steps {

			if (i + v) < (*g.data).Size() {
				c2 := g.Gene[i+v]
				if c != c2 {
					val, _ := (*g.data).Get(i, (i + v))
					total -= val
				} else {
					val, _ := (*g.data).Get(i, (i + v))
					total += val
				}
			}
			if (i - v) > 0 {
				c2 := g.Gene[i-v]
				if c != c2 {
					val, _ := (*g.data).Get(i, (i - v))
					total -= val
				} else {
					val, _ := (*g.data).Get(i, (i - v))
					total += val
				}
			}
		}

	}

	return float64(-total)
}

// EvalPhasing evaluates the quality of a phasing solution relative to a known
// correct phasing.
func EvalPhasing(phasing, key []bool) (float64, float64, float64, float64, error) {

	matches := 0.0
	comparisons := 0.0
	neighbor_matches := 0.0
	neighbor_comparisons := 0.0

	second_neighbor_matches := 0.0
	second_neighbor_comparisons := 0.0
	third_neighbor_matches := 0.0
	third_neighbor_comparisons := 0.0

	for i, _ := range phasing {
		for j, _ := range phasing {
			comparisons += 1
			if (i-j) == 1 || (j-i) == 1 {
				neighbor_comparisons += 1
			}

			if (i-j) == 2 || (j-i) == 2 {
				neighbor_comparisons += 1
			}
			if (i-j) == 3 || (j-i) == 3 {
				neighbor_comparisons += 1
			}

			pmatch := (phasing[i] == phasing[j])
			kmatch := (key[i] == key[j])
			if (pmatch && kmatch) || (!pmatch && !kmatch) {
				matches += 1
				if (i-j) == 1 || (j-i) == 1 {
					neighbor_matches += 1
				}

				if (i-j) == 2 || (j-i) == 2 {
					second_neighbor_matches += 1
				}
				if (i-j) == 3 || (j-i) == 3 {
					third_neighbor_matches += 1
				}

			}
		}
	}

	score := float64(matches / comparisons)
	firstScore := float64(neighbor_matches / neighbor_comparisons)
	secondScore := float64((neighbor_matches + second_neighbor_matches) / (neighbor_comparisons + second_neighbor_comparisons))
	thirdScore := float64((neighbor_matches + second_neighbor_matches + third_neighbor_matches) / (neighbor_comparisons + second_neighbor_comparisons + third_neighbor_comparisons))

	return score, firstScore, secondScore, thirdScore, nil

}

// Some of the hackiest code.....
func EvalPhasingDev(phasing, key map[string]bool) (float64, float64, float64, error) {

	matches := 0.0
	comparisons := 0.0
	matches300k := 0.0
	matches1mbp := 0.0
	comparisons300k := 0.0
	comparisons1mbp := 0.0

	for i, _ := range phasing {
		for j, _ := range phasing {

			indi, _ := strconv.Atoi(strings.Split(i, "_")[1])
			indj, _ := strconv.Atoi(strings.Split(j, "_")[1])
			inddiff := indi - indj
			if inddiff < 0 {
				inddiff = inddiff * -1
			}
			is3k := (inddiff <= 2)
			is1mbp := (inddiff <= 9)

			_, ok1 := key[i]
			_, ok2 := key[j]
			if !ok1 || !ok2 {
				continue
			}

			_, ok3 := phasing[i]
			_, ok4 := phasing[j]
			if !ok3 || !ok4 {
				continue
			}

			pmatch := (phasing[i] == phasing[j])
			kmatch := (key[i] == key[j])
			comparisons += 1
			if is3k {
				comparisons300k += 1
			}
			if is1mbp {
				comparisons1mbp += 1
			}

			if (pmatch && kmatch) || (!pmatch && !kmatch) {
				matches += 1
				if is3k {
					matches300k += 1
				}
				if is1mbp {
					matches1mbp += 1
				}
			}
		}
	}

	//score := float64(matches / comparisons)

	return float64(matches / comparisons), float64(matches300k / comparisons300k), float64(matches1mbp / comparisons1mbp), nil

}
