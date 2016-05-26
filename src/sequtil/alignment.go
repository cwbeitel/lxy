package util

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	//"github.com/codegangsta/cli"
)

// Alignment is an object which represents a sequence read alignment record with fields
// matching the SAM sequence alignment format.
type Alignment struct {
	qname string  // Query template NAME
	flag  SAMFlag // bitwise FLAG
	rname string  // Reference sequence NAME
	pos   int     // 1-based leftmost mapping POSition
	mapq  int     // MAPping Quality
	cigar CIGAR   // CIGAR string
	rnext string  // Ref. name of the mate/next read
	pnext int     // Position of the mate/next read
	tlen  int     // observed Template LENgth
	seq   string  // segment SEQuence
	qual  string  // ASCII of Phred-scaled base QUALity+33ASCII of Phred-scaled base QUALity+33
}

// SAMFlag represents in decoded (easily accessible) for a SAM bitwise flag.
type SAMFlag struct {
	multiseg        bool // template having multiple segments in sequencing
	allpropper      bool // each segment properly aligned according to the aligner
	unmapped        bool // segment unmapped
	nextunmapped    bool // next segment in the template unmapped
	reversecomp     bool // SEQ being reverse complemented
	nextreversecomp bool // SEQ of the next segment in the template being reverse complemented
	first           bool // the first segment in the template
	last            bool // the last segment in the template
	secondary       bool // secondary alignment
	nopass          bool // not passing filters, such as platform/vendor quality controls
	duplicate       bool // PCR or optical duplicate
	supplementary   bool // supplementary alignment
	code            int  // the integer encoded SAM flag
}

// parseSAMFlag takes an encoded SAM flag in string form and returns a structured SAMFlag
// object representing the information encoded in the flag, see SAMFlag.
func parseSAMFlag(encodedString string) (SAMFlag, error) {

	// Convert the string flag into an integer
	encoded, _ := strconv.Atoi(encodedString)

	return SAMFlag{
		(encoded&0x1 != 0),
		(encoded&0x2 != 0),
		(encoded&0x4 != 0),
		(encoded&0x8 != 0),
		(encoded&0x10 != 0),
		(encoded&0x20 != 0),
		(encoded&0x40 != 0),
		(encoded&0x80 != 0),
		(encoded&0x100 != 0),
		(encoded&0x200 != 0),
		(encoded&0x400 != 0),
		(encoded&0x800 != 0),
		encoded,
	}, nil

}

// CIGAR is an array of CIGARCode's which is used to represent the way a read
// aligns to a particular region.
//
// code description
// M 	alignment match (can be a sequence match or mismatch) insertion to the reference
// I 	insertion to the reference
// D 	deletion from the reference
// N 	skipped region from the reference
// S 	soft clipping (clipped sequences present in SEQ)
// H 	hard clipping (clipped sequences NOT present in SEQ) padding (silent deletion from padded reference)
// P 	padding (silent deletion from padded reference)
// =	sequence match
// X	sequence mismatch
type CIGAR []CIGARCode

// CIGARCode is a single code element of a CIGAR array.
type CIGARCode struct {
	value int
	code  string
}

// Add operates on a CIGAR object to add a new CIGARCode entry to the CIGAR code array.
func (c *CIGAR) Add(val int, code string) {
	cc := CIGARCode{val, code}
	(*c) = append(*c, cc)
}

// String takes a structured CIGAR object and returns a compacted string representation
// of the kind that can be read and written to a SAM file.
func (c *CIGAR) String() string {
	cstring := ""
	// Append each CIGARCode element to the output string in VALUECODE format.
	for _, cc := range *c {
		cstring = fmt.Sprintf("%s%d%s", cstring, cc.value, cc.code)
	}
	return cstring
}

