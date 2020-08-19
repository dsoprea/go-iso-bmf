package bmfcommon

import (
	"io"

	"encoding/binary"

	"github.com/dsoprea/go-logging"
)

var (
	resourceLogger = log.NewLogger("bmfcommon.resource")
)

// Resource defines a file structure.
type Resource struct {
	rs           io.ReadSeeker
	isFragmented bool

	// fullBoxIndex has all [known] boxes encountered in the stream.
	fullBoxIndex FullBoxIndex

	// LoadedBoxIndex contains this box's children.
	LoadedBoxIndex
}

// NewResource returns a new Resource struct.
func NewResource(rs io.ReadSeeker, size int64) (resource *Resource, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	// This has all [known] boxes encountered in the stream.
	fullBoxIndex := make(FullBoxIndex)

	resource = &Resource{
		rs:           rs,
		fullBoxIndex: fullBoxIndex,
	}

	resourceLogger.Debugf(nil, "Parsing stream with (%d) bytes.", size)

	boxes, err := readBoxes(resource, nil, int64(0), size)
	log.PanicIf(err)

	resourceLogger.Debugf(nil, "(%d) root boxes were found.", len(boxes))

	// This has the root boxes from the stream.
	resource.LoadedBoxIndex = boxes.Index()

	return resource, nil
}

// Index returns the complete index of the boxes found in the parsed file.
func (f *Resource) Index() FullBoxIndex {
	return f.fullBoxIndex
}

// readBytesAt reads a box at n and offset.
func (f *Resource) readBytesAt(offset int64, n int64) (b []byte, err error) {
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

// copyBytesAt seeksand then copies N bytes from the resource to the writer.
func (f *Resource) copyBytesAt(offset int64, n int64, w io.Writer) (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	_, err = f.rs.Seek(offset, io.SeekStart)
	log.PanicIf(err)

	_, err = io.CopyN(w, f.rs, n)
	log.PanicIf(err)

	return nil
}

// readBaseBox reads a box from an offset.
func (f *Resource) readBaseBox(offset int64) (box Box, err error) {
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
		// TODO(dustin): Come back to this when we have a supporting example.
		log.Panicf("box [%s] size is (0) and unhandled", boxType)
	}

	box = NewBox(boxType, offset, boxSize, headerSize, f)

	return box, nil
}

// ReadBaseBox reads the base box at the given offset. Supports testing.
func (f *Resource) ReadBaseBox(offset int64) (box Box, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	box, err = f.readBaseBox(offset)
	log.PanicIf(err)

	return box, nil
}

func readBox(f *Resource, parent CommonBox, offset int64) (cb CommonBox, known bool, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

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
		cbis.SetLoadedBoxIndex(boxes)
	}

	return cb, true, nil
}

func readBoxes(f *Resource, parent CommonBox, start int64, totalSize int64) (boxes Boxes, err error) {
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
		} else {
			// We insert nil entries to maintain the integrity of the child
			// list.
			boxes = append(boxes, nil)
		}

		name := cb.Name()
		size := cb.Size()

		resourceLogger.Debugf(nil, "[%s] Child (%d) box has name [%s] and size (%d).", parentName, i, name, size)

		offset += int64(size)
		i++
	}

	return boxes, nil
}
