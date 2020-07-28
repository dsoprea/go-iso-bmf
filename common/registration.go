package bmfcommon

import (
	"errors"
	"strings"
	"unicode"

	"github.com/dsoprea/go-logging"
)

var (
	// ErrNoChildren indicates that the given box does not support children.
	ErrNoChildren = errors.New("box does not support children")
)

var (
	boxMapping = make(map[string]BoxFactory)
)

// BoxNameIsValid returns true if no invalid characters are in the box-name.
// This is a strategy to determine if there is garbage at the end of the ISO
// 14496-12 data, since we'll just keep reading boxes until we reach the end of
// the allotment.
func BoxNameIsValid(name string) bool {
	// Trim right-side spacing. Spacing is valid on the right side, and this
	// will simplify things.
	name = strings.TrimRight(name, " ")

	// Name needs to be non-empty.
	if name == "" {
		return false
	}

	// Name needs to have only letters. Note that this will also fail if there
	// were spaces *in the middle* of the name.
	for _, r := range name {
		if unicode.IsLetter(r) == false && unicode.IsDigit(r) == false {
			return false
		}
	}

	return true
}

// GetParentBoxName returns the name of the given CB. If nil, returns "ROOT".
func GetParentBoxName(cb CommonBox) string {
	if cb == nil {
		return "ROOT"
	}

	return cb.Name()
}

// ChildBoxes is a simple wrapper that gets all children of the given type or
// panics if none.
func ChildBoxes(bci BoxChildIndexer, name string) (boxes []CommonBox) {
	boxes, err := bci.GetChildBoxes(name)
	log.PanicIf(err)

	// TODO(dustin): Add test

	return boxes
}

// bmfcommon.CommonBox is one parsed box.
type CommonBox interface {
	// TODO(dustin): Rename to Data()
	// ReadBoxData returns the bytes that were encapsulated in this box.
	ReadBoxData() (data []byte, err error)

	// Name returns the name of the box-type.
	Name() string

	// Size is the total size of the box on disk including standard eight-byte
	// header.
	Size() int64

	// InlineString returns an undecorated string of field names and values.
	InlineString() string

	// Parent is the parent box
	Parent() CommonBox
}

// BoxChildIndexer is a box that has children.
type BoxChildIndexer interface {
	// GetChildBoxes returns all found child boxes of the given type.
	GetChildBoxes(name string) (boxes []CommonBox, err error)

	// ChildrenTypes returns the names of the types of the children that were
	// found. Only registered types are recognized.
	ChildrenTypes() (names []string)
}

// BoxFactory knows how to construct a box struct.
type BoxFactory interface {
	// New reads, parses, loads, and returns the value struct given the common
	// box info.
	New(box Box) (cb CommonBox, childBoxSeriesOffset int, err error)

	// Name returns the name of the box-type that this factory knows how to
	// parse.
	Name() string
}

// RegisterBoxType registers the factory for a box-type.
func RegisterBoxType(bf BoxFactory) {

	// TODO(dustin): Add test

	name := bf.Name()

	if _, found := boxMapping[name]; found == true {
		log.Panicf("box-factory already registered: [%s]", name)
	}

	boxMapping[name] = bf
}

// GetFactory returns the factory for the given box-type. Will return `nil` if
// not known.
func GetFactory(name string) BoxFactory {

	// TODO(dustin): Add test

	return boxMapping[name]
}