// if it's star, handle that appropriately <- **
func parseCIGAR(cigarString string) (CIGAR, error) {
	// for now, must contain only S, H, or M
	///disallow := map[string]int{"I":1, "D":1, "N":1, "P":1, "=":1, "X":1,}
	cigar := CIGAR{}
	lastIndex := 0
	for i, c := range cigarString {
		if c > 57 {

			val, err := strconv.Atoi(cigarString[lastIndex:i])
			if err != nil {
				return CIGAR{}, fmt.Errorf("Error parsing CIGAR string.")
			}
			cigar.Add(val, string(c))
			lastIndex = i + 1
		}
	}
	return cigar, nil
}

// parseSAMLine takes a SAM formatted line and parses it into an Alignment object.
func parseSAMLine(line string) (Alignment, error) {
	arr := strings.Split(line, "\t")
	if len(arr) < 11 {
		return Alignment{}, fmt.Errorf("Incorrect length alignment")
	}

	qname := arr[0]

	samflag, esf := parseSAMFlag(arr[1])
	if esf != nil {
		return Alignment{}, esf
	}

	rname := arr[2]
	if rname == "*" {
		// didn't align, return null alignment for now
		return Alignment{}, fmt.Errorf("Read did not align")
	}

	pos, ep := strconv.Atoi(arr[3])
	if ep != nil {
		return Alignment{}, ep
	}

	mapq, em := strconv.Atoi(arr[4])
	if em != nil {
		return Alignment{}, em
	}

	c, ec := parseCIGAR(arr[5])
	if ec != nil {
		return Alignment{}, ec
	}

	rnext := arr[6]

	pnext, epn := strconv.Atoi(arr[7])
	if epn != nil {
		if arr[7] == "*" {
			pnext = -1
		} else {
			return Alignment{}, epn
		}
	}

	tlen, etl := strconv.Atoi(arr[8])
	if etl != nil {
		return Alignment{}, etl
	}

	seq := arr[9]
	qual := arr[10]

	a := Alignment{qname, samflag,
		rname, pos, mapq, c, rnext,
		pnext, tlen, seq, qual}

	return a, nil

}

// String generates a SAM format string record from an Alignment object.
func (a Alignment) String() string {
	stringCIGAR := a.cigar.String()
	stringAlignment := fmt.Sprintf("%s\t%d\t%s\t%d\t%d\t%s\t%s\t%d\t%d\t%s\t%s", a.qname, a.flag.code, a.rname, a.pos, a.mapq, stringCIGAR, a.rnext, a.pnext, a.tlen, a.seq, a.qual)
	return stringAlignment
}

func GenomePositions(a Alignment) (map[int]string, error) {

	hits := map[int]string{}
	offset := 0
	firstH := 1
	for _, c := range a.cigar {
		if c.code == "M" {
			for i := 0; i < c.value; i++ {
				ind := (a.pos + offset)
				hits[ind] = string(a.seq[offset])
				offset += 1
			}
		} else if c.code == "H" {
			if firstH == 0 {
				break
			}
			firstH = 0
		} else if c.code == "S" {
			for i := 0; i < c.value; i++ {
				offset += 1
			}
		} else {
			return map[int]string{}, fmt.Errorf("Unsupported code in CIGAR.")
		}
	}

	return hits, nil

}

func GetVariants(chrom string, gp map[int]string, vars *Variants) map[int]string {

	if _, ok := (*vars).data[chrom]; !ok {
		//e := fmt.Errorf("Variant set didn't contain any variants for contig %s\n", chrom)
		//fmt.Println(e)
		//fmt.Printf("couldn't find chromosome %s in variant set\n", chrom)
		return map[int]string{} // no variants for this chromosome...
	}

	for pos, v := range gp {
		if _, ok := (*vars).data[chrom][pos]; !ok {
			//fmt.Printf("didn't find position %d in variant set\n", pos)
			delete(gp, pos)
		} else {
			// mark whether the variant is reference, alternate, or something else
			//fmt.Println((*vars)[chrom][pos])
			if v == (*vars).data[chrom][pos].Ref {
				gp[pos] = "R"
				//fmt.Printf("Reference: %s %s %d %s\n", v, chrom, pos, (*vars)[chrom][pos][3])
			} else if v == (*vars).data[chrom][pos].Alt {
				gp[pos] = "A"
				//fmt.Printf("Alternate: %s %s %d %s\n", v, chrom, pos, (*vars)[chrom][pos][4])
			} else {
				gp[pos] = "N"
				//fmt.Printf("Other: %s %s %d\n", v, chrom, pos)
			}
		}
	}

	return gp

}

