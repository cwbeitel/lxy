package phase

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	util "sequtil"

	"github.com/golang/glog"
)

/*

gb build all; lxy phase prep varlinks --vcf data/GM12878/vars/bychr/blocks/vars.22.blk.vcf --sam data/GM12878/hic/split/22.sam --output data/GM12878/hic/split/22.links

*/

func PhaseCommand() cli.Command {
	return cli.Command{
		Name:  "phase",
		Usage: "Commands to infer and evaluate haplotype phasing solutions",
		Subcommands: []cli.Command{
			cli.Command{
				Name:  "infer",
				Usage: "Given a bipartite variant graph, determine the phase of the two partitions.",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "links",
						Value: "",
						Usage: "Path to the Hi-C links file.",
					},
					cli.StringFlag{
						Name:  "outpath",
						Value: "",
						Usage: "Output path for inferred phasing.",
					},
					cli.StringFlag{
						Name:  "key",
						Value: "",
						Usage: "Path to file specifying correct phasing of chromosome.",
					},
					cli.IntFlag{
						Name:  "iterations",
						Value: 1000,
						Usage: "Number of iterations or GA 'generations' to perform.",
					},
					cli.IntFlag{
						Name:  "subruntag",
						Value: 0,
						Usage: "The sub run tag to assign to this run.",
					},
					cli.Float64Flag{
						Name:  "optReportFreq",
						Value: 0.01,
						Usage: "Frequency of storing intermediate ordererings for later evaluation of progression of optimization.",
					},
				},
				Action: phaseInferCommand,
			},
			cli.Command{
				Name:  "eval",
				Usage: "Evaluate a haplotype phasing solution.",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "phasing",
						Value: "",
						Usage: "Path to phasing file.",
					},
					cli.StringFlag{
						Name:  "key",
						Value: "",
						Usage: "Path to phasing key file.",
					},
					cli.StringFlag{
						Name:  "output",
						Value: "",
						Usage: "Output path for evaluation stats file.",
					},
				},
				Action: evalPhasingCommand,
			},
			cli.Command{
				Name:  "prep",
				Usage: "Generate a variant links object from a set of aligned Hi-C reads.",
				Subcommands: []cli.Command{
					cli.Command{
						Name:  "varlinks",
						Usage: "Generate a variant links object from a set of aligned Hi-C reads.",
						Flags: []cli.Flag{
							cli.StringFlag{
								Name:  "sam",
								Value: "",
								Usage: "Path to sam file.",
							},
							cli.StringFlag{
								Name:  "vcf",
								Value: "",
								Usage: "Path to vcf file giving variants.",
							},
							cli.StringFlag{
								Name:  "output",
								Value: "",
								Usage: "Output path for links file.",
							},
						},
						Action: prepPhasingCommand,
					},
				},
			},
		},
	}

}

func phaseInferCommand(c *cli.Context) {

	if len(c.String("links")) == 0 {
		glog.Errorf("error: must provide a path to a links file with --links.")
		return
	}

	if _, err := os.Stat(c.String("links")); os.IsNotExist(err) {
		glog.Errorf("error: the specified links file does not exist: %s\n", c.String("links"))
		return
	}

	if len(c.String("outpath")) == 0 {
		glog.Errorf("error: must provide a path to an output file stem with --outpath.")
		return
	}

	if len(c.String("key")) == 0 {
		glog.Errorf("error: must provide a path to a key file with --key\n")
		return
	}

	links, e2 := util.LoadLinks(c.String("links"))
	if e2 != nil {
		glog.Errorf("Error loading links: %s\n", e2)
	}

	phasingOutPath := c.String("outpath") + "." + c.String("subruntag") + ".phasing.txt"
	phasing, _ := Phase(&links, phasingOutPath, c.Int("iterations"), c.Float64("optReportFreq"))

	key, ek := readPhasing(c.String("key"))
	if ek != nil {
		glog.Errorf("error: %s\n", ek)
		return
	}

	//score, neighborScore, n2score, n3score, e := EvalPhasingDev(phasing, key)
	score, score3k, score1mbp, e := EvalPhasingDev(phasing, key)
	if e != nil {
		glog.Errorf("error evaluating phasing: %s\n", e)
	}
	fmt.Printf("Evaluated phasing with scores %f, %f, %f\n", score, score3k, score1mbp)
	qscorepath := c.String("outpath") + "." + c.String("subruntag") + ".qscore.txt"
	out, err := os.Create(qscorepath)
	if err != nil {
		fmt.Printf("Couldn't open output file (%s) for writing: %s\n", qscorepath, err)
	}
	defer out.Close()
	//out.WriteString(fmt.Sprintf("score=%f;nscore=%f;n2score=%f;n3score=%f\n", score, neighborScore, n2score, n3score))
	out.WriteString(fmt.Sprintf("score=%f;300kbp_score=%f;1mbp_score=%f\n", score, score3k, score1mbp))

}

