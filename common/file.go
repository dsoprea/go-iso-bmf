package bmfcommon

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"encoding/binary"

	"github.com/dsoprea/go-logging"
)

var (
	fileLogger = log.NewLogger("mp4.bmfcommon.file")
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

// File defines a file structure.
type File struct {
	rs           io.ReadSeeker
	size         int64
	isFragmented bool

	fullBoxIndex FullBoxIndex

	// LoadedBoxIndex contains this boxes children.
	LoadedBoxIndex
}

// NewFile returns a new File struct.
func NewFile(rs io.ReadSeeker, size int64) *File {
	fullBoxIndex := make(FullBoxIndex)

	return &File{
		rs:           rs,
		size:         size,
		fullBoxIndex: fullBoxIndex,
	}
}

// Index returns the complete index of the boxes found in the parsed file.
func (f *File) Index() FullBoxIndex {

	// TODO(dustin): Add test

	return f.fullBoxIndex
}

// Parse reads an MP4 file for bmfcommon boxes.
func (f *File) Parse() (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	// TODO(dustin): !! Dump Parse() and move this to NewFile. This might break a lot of unit-tests.

	fileLogger.Debugf(nil, "Parsing stream with (%d) bytes.", f.size)

	boxes, err := readBoxes(f, nil, int64(0), f.size)
	log.PanicIf(err)

	fileLogger.Debugf(nil, "(%d) root boxes were found.", len(boxes))

	f.LoadedBoxIndex = boxes.Index()

	return nil
}

// readBytesAt reads a box at n and offset.
func (f *File) readBytesAt(offset int64, n int64) (b []byte, err error) {
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
func (f *File) readBoxAt(offset int64) (box Box, err error) {
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

	var headerSize int64

	if boxSize > 1 {
		// We have an alternative 32-bit box-size.

		headerSize = 8
	} else if boxSize == 1 {
		// We have an alternative 64-bit box-size. It follows the size and
		// type.

		headerSize = 16

		fileLogger.Debugf(nil,
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

func (f *File) ReadBaseBox(offset int64) (box Box, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	// TODO(dustin): Add test

	box, err = f.readBoxAt(offset)
	log.PanicIf(err)

	return box, nil
}

func readBox(f *File, parent CommonBox, offset int64) (cb CommonBox, known bool, err error) {
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
		fileLogger.Warningf(nil, "No factory registered for box-type [%s].", name)
		return box, false, nil
	}

	// Construct the type-specific box.

	cb, err = bf.New(box)
	log.PanicIf(err)

	f.fullBoxIndex.Add(cb)

	return cb, true, nil
}

func readBoxes(f *File, parent CommonBox, start int64, totalSize int64) (boxes Boxes, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	parentName := GetParentBoxName(parent)

	i := 0
	for offset := start; offset < start+totalSize; {
		fileLogger.Debugf(nil, "[%s] Reading child (%d) box at offset (0x%016x).", parentName, i, offset)

		cb, known, err := readBox(f, parent, offset)
		log.PanicIf(err)

		// We'll interpret everything as data. So, if there is good data
		// followed by garbage, we'll interpret the garbage as well. So, if we
		// see a box with an invalid name, skip it and stop reading any further.

		name := cb.Name()
		if BoxNameIsValid(name) == false {
			log.Panicf(
				"box (%d) in sequence starting at offset (0x%016x) looks like garbage",
				i, offset)
		}

		if known == true {
			boxes = append(boxes, cb)
		}

		boxSize := cb.Size()

		fileLogger.Debugf(nil, "[%s] Child (%d) box size is (%d).", parentName, i, boxSize)

		offset += int64(boxSize)
		i++
	}

	return boxes, nil
}
