package util

import (
	//"github.com/codegangsta/cli"
	"fmt"
    "os"
    "bufio"
    "strings"
    "strconv"
)


func subsetTagHeader(line string, count, windowsize int) string {

    arr := strings.Split(line, " ")
    arr[0] = arr[0] + "_" + strconv.Itoa(count) + " windowsize:" + strconv.Itoa(windowsize) 
    return strings.Join(arr, " ")

}

func Partition(inPath, outPath string, windowSize int) {

    // Open the ouput file for writing
    out, err := os.Create(outPath)
    if err != nil {
        fmt.Errorf("Couldn't open output file with path %s\n", outPath)
    }
    defer out.Close()

    // Open the input file
    in, err := os.Open(inPath)
    if err != nil {
        fmt.Errorf("Couldn't open input file with path %s\n", inPath)
    }
    defer in.Close()

    cache := []string{}
    s := bufio.NewScanner(in)
    blockCount := 0
    header := ""
    cacheLen := 0
    for s.Scan() {
        line := s.Text()
        if string(line[0]) == ">" {
            if len(cache) != 0 {
                out.WriteString(strings.Join(cache, "") + "\n")
                out.WriteString(subsetTagHeader(line, blockCount, windowSize) + "\n")
            }
            cache = []string{}
            blockCount = 0
            header = line
            cacheLen = 0
        } else {
            cache = append(cache, line)
            cacheLen += len(line)
            if cacheLen > windowSize {
                out.WriteString(subsetTagHeader(header, blockCount, windowSize)+ "\n")
                cacheLine := strings.Join(cache, "")
                cache = []string{}
                out.WriteString(cacheLine[:windowSize] + "\n")
                cache = append(cache, cacheLine[windowSize:])                
                cacheLen = len(cache[0])
                blockCount += 1
            }
        }
    }
}

func Mask(inPath, outPath string, vcf Variants) {

    // Open the ouput file for writing
    out, err1 := os.Create(outPath)
    if err1 != nil {
        fmt.Errorf("Couldn't open output file with path %s\n", outPath)
    }
    defer out.Close()

    // Open the input file
    in, err2 := os.Open(inPath)
    if err2 != nil {
        fmt.Errorf("Couldn't open input file with path %s\n", inPath)
    }
    defer in.Close()
   
    seqOffset := 0
    currentID := ""
    s := bufio.NewScanner(in)
    toggle := true
    for s.Scan() {
        line := s.Text()
        if string(line[0]) == ">" {
            fmt.Printf("masking contig... %s\n", line[1:])
            out.WriteString(line + "\n")
            currentID = strings.Split(line[1:], " ")[0]
            _, toggle = vcf.data[currentID]
            seqOffset = 0
        } else {
            if toggle {
                for i, _ := range line {
                    if _, ok := vcf.data[currentID][seqOffset + i]; ok {
                        line = line[:i] + "N" + line[i+1:]
                    }
                }                
            }
            out.WriteString(line + "\n")
            seqOffset += len(line)
        }
    }
}

