package util

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Variants is an object that stores a set of genome sequence
// variants indexed according to chromosome name and position.
type Variants struct {

	// Header field to store information on sample ID's and other key/value information
	header map[string][]string

	// Data field to store the mapping between chromosome, position, and variant record
	data map[string]map[int]variant
}

// variant is an object that represents a single genome sequence
// variant.
type variant struct {
	Chrom  string            // chromosome
	Pos    string            // position
	ID     string            // identifier
	Ref    string            // reference base(s)
	Alt    string            // alternate base(s)
	Qual   string            // quality
	Filter string            // filter status (PASS vs. semicolon-delmited list of failing filters)
	Info   map[string]string // additional information (semicolon-delimited list)

	// A mapping from sample ID to one of multiple fields, including GT for genotype, e.g.
	// {sample1: {GT: 0/1, GQ: 1}, sample2: {GT: 1/1, GQ: 1}}
	Additional map[string]map[string]string
}

func NewVariants() Variants {
	header := make(map[string][]string)
	data := make(map[string]map[int]variant)
	return Variants{header, data}
}

// String renders an alignment object as a single SAM format string.
func (v *variant) String() string {

	infostring := ""
	state := 0
	for k, v := range v.Info {
		infostring = infostring + k + "=" + v
		if state != 0 {
			infostring = infostring + ";"
		}
		state = 1
	}

	str := strings.Join([]string{v.Chrom, v.Pos, v.ID, v.Ref, v.Alt, v.Qual, v.Filter, infostring}, "\t")

	if len(v.Additional) > 0 {

		additional := "GT"

		for _, val := range (*v).Additional {

			field := ""

			for _, val2 := range val {
				if len(field) == 0 {
					field = val2
				} else {
					field = field + ":" + val2
				}
			}

			additional = additional + "\t" + field

		}

		str = str + "\t" + additional

	}

	return str
}

func (v *variant) ParseInfo(infoLine string) error {

	info := map[string]string{}

	pairs := strings.Split(infoLine, ";")
	for _, p := range pairs {
		if len(strings.Split(p, "=")) == 2 {
			arr := strings.Split(p, "=")
			info[arr[0]] = arr[1]
		} else {
			info["other"] = p
		}
	}

	v.Info = info

	return nil

}

// readVariant reads a single VCF line into a variant object
func readVariant(line string, samples []string) (variant, error) {

	// Check that the input line does not start with a hash, which means that it's a header line
	if string(line[0]) == "#" {
		return variant{}, fmt.Errorf("Can't parse header line into variant object: %s", line)
	}

	// Split the input line with tab delimiters
	arr := strings.Split(line, "\t")

	// Check that at least the minimum 8 fields are present
	if len(arr) < 8 {
		return variant{}, fmt.Errorf("Variant line does not include the required minimum number of fields.")
	}

	// Build a variant object to hold the eight mandatory fields
	v := variant{arr[0], arr[1], arr[2], arr[3], arr[4], arr[5], arr[6], nil, nil}
	v.ParseInfo(arr[7])

	// If we only have the mandatory 8 columns of information, return the record.
	//if len(arr) == 8 {
	//	return v, nil
	//}

	// If we have more than 9 fields the samples array provided should have a number of
	// entries equal to the number of non mandatory columns minus 1
	//if len(arr)-9 != len(samples) {
	//	return variant{}, fmt.Errorf("tried to parse a VCF line with non-mandatory sample fields but was provided a number of sample ids that did not match the number of columns in the line being parsed:\n %s", line)
	//}

	/*
		// Parse the format column
		format := strings.Split(arr[8], ":")

		// Instantiate the map for the additional column information
		additional := make(map[string]map[string]string)

		// For each sample
		for i, s := range samples {

			// Allocate the submap to store the information for this sample
			additional[s] = map[string]string{}

			// Calculate the index into the split line array for this sample
			ind := 8 + i + 1

			// Split the corresponding record (i.e. occurring at column 'ind')
			formatArr := strings.Split(arr[ind], ":")

			// For each format field, save it in the map
			for j, f := range formatArr {

				// Here format[j] is something like GT, GQ, ..., s is a sample ID like
				// NA12878, and f is something like 0/1 for GT, 1 for GQ, etc.
				form := format[j]
				additional[s][form] = f

			}

		}

		v.Additional = additional
	*/

	// Return the result
	return v, nil

}

