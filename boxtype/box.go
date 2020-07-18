package boxtype

import (
	"fmt"

	"encoding/binary"

	"github.com/dsoprea/go-logging"
)

var (
	boxLogger = log.NewLogger("mp4.boxtype.box")
)

const (
	// boxHeaderSize is the size of box header.
	boxHeaderSize = int64(8)
)

var (
	defaultEndianness binary.ByteOrder = binary.BigEndian
)

// Box defines an Atom Box structure.
type Box struct {
	name  string
	start int64
	size  int64
	file  *File
}

func newBox(name string, start, size int64, file *File) Box {

	// TODO(dustin): Add test

	return Box{
		name:  name,
		start: start,
		size:  size,
		file:  file,
	}
}

// InlineString returns an undecorated string of field names and values.
func (box Box) InlineString() string {

	// TODO(dustin): Add test

	return fmt.Sprintf("NAME=[%s] START=(%d) SIZE=(%d)", box.name, box.start, box.size)
}

// Name returns the box name.
func (box Box) Name() string {

	// TODO(dustin): Add test

	return box.name
}

// Size returns the box size.
func (box Box) Size() int64 {

	// TODO(dustin): Add test

	return box.size
}

// Start returns the box start offset.
func (box Box) Start() int64 {

	// TODO(dustin): Add test

	return box.start
}

func (box Box) readBoxes(startDisplace int) (boxes Boxes, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	// TODO(dustin): Add test

	start := box.Start() + boxHeaderSize + int64(startDisplace)
	size := box.Size() - boxHeaderSize - int64(startDisplace)

	boxes, err = readBoxes(box.file, start, size)

	log.PanicIf(err)

	return boxes, err
}

// ReadBoxData reads the box data from an boxtype box.
func (box Box) readBoxData() (data []byte, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	// TODO(dustin): Add test

	if box.size <= boxHeaderSize {
		return nil, nil
	}

	data, err = box.file.readBytesAt(
		box.start+boxHeaderSize,
		box.size-boxHeaderSize)

	log.PanicIf(err)

	return data, nil
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

	// TODO(dustin): Add test

	boxes, found := lbi[name]
	if found == false {
		log.Panicf("child box not found: [%s]", name)
	}

	return boxes, nil
}

// ChildrenTypes returns a slice with the names of all children with registered
// types.
func (lbi LoadedBoxIndex) ChildrenTypes() (names []string) {

	// TODO(dustin): Add test

	names = make([]string, len(lbi))
	i := 0
	for name, _ := range lbi {
		names[i] = name
		i++
	}

	return names
}

// Boxes is a slice of boxes that have been parsed and are ready to be acted on.
type Boxes []CommonBox

// Index returns a dictionary of boxes, keyed by name.
func (boxes Boxes) Index() (index LoadedBoxIndex) {

	// TODO(dustin): Add test

	index = make(LoadedBoxIndex)

	for _, box := range boxes {
		name := box.Name()
		if existing, found := index[name]; found == true {
			index[name] = append(existing, box)
		} else {
			index[name] = []CommonBox{box}
		}
	}

	return index
}
