package bmfcommon

import (
	"fmt"
	"io"
	"sort"
	"strings"

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

	boxes, err := readBoxes(f, nil, int64(0), f.size)
	log.PanicIf(err)

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
func (f *File) readBoxAt(offset int64) (boxSize uint32, boxType string, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	buf, err := f.readBytesAt(offset, BoxHeaderSize)
	log.PanicIf(err)

	boxSize = DefaultEndianness.Uint32(buf[0:4])
	boxType = string(buf[4:8])

	return boxSize, boxType, nil
}

func (f *File) ReadBaseBox(offset int64) (box Box, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	// TODO(dustin): Add test

	size, name, err := f.readBoxAt(offset)
	log.PanicIf(err)

	box = NewBox(name, offset, int64(size), f)

	return box, nil
}

func readBox(f *File, parent CommonBox, offset int64) (cb CommonBox, known bool, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	// TODO(dustin): Add test

	fileLogger.Debugf(nil, "Reading box at offset (0x%016x).", offset)

	box, err := f.ReadBaseBox(offset)
	log.PanicIf(err)

	box.parent = parent

	name := box.Name()

	bf := GetFactory(name)

	if bf == nil {
		boxLogger.Warningf(nil, "No factory registered for box-type [%s].", name)
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

	i := 0
	for offset := start; offset < start+totalSize; {
		fileLogger.Debugf(nil, "Reading box (%d) at offset (0x%016x).", i, offset)

		cb, known, err := readBox(f, parent, offset)
		log.PanicIf(err)

		if known == true {
			boxes = append(boxes, cb)
		} else {
			name := cb.Name()
			boxLogger.Warningf(nil, "No factory registered for box-type [%s].", name)
		}

		boxSize := cb.Size()

		offset += int64(boxSize)
		i++
	}

	return boxes, nil
}
