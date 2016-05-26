package util

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	//"github.com/gonum/matrix/mat64"
)

// Links stores a set of entity links, such as contigs or sequence
// variants.
type Links struct {

	// A map storing the mapping between entity names and entity ids
	idKey map[string]int

	// A map storing the mapping of entity ids to entity names
	idKeyRev map[int]string

	// The maximum id value
	maxid int

	// The entity ID / entity ID mapping, storing a real number association
	// value between the two.
	data map[int]map[int]float64
}

// NewLinks instantiates a new entity link object
func NewLinks() Links {
	l := Links{}

	// Initialize a map to store string ID to integer ID and int ID to string ID
	// mappings
	l.idKey = make(map[string]int)
	l.idKeyRev = make(map[int]string)

	// Set the current max ID to zero which will be used when encoding string ids
	// to integers, incrementing with each new encoded ID.
	l.maxid = 0

	// Instantiate an empty integer id, integer id, association value map to store
	// the link data.
	l.data = make(map[int]map[int]float64)

	return l
}

// addKey takes an entity name and assigns it an ID
func (l *Links) addKey(key string) {

	// Use the current max ID as the ID for this new item, recording both the
	// integer id to string id and string id to integer id mapping.
	maxid := (*l).maxid
	(*l).idKey[key] = maxid
	(*l).idKeyRev[maxid] = key

	// Increment the max id
	(*l).maxid = maxid + 1

}

// IntIDs returns an array of the integer IDs for the entities in a
// Links object.
func (l *Links) IntIDs() []int {

	// Allocate a slice to store the returned integer entity ids
	ret := make([]int, l.Size())

	// For each entity id in the entity id map, add it to the
	// entity id slice.
	i := 0
	for k, _ := range l.idKeyRev {
		ret[i] = k
		i += 1
	}
	return ret
}

// StringIDs returns an array of string IDs for the entities in a Links
// object.
func (l *Links) StringIDs() []string {

	// Allocate a slice to store the returned string entity names
	ret := make([]string, l.Size())

	// For each entity name in the entity name map, add it to the
	// entity name slice.
	i := 0
	for k, _ := range l.idKey {
		ret[i] = k
		i += 1
	}
	return ret
}

// ID returns the id for a specified entity name, creating a new one if
// the entity is not yet tracked in the Links object.
func (l *Links) ID(key string) int {

	// If the entity name is not in the entity name map, add it
	if _, ok := l.idKey[key]; !ok {
		l.addKey(key)
	}

	// Return the entity id for the specified entity name
	return l.idKey[key]

}

// Set sets the association value for a pair of entity integer ID's.
func (l *Links) Set(id1, id2 int, val float64) error {

	// TODO: Consider allowing this to operate from string instead of
	// integer ids to make code using it more readable.
	// TODO: Consider whether an error should be returned when setting for
	// two IDs not present in the Link set.

	if _, ok := l.idKeyRev[id1]; !ok {
		return fmt.Errorf("sequtil/links: cannot set value for integer id not known to Links object, %d", id1)
	}

	if _, ok := l.idKeyRev[id2]; !ok {
		return fmt.Errorf("sequtil/links: cannot set value for integer id not known to Links object, %d", id2)
	}

	// To prevent duplication in the map, always store the largest
	// of two ids as the top level key.
	if id1 > id2 {
		id2, id1 = id1, id2
	}

	// If the top level key does not exist in the map, allocate a
	// submap for it.
	if _, ok := (*l).data[id1]; !ok {
		(*l).data[id1] = make(map[int]float64)
	}

	// Set the value
	(*l).data[id1][id2] = val

	return nil

}

// Add adds a specified float value to the association value for a pair
// of entity ids.
func (l *Links) Add(id1, id2 int, val float64) error {

	// If either of the specified integer IDs does not exist in the Links object,
	// return an error to that effect.
	if _, ok := l.idKeyRev[id1]; !ok {
		return fmt.Errorf("sequtil/links: cannot add to value for integer id not known to Links object, %d", id1)
	}
	if _, ok := l.idKeyRev[id2]; !ok {
		return fmt.Errorf("sequtil/links: cannot add to value for integer id not known to Links object, %d", id2)
	}

	// To avoid duplication in the map, use the largest of the two id
	// values as the top-level key.
	if id1 > id2 {
		id2, id1 = id1, id2
	}

	// If the top-level key does not exist in the map, allocate a
	// submap for it.
	if _, ok := (*l).data[id1]; !ok {
		(*l).data[id1] = make(map[int]float64)
	}

	// Add the additional association value to the current association value.
	(*l).data[id1][id2] += val

	return nil

}

// AddMulti takes a partial set of variant or block links and adds them to the referenced links object.
func (l *Links) AddMulti() {

}

