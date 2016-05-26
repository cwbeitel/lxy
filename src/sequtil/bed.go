package util

import (
	"os"
	"fmt"
	//"strings"
)

// make a bed file to be visualized in the UCSC genome browser showing the alignment 
// position of individual reads.
func GenerateAlignmentTrack() {
	/*
	browser position chr1:1-10000
	track name=alignedreads description="Aligned Hi-C reads" visibility=2
	chr1   5000 5200
	chr1   5050 5250	
	chr1   5100 5300*/
}

// Make a bed file histogram of the read counts per window position to viz in UCSC genome browser
func CoverageTrackFromAlignment(alignmentPath, outPath, name, description string) error {
	// Initialize the bedgraph object
	bg := bedGraph{}
	// reading through the alignment, tabulate the number of reads aligning in each position
	bg.Write(outPath, name, description)
	return nil
}

type bedGraph struct {
	wsize int
	start int
	data map[string][]int
}

func (bg *bedGraph) Write(path, name, description string) {

	/*
	chr1  5000  5050  400
	chr1  5050  5100  405
	chr1  5100  5150  410
	chr1  5150  5200  415
	chr1  5200  5250  420
	chr1  5250  5300  415
	chr1  5300  5350  410
	chr1  5350  5400  405
	chr1  5400  5450  400
	*/

    // Open the ouput file for writing
    out, err := os.Create(path)
    if err != nil {
        fmt.Errorf("Couldn't open output bedgraph file with path %s\n", path)
    }
    defer out.Close()

    // Write header information
    out.WriteString("browser position chr1:1-10000")
    out.WriteString("browser hide all")
    out.WriteString("track type=bedGraph name=" + name + " description=" + description)

    for c, v := range bg.data {
    	for i, d := range v {
    		start := (i-1)*bg.wsize + 1
    		end := i*bg.wsize
    		out.WriteString(fmt.Sprintf("%s %d %d %d", c, start, end, d))
    	}
    }
}


