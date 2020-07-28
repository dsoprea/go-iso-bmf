package bmfcommon

import (
	"io"

	"encoding/binary"

	"github.com/dsoprea/go-logging"
)

var (
	resourceLogger = log.NewLogger("bmfcommon.resource")
)

// BmfResource defines a file structure.
type BmfResource struct {
	rs           io.ReadSeeker
	isFragmented bool

	fullBoxIndex FullBoxIndex

	// LoadedBoxIndex contains this boxes children.
	LoadedBoxIndex
}

// NewBmfResource returns a new BmfResource struct.
func NewBmfResource(rs io.ReadSeeker, size int64) *BmfResource {
	fullBoxIndex := make(FullBoxIndex)

	resource := &BmfResource{
		rs:           rs,
		fullBoxIndex: fullBoxIndex,
	}

	resourceLogger.Debugf(nil, "Parsing stream with (%d) bytes.", size)

	boxes, err := readBoxes(resource, nil, int64(0), size)
	log.PanicIf(err)

	resourceLogger.Debugf(nil, "(%d) root boxes were found.", len(boxes))

	resource.LoadedBoxIndex = boxes.Index()

	return resource
}

// Index returns the complete index of the boxes found in the parsed file.
func (f *BmfResource) Index() FullBoxIndex {

	// TODO(dustin): Add test

	return f.fullBoxIndex
}

// readBytesAt reads a box at n and offset.
func (f *BmfResource) readBytesAt(offset int64, n int64) (b []byte, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	b = make([]byte, n)

	_, err = f.rs.Seek(offset, io.SeekStart)
	log.PanicIf(err)

	_, err = f.rs.Read(b)
	log.PanicIf(err)

	return b, nil
}

// readBoxAt reads a box from an offset.
func (f *BmfResource) readBoxAt(offset int64) (box Box, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	_, err = f.rs.Seek(offset, io.SeekStart)
	log.PanicIf(err)

	// Read 32-bit box-size.

	var rawBoxSize uint32

	err = binary.Read(f.rs, DefaultEndianness, &rawBoxSize)
	log.PanicIf(err)

	boxSize := int64(rawBoxSize)

	// Read box-type.

	boxTypeRaw := make([]byte, 4)

	_, err = io.ReadFull(f.rs, boxTypeRaw)
	log.PanicIf(err)

	boxType := string(boxTypeRaw)

	// We'll interpret everything as data. So, if there is good data
	// followed by garbage, we may interpret the garbage as well. So, if we
	// see a box with an invalid name, panic as soon as possible.
	if BoxNameIsValid(boxType) == false {
		log.Panicf(
			"box starting at offset (0x%016x) looks like garbage",
			offset)
	}

	var headerSize int64

	if boxSize > 1 {
		// We have an alternative 32-bit box-size.

		headerSize = 8
	} else if boxSize == 1 {
		// We have an alternative 64-bit box-size. It follows the size and
		// type.

		headerSize = 16

		resourceLogger.Debugf(nil,
			"Box [%s] at offset (0x%016x) has a 64-bit size.",
			boxType, offset)

		var rawBoxSize uint64

		err = binary.Read(f.rs, DefaultEndianness, &rawBoxSize)
		log.PanicIf(err)

		if rawBoxSize > 0x7FFFFFFFFFFFFFFF {
			log.Panicf("box-size too large for int64")
		}

		boxSize = int64(rawBoxSize)
	} else if boxSize == 0 {
		// TODO(dustin): Come back to this.
		log.Panicf("box [%s] size is (0) and unhandled", boxType)
	}

	box = NewBox(boxType, offset, boxSize, headerSize, f)

	return box, nil
}

func (f *BmfResource) ReadBaseBox(offset int64) (box Box, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	// TODO(dustin): Add test
	// TODO(dustin): Drop this method

	box, err = f.readBoxAt(offset)
	log.PanicIf(err)

	return box, nil
}

func readBox(f *BmfResource, parent CommonBox, offset int64) (cb CommonBox, known bool, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	// TODO(dustin): Add test

	box, err := f.ReadBaseBox(offset)
	log.PanicIf(err)

	box.parent = parent

	name := box.Name()

	bf := GetFactory(name)

	if bf == nil {
		resourceLogger.Warningf(nil, "No factory registered for box-type [%s].", name)
		return box, false, nil
	}

	// Construct the type-specific box.

	cb, childBoxSeriesOffset, err := bf.New(box)
	log.PanicIf(err)

	f.fullBoxIndex.Add(cb)

	if childBoxSeriesOffset >= 0 {
		boxes, err := box.ReadBoxes(childBoxSeriesOffset, cb)
		log.PanicIf(err)

		cbis := cb.(ChildBoxIndexSetter)

		fbi := boxes.Index()
		cbis.SetLoadedBoxIndex(fbi)
	}

	return cb, true, nil
}

func readBoxes(f *BmfResource, parent CommonBox, start int64, totalSize int64) (boxes Boxes, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	parentName := GetParentBoxName(parent)

	i := 0
	for offset := start; offset < start+totalSize; {
		resourceLogger.Debugf(nil, "[%s] Reading child (%d) box at offset (0x%016x).", parentName, i, offset)

		cb, known, err := readBox(f, parent, offset)
		log.PanicIf(err)

		if known == true {
			boxes = append(boxes, cb)
		}

		name := cb.Name()
		size := cb.Size()

		resourceLogger.Debugf(nil, "[%s] Child (%d) box has name [%s] and size (%d).", parentName, i, name, size)

		offset += int64(size)
		i++
	}

	return boxes, nil
}
