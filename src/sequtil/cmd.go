package util

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"os"
	"os/exec"
)

func VarsCommand() cli.Command {

	return cli.Command{
		Name:  "vars",
		Usage: "A set of utility functions for working with variants.",
		Subcommands: []cli.Command{
			cli.Command{
				Name:   "simblocks",
				Usage:  "Given a known haplotype, simulate haplotype blocks.",
				Action: variantSimBlocksCommand,
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "input",
						Value: "",
						Usage: "Input VCF file with known phasing.",
					},
					cli.IntFlag{
						Name:  "blocksize",
						Value: 10000,
						Usage: "The size in basepairs of variant blocks to simulate.",
					},
					cli.StringFlag{
						Name:  "output",
						Value: "",
						Usage: "The destination path to which to write the output.",
					},
				},
			},
			cli.Command{
				Name:   "aggregate",
				Usage:  "Given variant links and known haplotype blocks, aggregate variant links into variant block links.",
				Action: aggregateVlinksOverBlocksCommand,
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "blocks",
						Value: "",
						Usage: "Input containing variant blocks.",
					},
					cli.StringFlag{
						Name:  "varlinks",
						Value: "",
						Usage: "Input containing variant links.",
					},
					cli.StringFlag{
						Name:  "output",
						Value: "",
						Usage: "The destination path to which to write the output.",
					},
				},
			},
			cli.Command{
				Name:   "stats",
				Usage:  "Variant statistics.",
				Action: variantStatsCommand,
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "input",
						Value: "",
						Usage: "Path to the fasta to be masked.",
					},
					cli.StringFlag{
						Name:  "output",
						Value: "",
						Usage: "Path to the vcf indicating positions to be masked.",
					},
				},
			},
		},
	}

}

func variantSimBlocksCommand(c *cli.Context) {
	// TODO
}

func aggregateVlinksOverBlocksCommand(c *cli.Context) {
	// TODO
}

func variantStatsCommand(c *cli.Context) {
	// TODO
}

func SeqCommand() cli.Command {
	return cli.Command{
		Name:  "seq",
		Usage: "A set of utility functions for supporting work with Hi-C data.",
		Subcommands: []cli.Command{
			cli.Command{
				Name:   "partition",
				Usage:  "Partition a fasta into a set of contigs of a specified size.",
				Action: partitionCommand,
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "fasta",
						Value: "",
						Usage: "Path to the fasta to be partitioned.",
					},
					cli.IntFlag{
						Name:  "windowsize",
						Value: 10000,
						Usage: "The size of the window to use to partition the fasta",
					},
					cli.StringFlag{
						Name:  "output",
						Value: "",
						Usage: "The destination path to which to write the output.",
					},
				},
			},
			cli.Command{
				Name:   "mask",
				Usage:  "Mask a fasta at a set of variant positions.",
				Action: maskCommand,
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "fasta",
						Value: "",
						Usage: "Path to the fasta to be masked.",
					},
					cli.StringFlag{
						Name:  "vcf",
						Value: "",
						Usage: "Path to the vcf indicating positions to be masked.",
					},
					cli.StringFlag{
						Name:  "output",
						Value: "",
						Usage: "The destination path to which to write the output.",
					},
				},
			},
			cli.Command{
				Name:   "align",
				Usage:  "Wrapper for BWA MEM sequence alignment.",
				Action: alignCommand,
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "input",
						Value: "",
						Usage: "FASTQ reads to align to the reference sequence.",
					},
					cli.StringFlag{
						Name:  "output",
						Value: "",
						Usage: "Path to which to write aligned reads.",
					},
					cli.StringFlag{
						Name:  "ref",
						Value: "",
						Usage: "Indexed reference sequence.",
					},
				},
			},
			cli.Command{
				Name:   "index",
				Usage:  "Wrapper for BWA sequence indexing.",
				Action: alignCommand,
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "input",
						Value: "",
						Usage: "FASTQ reads to align to the reference sequence.",
					},
					cli.StringFlag{
						Name:  "output",
						Value: "",
						Usage: "Path to which to write aligned reads.",
					},
					cli.StringFlag{
						Name:  "ref",
						Value: "",
						Usage: "Indexed reference sequence.",
					},
				},
			},
		},
	}

}

func alignCommand(c *cli.Context) {

	// Check that the inputs exist
	if len(c.String("input")) == 0 {
		fmt.Printf("error: must provide a fastq of reads to align to the reference with --input\n")
		return
	}

	if _, err := os.Stat(c.String("input")); os.IsNotExist(err) {
		fmt.Printf("error: the provided input reads file does not exist: %s\n", c.String("input"))
		return
	}

	if len(c.String("ref")) == 0 {
		fmt.Printf("error: must provide a reference to which to align with --ref\n")
		return
	}

	if _, err := os.Stat(c.String("ref")); os.IsNotExist(err) {
		fmt.Printf("error: the specified reference assembly does not exist: %s\n", c.String("ref"))
		return
	}

	// If the reference assembly has not yet been indexed, index it.
	if _, err := os.Stat(c.String("ref") + ".ann"); os.IsNotExist(err) {
		fmt.Printf("error: the reference assembly is not indexed, attempting to do so now...")
		indexCommand(c)
	}

	// Command string for the following
	// bwa mem -t 1 path/to/ref/stem.fa path/to/reads.fq > out.sam
	cmd := exec.Command("bwa", "mem", "-t4", c.String("ref"), c.String("input"))

	outfile, err := os.Create(c.String("output"))
	if err != nil {
		log.Fatal(err)
	}
	defer outfile.Close()
	cmd.Stdout = outfile
	cmd.Stderr = os.Stderr

	err2 := cmd.Start()
	if err2 != nil {
		log.Fatal(err2)
	}

	log.Debug("Performing sequence alignment with BWA...")
	err3 := cmd.Wait()
	if err3 != nil {
		log.Fatal("An error occurred during sequence alignment: %v", err3)
	}

}

func indexCommand(c *cli.Context) {

	// Check that the input exists
	if len(c.String("ref")) == 0 {
		fmt.Printf("error: must provide a reference to which to align with --ref\n")
		return
	}

	if _, err := os.Stat(c.String("ref")); os.IsNotExist(err) {
		fmt.Printf("error: the specified reference assembly does not exist: %s\n", c.String("ref"))
		return
	}

	// Construct the command
	commandString := "bwa index " + c.String("ref")

	// Exec the command and handle various outcomes
	cmd := exec.Command(commandString)
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Debug("Indexing reference assembly...")
	err = cmd.Wait()
	if err != nil {
		log.Fatal("An error occurred during reference indexing: %v", err)
	}

}

func partitionCommand(c *cli.Context) {
	Partition(c.String("fasta"), c.String("output"), c.Int("windowsize"))
}

func maskCommand(c *cli.Context) {
	vcf, _, _ := ReadVariants(c.String("vcf"))
	Mask(c.String("fasta"), c.String("output"), vcf)
}
