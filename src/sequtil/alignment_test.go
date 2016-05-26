package util

import (
	"reflect"
	"testing"
	//"path/filepath"
	"os"
)

func TestParseSAMFlag(t *testing.T) {

	if f, _ := parseSAMFlag("1"); !f.multiseg {
		t.Errorf("")
	}

	if f, _ := parseSAMFlag("2"); !f.allpropper {
		t.Errorf("")
	}

	if f, _ := parseSAMFlag("4"); !f.unmapped {
		t.Errorf("")
	}

	if f, _ := parseSAMFlag("8"); !f.nextunmapped {
		t.Errorf("")
	}

	if f, _ := parseSAMFlag("2048"); !f.supplementary {
		t.Errorf("")
	}

	if f, _ := parseSAMFlag("2049"); !f.supplementary || !f.multiseg {
		t.Errorf("")
	}

}

func TestParseCIGAR(t *testing.T) {

	test := map[string]CIGAR{
		"100M": CIGAR{CIGARCode{100, "M"}},
		"":     CIGAR{}, // not allowed
		//"H100": CIGAR{}, // Not allowed, expand to handle later
		//"100M100H100M": CIGAR{}, // Not allowed, expand to handle later
		"100H100M":     CIGAR{CIGARCode{100, "H"}, CIGARCode{100, "M"}},
		"100S100M100H": CIGAR{CIGARCode{100, "S"}, CIGARCode{100, "M"}, CIGARCode{100, "H"}},
	}

	for k, v := range test {

		c, e := parseCIGAR(k)
		if e != nil {
			if !reflect.DeepEqual(v, CIGAR{}) {
				t.Errorf("test error: unexpected CIGAR error case: %s\n", k)
			}
		} else {
			if !reflect.DeepEqual(c, v) {
				t.Errorf("test error: parsed and expected CIGAR structures do not match: %s, %s", c, v)
			}
		}

	}

}

func TestParseSAMLine(t *testing.T) {

	l1 := "SRR927086.7	2048	12	20766468	60	100H81M19H	*	0	0	CAAACGTGTGCACATCCNNGAGAGCCGTGAGCAACTTGCTCAGCANACNNCTCANCTTCCANGNCNTTCNCAAGCCCAGAG	<<<???@?@?@@@???<%%33=>???@??????????????????%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%	NM:i:10	MD:Z:17T0T26C2A0G4C6G1C1A3A11	AS:i:61	XS:i:0	SA:Z:10,52681560,+,100M100S,60,1;"
	sf, _ := parseSAMFlag("2048")
	a1 := Alignment{"SRR927086.7", sf, "12", 20766468, 60,
		CIGAR{CIGARCode{100, "H"}, CIGARCode{81, "M"}, CIGARCode{19, "H"}},
		"*", 0, 0,
		"CAAACGTGTGCACATCCNNGAGAGCCGTGAGCAACTTGCTCAGCANACNNCTCANCTTCCANGNCNTTCNCAAGCCCAGAG",
		"<<<???@?@?@@@???<%%33=>???@??????????????????%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%",
	}

	l2 := "SRR927086.6	0	11	95204105	0	99S101M	*	0	0	CAGGACATANGCGNNNGCAAGGACTTCATGTCCAAAACACCAAAAGCAATGGCAACAAAAGCCAAAATTGACAAATGAGATCTAATTAAACTAAAGAGCTTCTATATCTCTGTTTTGGTACCAGTACCATGCTGTTTTGGTTACTGTAGCCTTGTAGTATAGTTTGAAGTCAGGTAGTGTGATGCCTCCAGCTTTGTTCN	%%%%%%%%%%%%%%%%EEEHCHHHHDC@CDIIGGHHFGDGEGFEGGIIIH@IIHGIIIIIGIIIIHBCIIHIIHF>FIHHBHGBGIIFFFBD?BDDD?@@B@>B@B?B>@B@BBCCCCABDDDDB;FECFFIIIIFCGF@FBBGIFDBIIIIGEIGBDDFFEFGEGEG9IFFF;EFF<CCA:F?FCBDDFFDDDDDB:1NM:i:1	MD:Z:100T0	AS:i:100	XS:i:100	SA:Z:8,98312769,-,103M97S,0,4;"
	a2 := Alignment{
		"SRR927086.6", SAMFlag{}, "11", 95204105, 0,
		CIGAR{CIGARCode{99, "S"}, CIGARCode{101, "M"}},
		"*", 0, 0,
		"CAGGACATANGCGNNNGCAAGGACTTCATGTCCAAAACACCAAAAGCAATGGCAACAAAAGCCAAAATTGACAAATGAGATCTAATTAAACTAAAGAGCTTCTATATCTCTGTTTTGGTACCAGTACCATGCTGTTTTGGTTACTGTAGCCTTGTAGTATAGTTTGAAGTCAGGTAGTGTGATGCCTCCAGCTTTGTTCN",
		"%%%%%%%%%%%%%%%%EEEHCHHHHDC@CDIIGGHHFGDGEGFEGGIIIH@IIHGIIIIIGIIIIHBCIIHIIHF>FIHHBHGBGIIFFFBD?BDDD?@@B@>B@B?B>@B@BBCCCCABDDDDB;FECFFIIIIFCGF@FBBGIFDBIIIIGEIGBDDFFEFGEGEG9IFFF;EFF<CCA:F?FCBDDFFDDDDDB:1NM:i:1",
	}

	check := map[string]Alignment{
		l1: a1,
		l2: a2,
	}

	for k, v := range check {
		a, e := parseSAMLine(k)
		if e != nil {
			t.Errorf("test error: error when attempting to parse sam line: %s\n", e)
		}
		if !reflect.DeepEqual(a, v) {
			t.Errorf("test error: parsed sam line does not match key")
		}
	}

}