// add adds a single variant to a variant set object
func (vars *Variants) add(v variant) {
	chr := v.Chrom
	pos, _ := strconv.Atoi(v.Pos)
	if _, ok := (*vars).data[chr]; !ok {
		(*vars).data[chr] = map[int]variant{}
	}
	(*vars).data[chr][pos] = v
}

func parseSamples(line string) ([]string, error) {

	arr := strings.Split(line, "\t")

	// If we have fewer than the eight mandatory fields, the file is malformatted and
	// we will return an error.
	if len(arr) < 8 {
		return []string{}, fmt.Errorf("The provided file has a header line with fewer than the eight mandatory columns (%d columns were observed)", len(arr))
	}
	if len(arr) == 9 {
		return []string{}, fmt.Errorf("Observed greater than the eight mandatory columns but not more than 9, suggesting the input file is malformatted (i.e. a specification is provided in field 9 referencing fields 10+ which are not present.")
	}

	// Allocate a slice of the right size to store the resulting list of samples
	// and copy the sample ids to it.
	samples := make([]string, len(arr)-9)
	for i := 0; i < len(arr)-9; i++ {
		samples[i] = arr[(i + 9)]
	}

	return samples, nil

}

// ReadVariants loads a set of variants specified in a VCF file into a Variants struct
func ReadVariants(path string) (Variants, int, error) {

	vars := NewVariants()

	in, err := os.Open(path)
	if err != nil {
		fmt.Errorf("Couldn't open VCF file %s\n", path)
		return vars, 0, err
	}
	defer in.Close()
	s := bufio.NewScanner(in)

	samples := []string{}

	for s.Scan() {

		line := s.Text()

		if string(line[0]) == "#" {

			// When we see the line that begins with only one # symbol, it's the line that
			// has the list of sample ID's - parse it to obtain the sample ID array to be used
			// in subsequent variant parsing.
			if string(line[1]) != "#" {
				samples, err = parseSamples(line)
				if err != nil {
					return Variants{}, 0, err
				}
			}

		} else {
			v, e := readVariant(line, samples)
			//fmt.Println(v)
			if e != nil {
				fmt.Println(e)
				return Variants{}, 0, e
			}
			vars.add(v)
		}

	}

	vars.header = map[string][]string{
		"samples": samples,
	}

	num := len(vars.data)

	return vars, num, err

}

// TODO: write variants in appropriate sorted order
// WriteVariants writes a variant object to a VCF file on disk
func WriteVariants(v Variants, path string) error {

	MkdirForFile(path)

	out, err := os.Create(path)
	if err != nil {
		fmt.Errorf("Couldn't open VCF file %s\n", path)
	}
	defer out.Close()

	// Construct the header string, including the name of each sample if there is this
	// information in the header field of the variants object
	header := "#CHROM	POS	ID	REF	ALT	QUAL	FILTER	INFO"
	if _, ok := v.header["samples"]; ok {
		header = fmt.Sprintf("%s\t%s", header, "FORMAT")
		for _, s := range v.header["samples"] {
			header = header + "\t" + s
		}
	}

	// Write the header line to the output
	out.WriteString(header + "\n")

	// Print VCF body
	for _, value := range v.data {
		for _, value2 := range value {
			out.WriteString(value2.String() + "\n")
		}
	}

	return err

}

