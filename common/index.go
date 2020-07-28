package bmfcommon

import (
	"fmt"
	"sort"
	"strings"
)

// FullyQualifiedBoxName is the name of a box fully-qualified with its parents.
type FullyQualifiedBoxName []string

// String returns the name in dotted notation.
func (fqbn FullyQualifiedBoxName) String() string {
	return strings.Join(fqbn, ".")
}

// IndexedBoxEntry represents a unique entry in a FullBoxIndex.
type IndexedBoxEntry struct {
	// Name is a box-name fully-qualified with dots.
	NamePhrase string

	// SequenceNumber is the instance number (starting at zero). Unique entries
	// will only be zero but sequence will have identical names with ascending
	// sequence numbers.
	SequenceNumber int
}

// String returns a simply, stringified name for the index entry (dotted
// notation with numeric sequence number appended).
func (ibe IndexedBoxEntry) String() string {

	// TODO(dustin): Add test

	return fmt.Sprintf("%s(%d)", ibe.NamePhrase, ibe.SequenceNumber)
}

// fullBoxIndex describes all boxes encountered (immediately loaded, and loaded
// in the order encountered such that one box's parsing logic will be able to
// reference earlier siblings).
type FullBoxIndex map[IndexedBoxEntry]CommonBox

func (fbi FullBoxIndex) getBoxName(cb CommonBox) (fqbn FullyQualifiedBoxName) {

	// TODO(dustin): Add test

	fqbn = make(FullyQualifiedBoxName, 0)
	for current := cb; current != nil; current = current.Parent() {
		fqbn = append(FullyQualifiedBoxName{current.Name()}, fqbn...)
	}

	return fqbn
}

// Add adds one CommonBox to the index.
func (fbi FullBoxIndex) Add(cb CommonBox) {

	// TODO(dustin): Add test

	name := fbi.getBoxName(cb)

	for i := 0; ; i++ {
		ibe := IndexedBoxEntry{
			NamePhrase:     name.String(),
			SequenceNumber: i,
		}

		if _, found := fbi[ibe]; found == false {
			fbi[ibe] = cb
			break
		}
	}
}

// Dump prints the contents of the full box index.
func (fbi FullBoxIndex) Dump() {

	// TODO(dustin): Add test

	namePhrases := make([]string, len(fbi))
	flatIndex := make(map[string]IndexedBoxEntry)
	i := 0
	for ibe, _ := range fbi {
		namePhrase := ibe.String()
		flatIndex[namePhrase] = ibe
		namePhrases[i] = namePhrase
		i++
	}

	sort.Strings(namePhrases)

	for _, namePhrase := range namePhrases {
		ibe := flatIndex[namePhrase]
		cb := fbi[ibe]
		fmt.Printf("%s: [%s] %s\n", namePhrase, cb.Name(), cb.InlineString())
	}
}
