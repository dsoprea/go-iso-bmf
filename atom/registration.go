package atom

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
}

type BoxChildIndexer interface {
	GetChildBox(name string) (cb CommonBox, err error)
}

// MustGetChildBox is a simple wrapper that panics if the child box could not be
// gotten.
func MustGetChildBox(cb BoxChildIndexer, name string) (ccb CommonBox) {
	ccb, err := cb.GetChildBox(name)
	log.PanicIf(err)

	return ccb
}

type boxFactory interface {
	// New reads, parses, loads, and returns the value struct given the common
	// box info.
	New(box *Box) (cb CommonBox, err error)

	// Name returns the name of the box-type that this factory knows how to
	// parse.
	Name() string
}

func registerAtom(bf boxFactory) {
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