// GetVariantBlocks takes a map of positions and sequence at those positions and checks
// the Variants object to determine which variant block is indicated at each such position.
//
//
func GetVariantBlocks(chrom string, gp map[int]string, vars *Variants) map[int][]int {

	ret := make(map[int][]int)

	if len((*vars).data) == 0 {
		panic("The variants object should never be empty. There may have been an error in parsing the variants vcf.")
	}

	if _, ok := (*vars).data[chrom]; !ok {
		return map[int][]int{}
	}

	for pos, v := range gp {

		if _, ok := (*vars).data[chrom][pos]; !ok {
			delete(gp, pos)
		} else {
			//fmt.Println("met cond 1")
			//fmt.Println("chrom, pos:", chrom, pos, v)
			//fmt.Println((*vars).data[chrom][pos].Ref, (*vars).data[chrom][pos].Alt)
			if v == (*vars).data[chrom][pos].Ref {
				//fmt.Println("calling ref")
				if blockNum, ok := (*vars).data[chrom][pos].Info["BLOCK"]; ok {
					//fmt.Println("met condition")
					n, _ := strconv.Atoi(blockNum)
					if _, ok := ret[int(n)]; !ok {
						ret[n] = []int{0, 0, 0}
					}
					ret[n][0] += 1
				}
			} else if v == (*vars).data[chrom][pos].Alt {
				//fmt.Println("calling alt")
				if blockNum, ok := (*vars).data[chrom][pos].Info["BLOCK"]; ok {
					//fmt.Println("met condition")
					n, _ := strconv.Atoi(blockNum)
					if _, ok := ret[n]; !ok {
						ret[n] = []int{0, 0, 0}
					}
					ret[n][1] += 1
				}
			} else {
				//fmt.Println("calling N")
				if blockNum, ok := (*vars).data[chrom][pos].Info["BLOCK"]; ok {
					//fmt.Println("met condition")
					n, _ := strconv.Atoi(blockNum)
					if _, ok := ret[n]; !ok {
						ret[n] = []int{0, 0, 0}
					}
					ret[n][2] += 1
				}
			}
		}
	}

	//fmt.Println(ret)

	return ret

}