func TestGenomePositions(t *testing.T) {

	a1 := Alignment{"SRR927086.7", SAMFlag{}, "12", 20766468, 60,
		CIGAR{CIGARCode{100, "H"}, CIGARCode{4, "M"}, CIGARCode{19, "H"}},
		"*", 0, 0,
		"CAAA",
		"<<<???@?@?@@@???<%%33=>???@??????????????????%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%",
	}

	a2 := Alignment{
		"SRR927086.6", SAMFlag{}, "11", 95204105, 0,
		CIGAR{CIGARCode{1, "S"}, CIGARCode{3, "M"}},
		"*", 0, 0,
		"CAGG",
		"%%%%%%%%%%%%%%%%EEEHCHHHHDC@CDIIGGHHFGDGEGFEGGIIIH@IIHGIIIIIGIIIIHBCIIHIIHF>FIHHBHGBGIIFFFBD?BDDD?@@B@>B@B?B>@B@BBCCCCABDDDDB;FECFFIIIIFCGF@FBBGIFDBIIIIGEIGBDDFFEFGEGEG9IFFF;EFF<CCA:F?FCBDDFFDDDDDB:1NM:i:1",
	}
	a := []Alignment{a1, a2}

	gp1 := map[int]string{
		20766468: "C",
		20766469: "A",
		20766470: "A",
		20766471: "A",
	}

	gp2 := map[int]string{
		95204106: "A",
		95204107: "G",
		95204108: "G",
	}
	gp := []map[int]string{gp1, gp2}

	for i, v := range gp {
		check, _ := GenomePositions(a[i])
		if !reflect.DeepEqual(check, v) {
			t.Errorf("")
		}
	}

}

func TestGetVariants(t *testing.T) {

	header := map[string][]string{
		"samples": []string{
			"NA12878",
		},
	}

	varData := map[string]map[int]variant{
		"1": map[int]variant{
			20766468: variant{ // position 0 in 0-based
				"1", "20766468", "rs1", "A", "T", "0", "PASS", "SNP:NN", nil,
			},
			20766470: variant{
				"1", "20766470", "rs2", "G", "A", "0", "PASS", "SNP:NN", nil,
			},
			95204106: variant{ // position 2 in 0-based
				"1", "95204106", "rs3", "C", "G", "0", "PASS", "SNP:NN", nil,
			},
			95204107: variant{ // position 2 in 0-based
				"1", "95204107", "rs4", "T", "G", "0", "PASS", "SNP:NN", nil,
			},
		},
	}

	vars := Variants{
		header,
		varData,
	}

	gp1 := map[int]string{
		20766468: "C",
		20766469: "A",
		20766470: "A",
		20766471: "A",
	}

	gp2 := map[int]string{
		95204106: "C",
		95204107: "G",
		95204108: "G",
	}

	varpos1 := GetVariants("1", gp1, &vars)
	varpos2 := GetVariants("1", gp2, &vars)

	key1 := map[int]string{20766468: "N", 20766470: "A"}
	key2 := map[int]string{95204106: "R", 95204107: "A"}

	if !reflect.DeepEqual(key1, varpos1) || !reflect.DeepEqual(key2, varpos2) {
		t.Errorf("test error: error filtering positions with variants")
	}

}

