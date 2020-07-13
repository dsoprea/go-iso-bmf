package atom

import (
	"errors"

	"github.com/dsoprea/go-logging"
)

const (
	// boxHeaderSize is the size of box header.
	boxHeaderSize = int64(8)
)

// Box defines an Atom Box structure.
type Box struct {
	name  string
	size  int64
	start int64
	file  *File

	// UnsupportedBoxIndex will set the base Box implementation as not
	// supporting children by default. We need this so that Box satisfies the
	// CommonBox interface by default.
	UnsupportedBoxIndex
}

// Name returns the box name.
func (box *Box) Name() string {
	return box.name
}

// Size returns the box size.
func (box *Box) Size() int64 {
	return box.size
}

// Start returns the box start offset.
func (box *Box) Start() int64 {
	return box.start
}

func (box Box) readBoxes(startDisplace int) (boxes Boxes, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	start := box.Start() + boxHeaderSize + int64(startDisplace)
	stop := box.Size() - boxHeaderSize

	boxes, err = readBoxes(box.file, start, stop)

	log.PanicIf(err)

	return boxes, err
}

var (
	// ErrNoChildren indicates that the given box does not support children.
	ErrNoChildren = errors.New("box does not support children")
)

// CommonBox is one parsed box.
type CommonBox interface {
	// TODO(dustin): Rename to Data()
	// readBoxData returns the bytes that were encapsulated in this box.
	readBoxData() (data []byte, err error)

	// GetChildBox returns the given child box or panics. If box does not
	// support children this should return ErrNoChildren.
	GetChildBox(name string) (cb CommonBox, err error)
}

// MustGetChildBox is a simple wrapper that panics if the child box could not be
// gotten.
func MustGetChildBox(cb CommonBox, name string) (ccb CommonBox) {
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

// ReadBoxData reads the box data from an atom box.
func (b *Box) readBoxData() (data []byte, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	if b.size <= boxHeaderSize {
		return nil, nil
	}

	data, err = b.file.readBytesAt(b.size-boxHeaderSize, b.start+boxHeaderSize)
	log.PanicIf(err)

	return data, nil
}

// Boxes is a slice of boxes.
type Boxes []*Box

// LoadedBoxIndex provides a GetChildBox() method that returns a child box if
// present or panics with a descriptive error.
type LoadedBoxIndex map[string]*Box

// GetChildBox returns the given child box or panics. If box does not support
// children this should return ErrNoChildren.
func (lbi LoadedBoxIndex) GetChildBox(name string) (cb CommonBox, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	cb, found := lbi[name]
	if found == false {
		log.Panicf("child box not found: [%s]", name)
	}

	return cb, nil
}

// UnsupportedBoxIndex provides a GetChildBox() method that always panics.
type UnsupportedBoxIndex struct{}

// GetChildBox returns the given child box or panics uncontrollably.
func (UnsupportedBoxIndex) GetChildBox(name string) (cb CommonBox, err error) {
	return nil, ErrNoChildren
}

// Index returns a dictionary of boxes, keyed by name.
func (boxes Boxes) Index() (index LoadedBoxIndex) {

	// TODO(dustin): !! Can there be duplicates (read: sequences of boxes that may have more than one of the same type)?

	index = make(LoadedBoxIndex)

	for _, box := range boxes {
		index[box.name] = box
	}

	return index
}

func readBoxes(f *File, start int64, n int64) (boxes Boxes, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	// TODO(dustin): Make this a common method?

	for offset := start; offset < start+n; {
		size, name, err := f.readBoxAt(offset)
		log.PanicIf(err)

		b := &Box{
			name:  string(name),
			size:  int64(size),
			file:  f,
			start: offset,
		}

		boxes = append(boxes, b)
		offset += int64(size)
	}

	return boxes, nil
}