// VariantLinksFromSam parses a sam file, constructing a Links object
// representing simple counts of association between variants.
func VariantLinksFromSam(samPath, outPath string, vars Variants) {

	out, err1 := os.Create(outPath)
	if err1 != nil {
		fmt.Printf("Couldn't open output file (%s) for reading: %s\n", outPath, err1)
	}
	defer out.Close()

	in, err2 := os.Open(samPath)
	if err2 != nil {
		fmt.Printf("Couldn't open input file (%s) for reading: %s\n", samPath, err2)
	}
	defer in.Close()

	links := NewLinks()
	currentID := ""
	hits := []Alignment{}
	s := bufio.NewScanner(in)
	linkedVariants := 0
	lineCount := 0
	balance := 0
	//maxlines := 1000000

	for s.Scan() {
		line := s.Text()
		lineCount += 1
		/*if lineCount >  maxlines {
			break
		}*/
		if string(line[0]) != "@" {

			a, e := parseSAMLine(s.Text())
			if e != nil {
				//fmt.Printf("error parsing sam line, continuing: %s, %s\n", e, line)
				hits = []Alignment{}
				currentID = ""
				continue
			}

			if a.qname != currentID {
				if (len(hits) == 2) && (len(currentID) > 0) {
					if hits[0].rname == hits[1].rname { // same chromosome

						gp1, egp1 := GenomePositions(hits[0])
						gp2, egp2 := GenomePositions(hits[1])
						if (egp1 != nil) || (egp2 != nil) {
							//fmt.Println("continuing")
							//fmt.Println(gp1)
							//fmt.Println(gp2)
							hits = []Alignment{}
							currentID = ""
							continue
						}
						//fmt.Println(hits[0].rname)
						varPos1 := GetVariants(hits[0].rname, gp1, &vars)
						varPos2 := GetVariants(hits[0].rname, gp2, &vars)
						// requires both on same chromosome, otherwise bug in following
						ct1, bal1 := links.TabulateVariantLinks(hits[0].rname, varPos1, varPos2)
						ct2, bal2 := links.TabulateVariantLinks(hits[0].rname, varPos1, varPos1)
						ct3, bal3 := links.TabulateVariantLinks(hits[0].rname, varPos2, varPos2)
						ct := ct1 + ct2 + ct3
						bal := bal1 + bal2 + bal3
						if ct > 0 {
							//fmt.Println(varPos1)
							//fmt.Println(varPos2)
							//fmt.Println(hits)
							linkedVariants += ct
							balance += bal
							fmt.Printf("linked a total of %d variants in %d lines with balance %d\n", linkedVariants, lineCount, balance)
						}
						//fmt.Println(varPos1)
						//fmt.Println(varPos2)
					}
				}
				currentID = a.qname
				hits = []Alignment{}
			}

			hits = append(hits, a)

		}
	}

	links.Write(out)

}

type varBlock struct {

	// The name of the variant block
	blockName string

	// Whether the alignments to the variant block were called as reference, alternate, or neither
	// 0 for reference, 1 for alternate, and 2 for neither
	phase int

	refCount int // The number of reference variants matched in the block
	altCount int // The number of alternate variants matched in the block

}

/*

#CHROM	POS	ID	REF	ALT	QUAL	FILTER	INFO	FORMAT	NA12878
chr1	2	rs1234	T	G	757.12	PASS	SNP:99;BLOCK:b1	GT	0/1
chr1	5	rs1235	T	C	757.12	PASS	SNP:99;BLOCK:b1	GT	0/1
chr1	24	rs1236	G	T	757.12	PASS	SNP:99;BLOCK:b2	GT	0/1
chr1	27	rs1237	A	G	757.12	PASS	SNP:99;BLOCK:b2	GT	0/1
chr1	42	rs1238	T	A	757.12	PASS	SNP:99;BLOCK:b3	GT	0/1
chr1	46	rs1239	A	G	757.12	PASS	SNP:99;BLOCK:b3	GT	0/1
chr1	53	rs1240	G	C	757.12	PASS	SNP:99;BLOCK:b4	GT	0/1
chr1	57	rs1241	C	T	757.12	PASS	SNP:99;BLOCK:b4	GT	0/1
chr5	57	rs1242	C	T	757.12	PASS	SNP:99;BLOCK:b5	GT	0/1
chr5	58	rs1243	C	T	757.12	PASS	SNP:99;BLOCK:b5	GT	0/1

readid-1.1	16	chr1	4	60	3M	*	0	0	ATG	III
readid-1.2	16	chr1	23	60	3M	*	0	0	GTA	III
readid-2.1	16	chr1	41	60	3M	*	0	0	GAC	III
readid-2.2	16	chr1	52	60	3M	*	0	0	GCG	III
readid-3.1	16	chr1	1	60	3M	*	0	0	AGG	III
readid-3.2	16	chr1	23	60	3M	*	0	0	GAA	III
readid-4.1	16	chr1	45	60	3M	*	0	0	GAC	III
readid-4.2	16	chr1	56	60	3M	*	0	0	GCG	III
readid-5.1	16	chr1	100	60	3M	*	0	0	GAC	III
readid-5.2	16	chr2	200	60	3M	*	0	0	GCG	III
readid-6.1	16	chr2	100	60	3M	*	0	0	GAC	III
readid-6.2	16	chr1	200	60	3M	*	0	0	GCG	III
readid-7.1	16	chr2	100	60	3M	*	0	0	GAC	III
readid-7.2	16	chr3	200	60	3M	*	0	0	GCG	III

*/

