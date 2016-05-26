package scaff

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"

	util "sequtil"
)

/*
lxy scaff infer --links data/test/GM.1mbp.X.links --output data/test/scaff.real.longrun.out --key data/test/testkey.txt --viz data/test/GM.1mbp.X.png

*/

func ScaffoldCommand() cli.Command {
	return cli.Command{
		Name:  "scaff",
		Usage: "Order and orient contigs using Hi-C contact frequency data.",
		Subcommands: []cli.Command{
			cli.Command{
				Name:  "infer",
				Usage: "Infer a scaffolding from link data, e.g. lxy scaff infer --links data/test/GM.1mbp.links --output data/test/testscaffolding.txt --subset X --key data/test/testkey.txt --viz data/test/testorderfig.png",
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "debug",
						Usage: "Whether to print detailed debugging information.",
					},
					cli.BoolFlag{
						Name:  "heatmapCompare",
						Usage: "Whether to genereate heatmaps visualizing the consistency of the inferred order with the contact frequency data.",
					},
					cli.StringFlag{
						Name:  "links",
						Value: "",
						Usage: "Path to the Hi-C links file.",
					},
					cli.StringFlag{
						Name:  "subset",
						Value: "",
						Usage: "Tag on basis of which to subset links.",
					},
					cli.StringFlag{
						Name:  "outputPrefix",
						Value: "",
						Usage: "File path stem for output files.",
					},
					cli.StringFlag{
						Name:  "key",
						Value: "",
						Usage: "Path to scaffolding key file.",
					},
					cli.StringFlag{
						Name:  "viz",
						Value: "",
						Usage: "Path to output the ordering viz, requires --key.",
					},
					cli.IntFlag{
						Name:  "iterations",
						Value: 1000,
						Usage: "Number of iterations or GA 'generations' to perform.",
					},
					cli.Float64Flag{
						Name:  "optReportFreq",
						Value: 0.01,
						Usage: "Frequency of storing intermediate ordererings for later evaluation of progression of optimization.",
					},
					cli.IntFlag{
						Name:  "popSize",
						Value: 50,
						Usage: "Population size for genetic algorithm.",
					},
					cli.Float64Flag{
						Name:  "breedProb",
						Value: 0.8,
						Usage: "Probability of breeding between two members of population.",
					},
					cli.Float64Flag{
						Name:  "mutProb",
						Value: 0.8,
						Usage: "Probability of a mutation occurring.",
					},
				},
				Action: scaffoldInferCommand,
			},
			cli.Command{
				Name:  "eval",
				Usage: "Evaluate a scaffolding using a key ordering.",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "scaffolding",
						Value: "",
						Usage: "Path to scaffolding file.",
					},
					cli.StringFlag{
						Name:  "key",
						Value: "",
						Usage: "Path to scaffolding key file.",
					},
					cli.StringFlag{
						Name:  "output",
						Value: "",
						Usage: "Output path for evaluation stats file.",
					},
				},
				Action: evalScaffoldingCommand,
			},
			cli.Command{
				Name:  "prep",
				Usage: "Generate a links file from a set of aligned Hi-C reads.",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "sam",
						Value: "",
						Usage: "Path to sam file.",
					},
					cli.StringFlag{
						Name:  "output",
						Value: "",
						Usage: "Output path for links file.",
					},
				},
				Action: prepScaffoldingCommand,
			},
		},
	}
}

/*



 */

func scaffoldInferCommand(c *cli.Context) {

	if len(c.String("links")) == 0 {
		fmt.Printf("error: must provide a path to a links file with --links\n")
		return
	}
	links, _ := util.LoadLinks(c.String("links"))

	if _, err := os.Stat(c.String("links")); os.IsNotExist(err) {
		fmt.Printf("error: the specified links file does not exist: %s\n", c.String("links"))
		return
	}

	if len(c.String("outputPrefix")) == 0 {
		fmt.Printf("error: must provide an output prefix with --outputPrefix\n")
		return
	}
	scaffOutput := c.String("outputPrefix") + ".scaff.txt"

	if len(c.String("key")) == 0 {
		fmt.Printf("error (development): must provide a key contig ordering with --key\n")
		return
	}
	key := ReadScaffolding(c.String("key"))

	// Perform the scaffolding
	scaffolding, intermediateSolutions := Scaffold(&links, scaffOutput, c.Int("iterations"), c.Float64("optReportFreq"))
	// intermediate is a vector of iteration;metricquality;order

	// Evaluate the quality of the scaffolding
	score, nscore, _ := EvalScaffolding(scaffolding, key)
	fmt.Printf("Evaluated scaffolding with score %f and neighbor score %f\n", score, nscore)
	qscorepath := c.String("outputPrefix") + ".qscore.txt"
	out, err := os.Create(qscorepath)
	if err != nil {
		fmt.Printf("Couldn't open output file (%s) for writing: %s\n", qscorepath, err)
	}
	defer out.Close()
	out.WriteString(fmt.Sprintf("score=%f;score_neighbor=%f\n", score, nscore))
	// -----------

	// Visualize scaffolding quality
	orderVizOutput := c.String("outputPrefix") + ".orderviz.png"
	VisualizeScaffolding(scaffOutput, c.String("key"), orderVizOutput)
	if len(c.String("heatmapCompare")) != 0 {
		trueHeatmapOutpath := c.String("outputPrefix") + ".heat.true.png"
		inferredHeatmapOutpath := c.String("outputPrefix") + ".heat.inferred.png"
		// Visualize the true order heatmap
		VisualizeHeatmap(c.String("links"), trueHeatmapOutpath, "true", c.String("key"))
		// Visualize the inferred order heatmap
		VisualizeHeatmap(c.String("links"), inferredHeatmapOutpath, "inferred", scaffOutput)
	}

	// Visualize the improvement in the quality score of intermediate scaffoldings through the course of the optimization
	optDataPath := c.String("outputPrefix") + ".opt.txt"
	optOut, err := os.Create(optDataPath)
	if err != nil {
		fmt.Printf("Couldn't open output file (%s) for writing: %s\n", optDataPath, err)
	}
	defer optOut.Close()
	for _, intermediateSolution := range intermediateSolutions {
		score, _, _ := EvalScaffolding(intermediateSolution.order, key)
		optOut.WriteString(fmt.Sprintf("%d %f %f\n", intermediateSolution.iteration, intermediateSolution.score, score))
	}

	VisualizeOptimization(optDataPath, "test", c.String("outputPrefix")+".opt.png")

}

func evalScaffoldingCommand(c *cli.Context) {

	if len(c.String("scaffolding")) == 0 {
		fmt.Println("error: need a scaffolding solution to evaluate")
		return
	}

	if len(c.String("key")) == 0 {
		fmt.Println("error: need an ordering key to evaluate scaffolding solution")
		return
	}

	scaff := ReadScaffolding(c.String("scaffolding"))
	key := ReadScaffolding(c.String("key"))
	score, neighborScore, _ := EvalScaffolding(scaff, key)
	fmt.Printf("%f\n%f\n", score, neighborScore)

}

func prepScaffoldingCommand(c *cli.Context) {

	util.ScaffoldLinksFromSam(c.String("sam"), c.String("output"))

}