// Get returns the association value for a pair of entity ids from a
// Links object.
//
// Get will return an error if either of the specified entity ids are not
// present in the Links object. If they are both present but there is no
// association value entry, a zero value is returned.
func (l *Links) Get(id1, id2 int) (float64, error) {

	// If the first specified ID does not exist in the links object, return an error to this effect.
	if _, ok := l.idKeyRev[id1]; !ok {
		return -1, fmt.Errorf("sequtil/links: cannot get value for integer id not known to Links object, %d", id1)
	}

	// If the second specified ID does not exist in the links object, return an error to this effect.
	if _, ok := l.idKeyRev[id2]; !ok {
		return -1, fmt.Errorf("sequtil/links: cannot get value for integer id not known to Links object, %d", id2)
	}

	// Given that the map only contains one copy of each pair association
	// value and that the larger of the two is used as the key for the top
	// level of the map, set the tw
	if id1 > id2 {
		id2, id1 = id1, id2
	}

	// If either no links exist in the links object for the first id or if no links
	// exist between the two ids, return a value of zero.
	if _, ok := (*l).data[id1]; !ok {
		return 0, nil
	} else if _, ok := (*l).data[id1][id2]; !ok {
		return 0, nil
	}

	// Here we know the two ids exist in the map so index into it using these and obtain
	// the value to return.
	return (*l).data[id1][id2], nil

}

// Print prints the full set of entity links in a Links object as
// intID1, intID2, value triplets.
func (l *Links) Print() {

	// For each key in the top level map of the links data map
	for k, _ := range (*l).data {

		// For each key in the second level of the links data map
		for k2, v := range (*l).data[k] {

			// Print the id1, id2, value triplet
			fmt.Println(k, k2, v)

		}
	}
}

// Write takes an writable os.File pointer and writes the stringID1, stringID2, value
// triplets for each association value inthe referenced Links object.
func (l *Links) Write(out *os.File) {

	// Compile a header line which will store the mapping between string keys and iteger
	// keys, with key value separated by ':' and pairs separated by commas
	// i.e. s:i,s:i,s:i,...
	header := "#"
	for k, v := range (*l).idKey {
		header = header + " " + k + ":" + strconv.Itoa(v)
	}
	header += "\n"
	// Write the header string to the output file.
	out.WriteString(header)

	// For each key in the top level Links data map
	for k, _ := range (*l).data {
		// For each key in the second level Links data map
		for k2, v := range (*l).data[k] {
			// Write the string id, string id, float64 triplet to the output file,
			// delimited by a space.
			out.WriteString(fmt.Sprintf("%s %s %f\n", l.idKeyRev[k], l.idKeyRev[k2], v))
		}
	}
}

// Size returns the size of the Links object.
func (l *Links) Size() int {
	return len(l.idKey)
}

// Decode takes an array of entity integer ids and returns an array of entity
// string names.
//
// If the provided entity id array contains an id not tracked in the Links object,
// a non-nil error will be returned.
func (l *Links) Decode(in []int) ([]string, error) {
	out := make([]string, len(in))
	for i, j := range in {
		val, ok := (*l).idKeyRev[j]
		if !ok {
			return []string{}, fmt.Errorf("Error decoding links, an input id was not recognized.")
		}
		out[i] = val
	}
	return out, nil
}

func (l *Links) DecodePhasing(in []bool) (map[string]bool, error) {
	out := make(map[string]bool, len(in))
	for i, j := range in {
		val, ok := (*l).idKeyRev[i]
		if !ok {
			return map[string]bool{}, fmt.Errorf("Error decoding links, an input id was not recognized.")
		}
		out[val] = j
	}
	return out, nil

}

// Subset deletes all elements in a Links object besides those with string ids containing
// a specified prefix or "tag".
func (l *Links) Subset(tag string) {

	// For each element in the Links object
	for s, i := range (*l).idKey {

		// if first part of s doesnt match the tag
		tagplus := tag + "_"
		if s[:len(tagplus)] != tagplus {
			// delete the corresponding entry from both the integer id to string and string
			// to integer id map.
			delete((*l).idKeyRev, i)
			delete((*l).idKey, s)
		}

	}

}