// VariantLinksFromSam parses a sam file, constructing a Links object
// representing simple counts of association between variants.
func BlockLinksFromSam(samPath, outPath string, vars Variants) {

	out, err1 := os.Create(outPath)
	if err1 != nil {
		fmt.Printf("Couldn't open output file (%s) for reading: %s\n", outPath, err1)
	}
	defer out.Close()

	in, err2 := os.Open(samPath)
	if err2 != nil {
		fmt.Printf("Couldn't open input file (%s) for reading: %s\n", samPath, err2)
	}
	defer in.Close()

	links := NewLinks()
	currentID := ""
	hits := []Alignment{}
	s := bufio.NewScanner(in)
	linkedVariants := 0
	lineCount := 0
	balance := 0
	//maxlines := 1000000

	for s.Scan() {
		line := s.Text()
		lineCount += 1
		/*if lineCount >  maxlines {
			break
		}*/
		if string(line[0]) != "@" {

			a, e := parseSAMLine(s.Text())
			if e != nil {
				//fmt.Printf("error parsing sam line, continuing: %s, %s\n", e, line)
				hits = []Alignment{}
				currentID = ""
				continue
			}

			if a.qname != currentID {
				if (len(hits) == 2) && (len(currentID) > 0) {
					if hits[0].rname == hits[1].rname { // same chromosome

						gp1, egp1 := GenomePositions(hits[0])
						gp2, egp2 := GenomePositions(hits[1])

						//fmt.Println(gp1, gp2)

						if (egp1 != nil) || (egp2 != nil) {
							//fmt.Println("continuing")
							//fmt.Println(gp1)
							//fmt.Println(gp2)
							hits = []Alignment{}
							currentID = ""
							continue
						}

						// Given the above genome positions, first map the variants at those positions to
						// known variant blocks.

						// The trick to parsing block links instead of variant links is passing a varpos
						// object into the subsequent step that indicates whether each block that was inter-
						// sected was so at a reference or alternate position.
						varPos1 := GetVariantBlocks(hits[0].rname, gp1, &vars)
						varPos2 := GetVariantBlocks(hits[1].rname, gp2, &vars)

						// requires both on same chromosome, otherwise bug in following
						ct1, bal1 := links.TabulateVariantBlockLinks(hits[0].rname, varPos1, varPos2)
						//fmt.Println(ct1)
						ct := ct1
						bal := bal1
						if ct > 0 {
							//fmt.Println(varPos1)
							//fmt.Println(varPos2)
							//fmt.Println(hits)
							linkedVariants += ct
							balance += bal
							//fmt.Printf("linked a total of %d variants in %d lines with balance %d\n", linkedVariants, lineCount, balance)
						}
						//fmt.Println(varPos1)
						//fmt.Println(varPos2)
					}
				}
				currentID = a.qname
				hits = []Alignment{}
			}

			hits = append(hits, a)

		}
	}

	links.Write(out)

}