func evalPhasingCommand(c *cli.Context) {

	if len(c.String("phasing")) == 0 {
		glog.Errorf("error: must provide a path to a phasing file with --phasing\n")
		return
	}

	if _, err := os.Stat(c.String("phasing")); os.IsNotExist(err) {
		glog.Errorf("error: the specified phasing file does not exist: %s\n", c.String("phasing"))
		return
	}

	if len(c.String("key")) == 0 {
		glog.Errorf("error: must provide a path to a key file with --key\n")
		return
	}

	if _, err := os.Stat(c.String("key")); os.IsNotExist(err) {
		glog.Errorf("error: the specified key file does not exist: %s\n", c.String("key"))
		return
	}

	phasing, ep := readPhasing(c.String("phasing"))
	if ep != nil {
		glog.Errorf("error: %s", ep)
		return
	}

	key, ek := readPhasing(c.String("key"))
	if ek != nil {
		glog.Errorf("error: %s\n", ek)
		return
	}

	//score, nscore, _, _, e := EvalPhasingDev(phasing, key)
	score, score3k, score1mbp, e := EvalPhasingDev(phasing, key)
	if e != nil {
		glog.Errorf("error evaluating phasing: %s\n", e)
	}

	//fmt.Printf("Evaluated phasing with score %f and neighbor score %f\n", score, nscore)
	fmt.Printf("Evaluated phasing with scores %f, %f, %f\n", score, score3k, score1mbp)

	if len(c.String("output")) != 0 {

	}

}

func prepPhasingCommand(c *cli.Context) {

	if len(c.String("vcf")) == 0 {
		glog.Errorf("error: must provide a path to a vcf file with --vcf\n")
		return
	}

	if _, err := os.Stat(c.String("vcf")); os.IsNotExist(err) {
		glog.Errorf("error: the specified vcf file does not exist: %s\n", c.String("vcf"))
		return
	}

	if len(c.String("sam")) == 0 {
		glog.Errorf("error: must provide a path to a sam file with --sam\n")
		return
	}

	if _, err := os.Stat(c.String("sam")); os.IsNotExist(err) {
		glog.Errorf("error: the specified sam file does not exist: %s\n", c.String("sam"))
		return
	}

	if len(c.String("output")) == 0 {
		glog.Errorf("error: must provide a path to an output destination with --output\n")
		return
	}

	fmt.Println("Reading variants...")
	vcf, num, _ := util.ReadVariants(c.String("vcf"))

	if num == 0 {
		glog.Errorf("error: the length of the variant data vector should be nonzero in order to perform haplotype phasing.\n")
		return
	}

	fmt.Println("Parsing sam to variant links...")
	//util.VariantLinksFromSam(c.String("sam"), c.String("output"), vcf)
	util.BlockLinksFromSam(c.String("sam"), c.String("output"), vcf)

}