// SimPhasedBlocks take a path to a set of sequence variants with known phase and simulates
// blocks of known phase of a specified block and gap size.
func SimPhasedBlocks(varPath, outPath string, size, gap int) error {

	// Note: Assumes the input is sorted.
	// TODO: Add a check for whether a provided input is indeed sorted.

	// Open the output file for writing
	out, err := os.Create(outPath)
	if err != nil {
		fmt.Errorf("Couldn't create output file %s to which to write variant block simulation\n", outPath)
	}
	defer out.Close()

	// Open the input file and instantiate a bufio reader
	in, err2 := os.Open(varPath)
	if err2 != nil {
		fmt.Errorf("Couldn't open input variants file for variant block simulation\n", outPath)
	}
	defer in.Close()

	s := bufio.NewScanner(in)

	// Initialize a counter to keep track of the start of the current block
	currentStart := 0
	currentEnd := size
	blockNum := 1
	samples := []string{}
	currentChrom := ""
	sawVariantsInBlock := false

	wroteHeader := false
	header := "#CHROM	POS	ID	REF	ALT	QUAL	FILTER	INFO"

	// For each line in the input variants file, determine the block to which it belongs
	// and write to the output.
	for s.Scan() {

		// Get the line
		line := s.Text()
		if string(line[0]) == "#" && string(line[1]) != "#" {
			samples, err = parseSamples(line)
			if err != nil {
				return err
			}

			header = header + "\t" + "FORMAT"
			for _, sample := range samples {
				header = header + "\t" + sample
			}
			out.WriteString(header + "\n")
			wroteHeader = true
		} else if string(line[0]) != "#" {

			if !wroteHeader {
				out.WriteString(header + "\n")
				wroteHeader = true
			}

			variant, _ := readVariant(line, samples)
			varPos, _ := strconv.Atoi(variant.Pos)

			// If we are either at the first record or the start of a new chromosome,
			// increment the block counter and reset the start and end positions. Also
			// note that the currentChrom has changed.
			if variant.Chrom != currentChrom {

				if sawVariantsInBlock {
					blockNum += 1
				}
				currentStart = 0
				currentEnd = size
				currentChrom = variant.Chrom
				sawVariantsInBlock = false
			}

			// If the position of the current variant is beyond the current block,
			// advance the current block until either the variant falls within the
			// block or the start of the block occurs after the variant.
			for {

				if varPos < currentStart {
					break
				} else if varPos < currentEnd {
					sawVariantsInBlock = true
					block := fmt.Sprintf("BLOCK:b%d", blockNum)
					//variant.Info = variant.Info + ";" + block
					variant.Info["BLOCK"] = string(block)
					out.WriteString(variant.String() + "\n")
					break
				}

				currentStart, currentEnd = currentEnd+gap, currentEnd+gap+size

				if sawVariantsInBlock {
					blockNum += 1
				}
			}
		}
	}

	return nil
}

func PartitionVariantsByContig(inputPath, outstem string) error {

	// Set a limit on the number of allowed open files
	maxOpenFiles := 200

	// Initialize a map to store the output file handles
	fhMap := map[string]*os.File{}
	fhCount := 0

	// Initialize an appropriate header to be written to output files on first open
	//header := "test"

	// Open the input file for reading
	infh, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("error: PartitionVariantsByContig(%s, %s) failed at attempt to open input file", inputPath, outstem)
	}
	defer infh.Close()
	s := bufio.NewScanner(infh)

	samples := []string{}

	for s.Scan() {

		// Get the line
		line := s.Text()
		if string(line[0]) == "#" && string(line[1]) != "#" {
			e := fmt.Errorf("")
			samples, e = parseSamples(line)
			if e != nil {
				return e
			}
		} else if string(line[0]) != "#" {

			v, _ := readVariant(line, samples)

			// If the observed chromosome name is not in the name/file handle map,
			// and we have not reached the limit of open files, create a new file handle
			// and add it to the map.
			if _, ok := fhMap[v.Chrom]; !ok {
				if fhCount >= maxOpenFiles {
					return fmt.Errorf("Too many open files for PartitionVariantsByContig to proceed. Either increase the limit or try a different approach.")
				}
				outpath := filepath.Join(outstem, v.Chrom, ".vcf")
				f, e := os.Create(outpath)
				if e != nil {
					return e
				}
				defer f.Close()
				fhMap[v.Chrom] = f
				fhCount += 1
			}

			fh := fhMap[v.Chrom]
			fh.WriteString(v.String() + "\n")

		}

	}

	return nil

}
