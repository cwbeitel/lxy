package util

import (
	"testing"
	"fmt"
	"reflect"
	"path/filepath"
)


func TestVariantIO(t *testing.T) {
	
	//fmt.Println("testing: util/WriteVariants and util/ReadVariants")

	path := "/tmp/lxy/test/testvariantio.vcf"

	// Build a variant header object
	header := map[string][]string{
		"samples": []string {
			"NA12878",
		},
	} 

	additional := map[string]map[string]string {
		"NA12878": map[string]string{
			"GT": "0/1",
		},
	}

	// Build a map of variant data
	vardata := map[string]map[int]variant{
		"chr1": {
			1: variant{
				Chrom: "chr1",
				Pos: "1",
				ID: "rs1",
				Ref: "C",
				Alt: "A",
				Qual: "0",
				Filter: "PASS",
				Info: "SNP:99;BLOCK=b1", 
				Additional: additional,
			},
			5: variant{"chr1", "5", "rs2", "C", "A", "0", "PASS", "SNP:99", additional},
			10: variant{"chr1", "10", "rs3", "C", "A", "0", "PASS", "SNP:99;BLOCK=b1", additional},
		},
		"chr2": {
			1: variant{"chr2", "1", "rs4", "C", "A", "0", "PASS", "SNP:99;BLOCK=b2", additional},
			12: variant{"chr2", "12", "rs5", "C", "A", "0", "PASS", "SNP:99;BLOCK=b2", additional},
			4: variant{"chr2", "4", "rs6", "C", "A", "0", "PASS", "SNP:99", additional},
		},
	}

	v := Variants{header, vardata}

	err := WriteVariants(v, path)
	if err != nil {
		t.Errorf("test error: TestVariantIO, WriteVariants failed\n")
	}

	v2, err2 := ReadVariants(path)
	if err2 != nil {
		t.Errorf("test error: TestVariantIO, ReadVariants failed\n")
	}

	if !reflect.DeepEqual(v, v2) {
		t.Errorf("test error: variant objects before and after reading do not match")
		fmt.Println(v)
		fmt.Println(v2)
	}

}


func TestReadVariant(t *testing.T) {

	additional := map[string]map[string]string {
		"NA12878": map[string]string{
			"GT": "0/1",
		},
	}
	v := variant{"chr1", "2", "rs1234", "T", "G", "757.12", "PASS", "SNP:99;BLOCK=b1", additional}

	vcfLine := "chr1	2	rs1234	T	G	757.12	PASS	SNP:99;BLOCK=b1	GT	0/1"

	v2, e := readVariant(vcfLine, []string{"NA12878"})
	if e != nil {
		t.Error(e)
	}

	if !reflect.DeepEqual(v, v2) {
		t.Errorf("Loaded and key variant objects do not match.")
		fmt.Println(v)
		fmt.Println(v2)
	}

}


func TestReadVariants(t *testing.T) {

	// Test that variants can be read correctly from the toy.vcf file

	// Build the header portion of the key variants object
	header := map[string][]string{
		"samples": []string{
			"NA12878",
		},
	}

	additional := map[string]map[string]string {
		"NA12878": map[string]string{
			"GT": "0/1",
		},
	}

	// Build the data portion of the key variants object
	varData := map[string]map[int]variant{
		"chr1": map[int]variant{
			2: variant{"chr1", "2", "rs1234", "T", "G", "757.12", "PASS", "SNP:99;BLOCK:b1", additional},
			5: variant{"chr1", "5", "rs1235", "T", "C", "757.12", "PASS", "SNP:99;BLOCK:b1", additional},
			24: variant{"chr1", "24", "rs1236",	"G", "T", "757.12", "PASS", "SNP:99;BLOCK:b2", additional},
			27: variant{"chr1", "27", "rs1237",	"A", "G", "757.12", "PASS", "SNP:99;BLOCK:b2", additional},
			42: variant{"chr1", "42", "rs1238",	"T", "A", "757.12", "PASS", "SNP:99;BLOCK:b3", additional},
			46: variant{"chr1", "46", "rs1239",	"A", "G", "757.12", "PASS", "SNP:99;BLOCK:b3", additional},
			53: variant{"chr1", "53", "rs1240",	"G", "C", "757.12", "PASS", "SNP:99;BLOCK:b4", additional},
			57: variant{"chr1", "57", "rs1241",	"C", "T", "757.12", "PASS", "SNP:99;BLOCK:b4", additional},
		},
		"chr5": map[int]variant{
			57: variant{"chr5", "57", "rs1242", "C", "T", "757.12", "PASS", "SNP:99;BLOCK:b5", additional},
			58: variant{"chr5", "58", "rs1243", "C", "T", "757.12", "PASS", "SNP:99;BLOCK:b5", additional},
		},
	}

	varKey := Variants{header, varData}

	path := filepath.Join(cwd(t), "_testdata", "toy.vcf")
	vars, err := ReadVariants(path)
	if err != nil {
		t.Errorf("ReadVariants(%s) failed with error %s", path, err)
	}

	if !reflect.DeepEqual(vars, varKey) {
		t.Errorf("ReadVarants(%s) loaded a variant object from a file that did not match what was expected.", path)
		fmt.Println(vars)
		fmt.Println(varKey)
	}

}

func TestSimPhasedBlocks(t *testing.T) {

	input := filepath.Join(cwd(t), "_testdata", "toy.noblocks.vcf")
	output := filepath.Join(cwd(t), "_testdata", "test.vcf")

	e := SimPhasedBlocks(input, output, 16, 0)
	if e != nil {
		t.Fatal(e)
	}

	vars, e2 := ReadVariants(output)
	if e2 != nil {
		t.Error(e2)
	}
	varsKey, e3 := ReadVariants(filepath.Join(cwd(t), "_testdata", "toy.vcf"))
	if e3 != nil {
		t.Error(e3)
	}

	if !reflect.DeepEqual(vars, varsKey) {
		t.Errorf("SimPhasedBlocks(%s, %s, 20, 10) yielded a simulated haplotype block structure that didn't match the key", input, output)
		fmt.Println(vars)
		fmt.Println(varsKey)
	}

}


func TestPartitionVariantsByContig(t *testing.T){

	inputPath := filepath.Join(cwd(t), "_testdata", "toy.vcf")
	outstem := filepath.Join(cwd(t), "_testdata", "partitioned")
	PartitionVariantsByContig(inputPath, outstem)

	// Verify partitioning by partitioning the toy files by hand and comparing each of the resulting
	// subset files.

}



