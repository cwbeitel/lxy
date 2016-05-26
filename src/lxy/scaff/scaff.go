package scaff

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/thoj/go-galib"

	util "sequtil"
)

var scores int

/*
e.g. lxy scaff infer --links data/test/GM.1mbp.X.links --output data/test/scaff.real.longrun.out --key data/test/testkey.txt --viz data/test/GM.1mbp.X.png
*/

func WriteScaffolding(scaffolding []string, path string) error {

	out, err := os.Create(path)
	if err != nil {
		fmt.Printf("Couldn't open output file (%s) for writing: %s\n", path, err)
	}
	defer out.Close()

	for _, v := range scaffolding {
		//fmt.Println(v)
		out.WriteString(v + "\n")
	}

	return err

}

func ReadScaffolding(path string) []string {

	in, err := os.Open(path)
	if err != nil {
		fmt.Printf("Couldn't open input file (%s) for reading: %s\n", path, err)
	}
	defer in.Close()

	scaff := []string{}
	s := bufio.NewScanner(in)
	for s.Scan() {
		scaff = append(scaff, s.Text())
	}

	return scaff

}

type Intermediate struct {
	iteration int
	score     float64
	order     []string
}

// Scaffold takes a set of contig links and infers an order for them,
// writing the result to an output path and returning an array of contig
// names in the inferred order.
func Scaffold(links *util.Links, outPath string, iterations int, optReportFreq float64) ([]string, []Intermediate) {

	rand.Seed(time.Now().UTC().UnixNano())

	m := ga.NewMultiMutator()
	msh := new(ga.GAShiftMutator)
	msw := new(ga.GASwitchMutator)
	inv := new(GAInvertMutator)
	m.Add(msh)
	m.Add(msw)
	m.Add(inv)

	param := ga.GAParameter{
		Initializer: new(ga.GARandomInitializer),
		Selector:    ga.NewGATournamentSelector(0.7, 5),
		Breeder:     new(ga.GA2PointBreeder),
		Mutator:     m,
		PMutate:     0.6,
		PBreed:      0.2}

	gao := ga.NewGAParallel(param, 7)
	init := (*links).IntIDs()
	genome := NewOrderedIntGenome(init, score)

	(*genome).data = links

	gao.Init(40, genome)

	reportEvery := int(math.Floor(float64(iterations)*optReportFreq) + 1.0)
	intermedArraySize := int(math.Floor(float64(iterations) / float64(reportEvery)))
	intermediateSolutions := make([]Intermediate, intermedArraySize)

	numiter := iterations
	ct := 0
	intermed := 0
	intermedCt := 0

	for {

		ct += 1
		gao.Optimize(1)
		best := gao.Best().(*GAOrderedIntGenome)

		if ct >= numiter {
			break
		}
		if intermed > reportEvery {
			order, _ := (*links).Decode(best.Gene)
			intermediateSolutions[intermedCt] = Intermediate{ct, best.Score(), order}
			intermedCt += 1
			intermed = 0

			fmt.Println("best:", best.Score())
			fmt.Printf("Doing iteration %d (of %d)\n", ct, numiter)

		}
		intermed += 1

	}

	best := gao.Best().(*GAOrderedIntGenome)
	scaffolding, _ := (*links).Decode(best.Gene)
	err := WriteScaffolding(scaffolding, outPath)
	if err != nil {
		fmt.Printf("Error writing scaffolding: ", err)
	}
	fmt.Println(scaffolding)

	fmt.Printf("Calls to score = %d\n", scores)
	fmt.Printf("%s\n", m.Stats())

	return scaffolding, intermediateSolutions

}

func VisualizeScaffolding(scaffPath, keyPath, outPath string) {

	// DEV
	lxyScriptsDir := "/Users/cb/code/src/github.com/cb01/core/lxy/scripts"

	cmd := exec.Command("python", lxyScriptsDir+"/scaffplot.py", "--inferred", scaffPath, "--actual", keyPath, "--outpath", outPath)
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Debug("Visualizing contig order dotplot...")
	err = cmd.Wait()
	if err != nil {
		log.Fatal("Command finished with error: %v", err)
	}

}

func VisualizeHeatmap(linksPath, outPath, label, orderPath string) {

	lxyScriptsDir := "/Users/cb/code/src/github.com/cb01/core/lxy/scripts"

	cmd := exec.Command("python", lxyScriptsDir+"/heatmap.py", "--inpath", linksPath, "--outpath", outPath, "--label", label, "--ordering", orderPath)
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Debug("Visualizing heatmap...")
	err = cmd.Wait()
	if err != nil {
		log.Fatal("Command finished with error: %v", err)
	}

}

func VisualizeOptimization(scores, label, outpath string) {

	lxyScriptsDir := "/Users/cb/code/src/github.com/cb01/core/lxy/scripts"

	cmd := exec.Command("python", lxyScriptsDir+"/optimizationviz.py", "--scores", scores, "--label", label, "--outpath", outpath)
	fmt.Println(cmd)
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Debug("Visualizing optimization progression...")
	err = cmd.Wait()
	if err != nil {
		log.Fatal("Command finished with error: %v", err)
	}

}

