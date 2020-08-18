package bmfcommon

import (
	"fmt"
	"sort"
	"strings"

	"github.com/dsoprea/go-logging"
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
	return fmt.Sprintf("%s(%d)", ibe.NamePhrase, ibe.SequenceNumber)
}

// FullBoxIndex describes all boxes encountered (immediately loaded, and loaded
// in the order encountered such that one box's parsing logic will be able to
// reference earlier siblings).
type FullBoxIndex map[IndexedBoxEntry]CommonBox

func (fbi FullBoxIndex) getBoxName(cb CommonBox) (fqbn FullyQualifiedBoxName) {
	fqbn = make(FullyQualifiedBoxName, 0)
	for current := cb; current != nil; current = current.Parent() {
		fqbn = append(FullyQualifiedBoxName{current.Name()}, fqbn...)
	}

	return fqbn
}

// Add adds one CommonBox to the index.
func (fbi FullBoxIndex) Add(cb CommonBox) {
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
	namePhrases := make([]string, len(fbi))
	flatIndex := make(map[string]IndexedBoxEntry)
	i := 0
	for ibe := range fbi {
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

// LoadedBoxIndex provides a GetChildBoxes() method that returns a child box if
// present or panics with a descriptive error.
type LoadedBoxIndex map[string][]CommonBox

// GetChildBoxes returns the given child box or panics. If box does not support
// children this should return ErrNoChildren.
func (lbi LoadedBoxIndex) GetChildBoxes(name string) (boxes []CommonBox, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	boxes, found := lbi[name]
	if found == false {
		log.Panicf("child box not found: [%s]", name)
	}

	return boxes, nil
}

// ChildrenTypes returns a slice with the names of all children with registered
// types.
func (lbi LoadedBoxIndex) ChildrenTypes() (names []string) {
	names = make([]string, len(lbi))
	i := 0
	for name := range lbi {
		names[i] = name
		i++
	}

	sort.Strings(names)

	return names
}

// ChildBoxIndexSetter is a box that is known to support children and will be
// called, after they were parsed, to store them.
type ChildBoxIndexSetter interface {
	// SetLoadedBoxIndex sets the child boxes after a box has been manufactured
	// and the children have been parsed. This allows parent boxes to be
	// registered before the child boxes can look for them.
	SetLoadedBoxIndex(boxes Boxes)
}