// TabulateVariantLinks tabulates the in-phase (+1) or out-of-phase (-1) links between
// variants in a pair of reads. The function takes the name of a chromosome and a pair of
// maps which store whether a ref or alt variant was observed in each position in each read
// where there is known to be a variant.
func (l *Links) TabulateVariantLinks(chrom string, fgp1, fgp2 map[int]string) (int, int) {

	ct := 0
	balance := 0

	// For each position key in the first position / ref-alt map
	for k1, _ := range fgp1 {

		// For each position key in the second position / ref-alt map
		for k2, _ := range fgp2 {

			// Consider only cases where the first key is greater than the second,
			// avoiding counting self-associations and double counting.
			if k1 > k2 {

				// Construct the string ID from the chromosome number and the key value
				id1 := fmt.Sprintf("%s_%d", chrom, k1)
				id2 := fmt.Sprintf("%s_%d", chrom, k2)

				// If neither ref-alt value is "N"
				if (fgp1[k1] != "N") && (fgp2[k2] != "N") {

					// If the two positions are in phase, i.e. both ref or both alt,
					// add 1 to the association count between the two.
					if fgp1[k1] == fgp2[k2] {
						(*l).Add((*l).ID(id1), (*l).ID(id2), 1) // in phase
						balance += 1
						ct += 1
					} else {
						// Otherwise, the two are out of phase, i.e. one is ref and
						// the other is alt, therefore subtract one from the association
						// count between the two.
						(*l).Add((*l).ID(id1), (*l).ID(id2), -1) // out of phase
						balance -= 1
						ct += 1
					}
				}
			}
		}
	}

	return ct, balance

}

// VariantBlocksForPositions takes an array of genome positions and a connection to a variants
// database and returns the variant blocks (and their phase) to which the positions correpond.
func (l *Links) TabulateVariantBlockLinks(chrom string, fgp1, fgp2 map[int][]int) (int, int) {

	ct := 0
	balance := 0

	// For each position key in the first position / ref-alt map
	for k1, _ := range fgp1 {

		// For each position key in the second position / ref-alt map
		for k2, _ := range fgp2 {

			// Consider only cases where the first key is greater than the second,
			// avoiding counting self-associations and double counting.
			if k1 > k2 {

				// Construct the string ID from the chromosome number and the key value
				id1 := fmt.Sprintf("%s_%d", chrom, k1)
				id2 := fmt.Sprintf("%s_%d", chrom, k2)

				// If neither ref-alt value is "N"
				if (fgp1[k1][2] == 0) && (fgp2[k2][2] == 0) {

					inPhaseCaseOne := fgp1[k1][0] > fgp1[k1][1] && fgp2[k2][0] > fgp2[k2][1]
					inPhaseCaseTwo := fgp1[k1][0] < fgp1[k1][1] && fgp2[k2][0] < fgp2[k2][1]

					//fmt.Println(chrom, k1, k2, fgp1[k1], fgp2[k2])

					// If the two positions are in phase, i.e. both ref or both alt,
					// add 1 to the association count between the two.
					if inPhaseCaseOne || inPhaseCaseTwo {
						(*l).Add((*l).ID(id1), (*l).ID(id2), 1) // in phase
						balance += 1
						ct += 1
					} else {
						// Otherwise, the two are out of phase, i.e. one is ref and
						// the other is alt, therefore subtract one from the association
						// count between the two.
						(*l).Add((*l).ID(id1), (*l).ID(id2), -1) // out of phase
						balance -= 1
						ct += 1
					}
				}
			}
		}
	}

	return ct, balance

}

// LoadLinks loads a set of links from a specified path into a Links object
//
// An error will be returned if no file exists at the specified path or if
// after reading from this file no links were tabulated (meaning either the file
// was empty or contained only header lines).
func LoadLinks(linksPath string) (Links, error) {

	// Instantiate a new links object
	links := NewLinks()

	// Open the input links file for reading, if possible
	lf, err := os.Open(linksPath)
	defer lf.Close()
	if err != nil {
		return Links{}, fmt.Errorf("Couldn't open input file with path %s\n", linksPath)
	}

	// Instantiate a new file scanner
	s := bufio.NewScanner(lf)

	ct := 0

	// For each line in the file
	for s.Scan() {

		// Get the line
		line := s.Text()
		ct += 1
		if string(line[0]) != "#" {
			// If the line is not a header line

			// Partition the line into space-delimited tokens
			arr := strings.Split(line, " ")

			// If the first string ID is not present in the Links object ID set,
			// add it and assign it an integer ID
			if _, ok := links.idKey[arr[0]]; !ok {
				links.addKey(arr[0])
			}

			// If the second string ID is not present in the Links object ID set,
			// add it and assign it an integer ID
			if _, ok := links.idKey[arr[1]]; !ok {
				links.addKey(arr[1])
			}

			// Obtain the integer ID's for each of the string string IDs
			id1 := links.idKey[arr[0]]
			id2 := links.idKey[arr[1]]

			// Parse the string association value into a float64
			val, _ := strconv.ParseFloat(arr[2], 64)

			// Set the association value between the two integer IDs to be the
			// parsed float64
			links.Set(id1, id2, val)

		}

	}

	// If we reach the end of the links file and no links have been tabulated, the file was either
	// empty or contained only a header, return the empty links object and an error to this effect.
	if links.Size() <= 0 {
		return links, fmt.Errorf("read %d non-header lines from a links file and resulted in an empty links object.", ct)
	}

	return links, nil

}