/*
// VariantLinksFromSam parses a sam file, constructing a Links object
// representing simple counts of association between variants.
func BlockLinksFromSam(samPath, outPath string, vars Variants) {

	out, err1 := os.Create(outPath)
	if err1 != nil {
		fmt.Printf("Couldn't open output file (%s) for reading: %s\n", outPath, err1)
	}
	defer out.Close()

	in, err2 := os.Open(samPath)
	if err2 != nil {
		fmt.Printf("Couldn't open input file (%s) for reading: %s\n", samPath, err2)
	}
	defer in.Close()

	links := NewLinks()
	currentID := ""
	hits := []Alignment{}
	s := bufio.NewScanner(in)
	linkedVariants := 0
	lineCount := 0

	for s.Scan() {
		line := s.Text()
		lineCount += 1
		if string(line[0]) != "@" {

			a, e := parseSAMLine(s.Text())

			if e != nil {
				hits = []Alignment{}
				currentID = ""
				continue
			}

			if a.qname != currentID {
				if (len(hits) == 2) && (len(currentID) > 0) {
					if hits[0].rname == hits[1].rname { // same chromosome

						blockLinks, e := BlockLinksForAlignmentPair(hits[0], hits[1], &vars)

						if e != nil {
							continue
						}

						links.AddMulti(blockLinks)

					}
				}

				currentID = a.qname
				hits = []Alignment{}
			}

			hits = append(hits, a)

		}
	}

	links.Write(out)

}*/

// ScaffoldLinksFromSam parses a sam file, constructing a Links object
// representing simple counts of association between contigs.
func ScaffoldLinksFromSam(samPath, outPath string) {

	// assumes sam is sorted by id
	// for now, building links map in memory but could
	// also write all to disk then sort | uniq -c

	out, err1 := os.Create(outPath)
	if err1 != nil {
		fmt.Printf("Couldn't open output file (%s) for reading: %s\n", outPath, err1)
	}
	defer out.Close()

	in, err2 := os.Open(samPath)
	if err2 != nil {
		fmt.Printf("Couldn't open input file (%s) for reading: %s\n", samPath, err2)
	}
	defer in.Close()

	links := NewLinks()
	currentID := ""
	chrHits := []string{}
	//posHits := []int{}
	s := bufio.NewScanner(in)

	for s.Scan() {
		line := s.Text()
		if string(line[0]) != "@" {

			arr := strings.Split(s.Text(), "\t")

			if arr[0] != currentID {

				if len(chrHits) == 2 {
					id1 := links.ID(chrHits[0])
					id2 := links.ID(chrHits[1])
					links.Add(id1, id2, 1)
				}

				currentID = arr[0]
				chrHits = []string{}
				//posHits = []string{}

			}

			chrHits = append(chrHits, arr[2])
			//posHits = append(posHits, arr[3])

		}
	}

	links.Write(out)

}

// PartitionAlignmentsByContig takes a path of an alignments file and partitions
// it into separate files with the alignments for each chromosome each in a separate
// file. The output is written to the output directory with the specified output stem
// with the chromosome name included, e.g. /path/to/file/outstem.chrNN.[sam/bam]
func PartitionAlignmentsByContig(inputPath, outstem string) error {

	// Set a limit on the number of allowed open files
	maxOpenFiles := 200

	// Initialize a map to store the output file handles
	fhMap := map[string]*os.File{}
	fhCount := 0

	// Open the input file for reading
	infh, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("error: PartitionAlignmentsByContig(%s, %s) failed at attempt to open input file", inputPath, outstem)
	}
	defer infh.Close()
	s := bufio.NewScanner(infh)

	for s.Scan() {

		// Get the line
		line := s.Text()
		if string(line[0]) != "#" {

			alignment, e := parseSAMLine(line)
			if e != nil {
				return e
			}

			// If the observed chromosome name is not in the name/file handle map,
			// and we have not reached the limit of open files, create a new file handle
			// and add it to the map.
			if _, ok := fhMap[alignment.rname]; !ok {
				if fhCount >= maxOpenFiles {
					return fmt.Errorf("Too many open files for PartitionAlignmentsByContig to proceed. Either increase the limit or try a different approach.")
				}
				outpath := filepath.Join(outstem, alignment.rname, ".vcf")
				f, err := os.Create(outpath)
				if err != nil {
					return err
				}
				defer f.Close()
				fhMap[alignment.rname] = f
				fhCount += 1
			}

			fh := fhMap[alignment.rname]
			fh.WriteString(alignment.String() + "\n")

		}

	}

	// Return nil, signaling an error free run
	return nil

}
