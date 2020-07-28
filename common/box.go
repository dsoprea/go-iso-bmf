package bmfcommon

import (
	"fmt"

	"encoding/binary"

	"github.com/dsoprea/go-logging"
)

var (
	boxLogger = log.NewLogger("bmfcommon.box")
)

var (
	// DefaultEndianness is the default endianness of stored integers.
	DefaultEndianness binary.ByteOrder = binary.BigEndian
)

// Box defines an Atom Box structure.
type Box struct {
	name       string
	start      int64
	size       int64
	headerSize int64
	resource   *BmfResource

	parent CommonBox
}

func NewBox(name string, start, size, headerSize int64, resource *BmfResource) Box {
	return Box{
		name:       name,
		start:      start,
		size:       size,
		headerSize: headerSize,
		resource:   resource,
	}
}

// InlineString returns an undecorated string of field names and values.
func (box Box) InlineString() string {
	parentName := GetParentBoxName(box.parent)

	return fmt.Sprintf(
		"NAME=[%s] PARENT=[%s] START=(0x%016x) SIZE=(%d)",
		box.name, parentName, box.start, box.size)
}

// Name returns the box name.
func (box Box) Name() string {
	return box.name
}

// Start returns the box start offset.
func (box Box) Start() int64 {
	return box.start
}

// Size returns the box size.
func (box Box) Size() int64 {
	return box.size
}

// HeaderSize is the effective size of the header.
func (box Box) HeaderSize() int64 {
	return box.headerSize
}

// Parent returns the parent box.
func (box Box) Parent() CommonBox {
	return box.parent
}

// Index returns the FullBoxIndex for the resource. It contains all previously-
// loaded boxes.
func (box Box) Index() FullBoxIndex {
	return box.resource.Index()
}

// ReadBytesAt reads a box at n and offset.
func (box Box) ReadBytesAt(offset int64, n int64) (b []byte, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	b, err = box.resource.readBytesAt(offset, n)
	log.PanicIf(err)

	return b, nil
}

// ReadBoxes bridges to the lower-level function that knows how to extract child-
// boxes. This also asserts that all box names look valid.
func (box Box) ReadBoxes(startDisplace int, parent CommonBox) (boxes Boxes, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	start := box.Start() + box.HeaderSize() + int64(startDisplace)
	size := box.Size() - box.HeaderSize() - int64(startDisplace)

	boxes, err = readBoxes(box.resource, parent, start, size)
	log.PanicIf(err)

	return boxes, err
}

// ReadBoxData reads the box data from an bmfcommon box.
func (box Box) ReadBoxData() (data []byte, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	headerSize := box.HeaderSize()

	if box.size < headerSize {
		log.Panicf(
			"box [%s] total-size (%d) is smaller then box header-size (%d)",
			box.Name(), box.size, headerSize)
	}

	if box.size == headerSize {
		return nil, nil
	}

	data, err = box.resource.readBytesAt(
		box.start+headerSize,
		box.size-headerSize)

	log.PanicIf(err)

	return data, nil
}

// Boxes is a slice of boxes that have been parsed and are ready to be acted on.
type Boxes []CommonBox

// Index returns a dictionary of boxes, keyed by name.
func (boxes Boxes) Index() (index LoadedBoxIndex) {
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