/*
func TestVariantBlocksForPositions(t *testing.T) {

	header := map[string][]string {
		"samples": []string {
			"NA12878",
		},
	}

	additional1 := map[string]map[string]string {
		"NA12878": map[string]string{
			"GT": "0/1",
		},
	}

	additional2 := map[string]map[string]string {
		"NA12878": map[string]string{
			"GT": "1/0",
		},
	}

	// The first block is in phase with the reference while the second block has an inversion
	// between 95204106 and 95204107 relative to the reference.
	//
	// 				reference 	alternate 	NA12878
	// 20766468			A 			T 		 A  T 		<- b1
	// 20766470			G 			A 		 G  A 		<- b1
	//
	// 95204106			C 			G 		 G  C 		<- b2
	// 95204107 		T 			G 		 T  G 		<- b2

	varData := map[string]map[int]variant {
		"1": map[int]variant {
			20766468: variant{ 	// position 0 in 0-based
				"1", "20766468", "rs1", "A", "T", "0", "PASS", "SNP:NN;BLOCK:b1", additional1,
			},
			20766470: variant{
				"1", "20766470", "rs2", "G", "A", "0", "PASS", "SNP:NN;BLOCK:b1", additional1,
			},
			95204106: variant{		// position 2 in 0-based
				"1", "95204106", "rs3", "C", "G", "0", "PASS", "SNP:NN;BLOCK:b2", additional2,
			},
			95204107: variant{		// position 2 in 0-based
				"1", "95204107", "rs4", "T", "G", "0", "PASS", "SNP:NN;BLOCK:b2", additional1,
			},
		},
	}

	vars := Variants{
		header,
		varData,
	}

	gp1 := map[int]string {
		20766468: "C",
		20766469: "A",
		20766470: "A",
		20766471: "A",
	}

	gp2 := map[int]string {
		95204106: "C",
		95204107: "G",
		95204108: "G",
	}

	varpos1 := VariantBlocksForPositions("1", gp1, &vars)
	varpos2 := VariantBlocksForPositions("1", gp2, &vars)

	key1 := []varBlock{
			varBlock{},
		}
	key2 := []varBlock{
			varBlock{block2, phase2, ref2, alt2},
		}

	if !reflect.DeepEqual(key1, varpos1) || !reflect.DeepEqual(key2, varpos2) {
		t.Errorf("test error: error filtering positions with variants")
	}

}
*/

/*
func TestBlockLinksFromSam(t *testing.T) {

	samPath := filepath.Join(cwd(t), "_testdata", "toy.sam")
	outPath := filepath.Join(cwd(t), "_testdata", "output.var.links")
	vars, _ := ReadVariants(filepath.Join(cwd(t), "_testdata", "toy.vcf"))

	BlockLinksFromSam(samPath, outPath, vars)

	linksTest, _ := LoadLinks(outPath)
	linksKey, _ := LoadLinks(filepath.Join(cwd(t), "_testdata", "toy.varblock.links"))

	if !reflect.DeepEqual(linksTest, linksKey) {
		t.Errorf("VariantBlockLinksFromSAM(%s, %s) yielded links non-identical to contig links key", samPath, outPath)
	}

}
*/

/*
func TestVariantLinksFromSAM(t *testing.T) {

	samPath := filepath.Join(cwd(t), "_testdata", "toy.sam")
	outPath := filepath.Join(cwd(t), "_testdata", "output.var.links")
	vars, _ := ReadVariants(filepath.Join(cwd(t), "_testdata", "toy.var.links"))

	VariantLinksFromSam(samPath, outPath, vars)

	linksTest, _ := LoadLinks(outPath)
	linksKey, _ := LoadLinks(filepath.Join(cwd(t), "_testdata", "toy.var.links"))

	if !reflect.DeepEqual(linksTest, linksKey) {
		t.Errorf("VariantLinksFromSam(%s, %s) yielded links non-identical to contig links key", samPath, outPath)
	}

}

func TestScaffoldLinksFromSam(t *testing.T) {

	samPath := filepath.Join(cwd(t), "_testdata", "toy.sam")
	outPath := filepath.Join(cwd(t), "_testdata", "output.ctg.links")

	ScaffoldLinksFromSam(samPath, outPath)

	linksTest, _ := LoadLinks(outPath)
	linksKey, _ := LoadLinks(filepath.Join(cwd(t), "_testdata", "toy.ctg.links"))

	if !reflect.DeepEqual(linksTest, linksKey) {
		t.Errorf("ScaffoldLinksFromSam(%s, %s) yielded links non-identical to contig links key", samPath, outPath)
	}

}
*/

/* A better structure for tests is the following:
tests := []struct {
	path string
	err error
}{
	{path: filepath.Join(cwd(t), "_testdata", "splitme.sam"), err: nil},
	{path: filepath.Join(cwd(t), "_testdata", "splitme.sam1"), err: fmt.Errorf('error: input file does not exist, %s', filepath.Join(cwd(t), "_testdata", "splitme.sam1")},
}
for _, tt := range tests {
	err := PartitionAlignmentsByContig(tt.path, outstem)
	if !reflect.DeepEqual(err, tt.err){
		t.Errorf("PartitionAlignmentsByChromosome(%s) failed with error %s, expected error %s", tt.path, err, tt.err)
	}
}
*/

/*
func TestPartitionAlignmentsByContig(t *testing.T) {

	// Get the path of the test data
	outdir := filepath.Join(cwd(t), "_testdata", "split")
	if err := os.MkdirAll(outdir, 0777); err != nil {
		t.Error(err)
	}

	outstem := filepath.Join(outdir, "splitme")

	inpath := filepath.Join(cwd(t), "_testdata", "splitme.sam")
	if e := PartitionAlignmentsByContig(inpath, outstem); e != nil {
		t.Errorf("PartitionAlignmentsByContig(%s)", inpath)
	}
	if e := PartitionAlignmentsByContig(filepath.Join(cwd(t), "_testdata", "splitme.sam1"), outstem); e == nil {
		t.Errorf("Got nil error when non nil was expected, i.e. when operating on a non-existent file.", e)
	}

	// Verify partitioning by partitioning the toy files by hand and comparing each of the resulting
	// subset files.

}
*/

func cwd(t *testing.T) string {
	cwd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	return cwd
}