/*
func scoreExp(g *GAOrderedIntGenome) float64 {

	var total float64
	downsampleFac := 0.1
	num := int(math.Floor(float64(len(g.Gene)) * downsampleFac))
	geneIndices := make([]int, num)
	for i := 0; i < num; i++ {
		geneIndices[i] = rand.Intn(len(g.Gene))
	}

	maxDist := math.Floor(float64(len(g.Gene) / 2))

	// Choose five random contigs to evaluate relative to the test
	// set...
	for _, i := range geneIndices {

		c := g.Gene[i]

		for sample := 0; sample < 5; sample++ {

			offset := int(rand.Intn(int(maxDist)))

			//fmt.Println(offset)

			if (i + offset) < (*g.data).Size() {
				val, _ := (*g.data).Get(c, g.Gene[i+offset])
				total += val
			}
			if (i - offset) > 0 {
				val, _ := (*g.data).Get(c, g.Gene[i-offset])
				total += val
			}

		}

	}

	scores++

	return float64(-total)

}
*/

// score determines the quality score for a scaffolding genome solution. Presently,
// the score sums the Hi-C links between a contig and its neighbors out to 1, 2, 3,
// 5, 11, and 20 steps, discounting the more distant steps. This is to enforce the
// expectation that a scaffolding that is in order relative to one that is out of
// order will have more Hi-C links to nearby contigs.
func score(g *GAOrderedIntGenome) float64 {

	var total float64

	for i, c := range g.Gene {

		if (i + 1) < (*g.data).Size() {
			val, _ := (*g.data).Get(c, g.Gene[i+1])
			total += val
		}
		if (i - 1) > 0 {
			val, _ := (*g.data).Get(c, g.Gene[i-1])
			total += val
		}
		if (i + 2) < (*g.data).Size() {
			val, _ := (*g.data).Get(c, g.Gene[i+2])
			total += 0.5 * val
		}
		if (i - 2) > 0 {
			val, _ := (*g.data).Get(c, g.Gene[i-2])
			total += 0.5 * val
		}
		if (i + 3) < (*g.data).Size() {
			val, _ := (*g.data).Get(c, g.Gene[i+3])
			total += 0.33 * val
		}
		if (i - 3) > 0 {
			val, _ := (*g.data).Get(c, g.Gene[i-3])
			total += 0.33 * val
		}

		if (i + 5) < (*g.data).Size() {
			val, _ := (*g.data).Get(c, g.Gene[i+5])
			total += 0.2 * val
		}
		if (i - 5) > 0 {
			val, _ := (*g.data).Get(c, g.Gene[i-5])
			total += 0.2 * val
		}

		if (i + 11) < (*g.data).Size() {
			val, _ := (*g.data).Get(c, g.Gene[i+11])
			total += 0.1 * val
		}
		if (i - 11) > 0 {
			val, _ := (*g.data).Get(c, g.Gene[i-11])
			total += 0.1 * val
		}

		if (i + 20) < (*g.data).Size() {
			val, _ := (*g.data).Get(c, g.Gene[i+20])
			total += 0.05 * val
		}
		if (i - 20) > 0 {
			val, _ := (*g.data).Get(c, g.Gene[i-20])
			total += 0.05 * val
		}

	}

	scores++

	return float64(-total)

}

func distScore(d, c1, c2 int) int {

	actualDist := (c1 - c2)
	if actualDist < 0 {
		actualDist = -1 * actualDist
	}

	score := d - actualDist
	if score < 0 {
		return -1 * score
	}
	return score

}

// EvalScaffolding evaluates the quality of a scaffolding solution relative to a known
// correct scaffolding.
func EvalScaffolding(scaff []string, key []string) (float64, float64, error) {

	// Cache the order of each element in the key
	keyOrder := map[string]int{}
	keyOrderReverse := map[string]int{}
	for i, v := range key {
		keyOrder[v] = i
		keyOrderReverse[v] = len(key) - i
	}

	comparisons := 0.0
	correct := 0.0
	correctReverse := 0.0
	neighbors := 0.0
	neighborsCorrect := 0.0
	neighborsCorrectReverse := 0.0

	// Forward

	// For each triplet in the inferred solution, check whether
	// that triplet is in the same order in the key.
	for i, v1 := range scaff {
		for j, v2 := range scaff {
			for k, v3 := range scaff {
				if (i < j) && (k > j) {
					comparisons += 1
					if (keyOrder[v1] < keyOrder[v2]) && (keyOrder[v2] < keyOrder[v3]) {
						correct += 1
					}
					if (keyOrderReverse[v1] < keyOrderReverse[v2]) && (keyOrderReverse[v2] < keyOrderReverse[v3]) {
						correctReverse += 1
					}
					if (j == i+1) && (k == j+1) {
						neighbors += 1
						if (keyOrder[v1] < keyOrder[v2]) && (keyOrder[v2] < keyOrder[v3]) {
							neighborsCorrect += 1
						}
						if (keyOrderReverse[v1] < keyOrderReverse[v2]) && (keyOrderReverse[v2] < keyOrderReverse[v3]) {
							neighborsCorrectReverse += 1
						}
					}
				}
			}
		}
	}

	qscoreF := float64(correct / comparisons)
	qscoreR := float64(correctReverse / comparisons)

	if qscoreF > qscoreR {
		return qscoreF, float64(neighborsCorrect / neighbors), nil
	} else {
		return qscoreR, float64(neighborsCorrectReverse / neighbors), nil
	}

}
