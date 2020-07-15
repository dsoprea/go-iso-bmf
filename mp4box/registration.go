package mp4box

import (
	"errors"

	"github.com/dsoprea/go-logging"
)

var (
	// ErrNoChildren indicates that the given box does not support children.
	ErrNoChildren = errors.New("box does not support children")
)

var (
	boxMapping = make(map[string]boxFactory)
)

// CommonBox is one parsed box.
type CommonBox interface {
	// TODO(dustin): Rename to Data()
	// readBoxData returns the bytes that were encapsulated in this box.
	readBoxData() (data []byte, err error)

	// Name returns the name of the box-type.
	Name() string

	// Size is the total size of the box on disk including standard eight-byte
	// header.
	Size() int64

	// InlineString returns an undecorated string of field names and values.
	InlineString() string
}

// BoxChildIndexer is a box that has children.
type BoxChildIndexer interface {
	// GetChildBoxes returns all found child boxes of the given type.
	GetChildBoxes(name string) (boxes []CommonBox, err error)

	// ChildrenTypes returns the names of the types of the children that were
	// found. Only registered types are recognized.
	ChildrenTypes() (names []string)
}

// ChildBoxes is a simple wrapper that gets all children of the given type or
// panics if none.
func ChildBoxes(bci BoxChildIndexer, name string) (boxes []CommonBox) {
	boxes, err := bci.GetChildBoxes(name)
	log.PanicIf(err)

	return boxes
}

type boxFactory interface {
	// New reads, parses, loads, and returns the value struct given the common
	// box info.
	New(box Box) (cb CommonBox, err error)

	// Name returns the name of the box-type that this factory knows how to
	// parse.
	Name() string
}

func RegisterBoxType(bf boxFactory) {
	name := bf.Name()

	if _, found := boxMapping[name]; found == true {
		log.Panicf("box-factory already registered: [%s]", name)
	}

	boxMapping[name] = bf
}

// GetFactory returns the factory for the given box-type. Will return `nil` if
// not known.
func GetFactory(name string) boxFactory {
	return boxMapping[name]
}
