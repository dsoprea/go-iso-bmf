package bmfcommon

import (
	"io"

	"github.com/dsoprea/go-logging"
)

var (
	fileLogger = log.NewLogger("mp4.bmfcommon.file")
)

// File defines a file structure.
type File struct {
	rs io.ReadSeeker

	size int64

	isFragmented bool

	LoadedBoxIndex
}

// NewFile returns a new File struct.
func NewFile(rs io.ReadSeeker, size int64) *File {
	return &File{
		rs:   rs,
		size: size,
	}
}

// Parse reads an MP4 file for bmfcommon boxes.
func (f *File) Parse() (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	boxes, err := readBoxes(f, int64(0), f.size)
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

func readBox(f *File, offset int64) (cb CommonBox, known bool, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	// TODO(dustin): Add test

	fileLogger.Debugf(nil, "Reading box at offset (0x%016x).", offset)

	box, err := f.ReadBaseBox(offset)
	log.PanicIf(err)

	name := box.Name()

	bf := GetFactory(name)

	if bf == nil {
		boxLogger.Warningf(nil, "No factory registered for box-type [%s].", name)
		return box, false, nil
	}

	// Construct the type-specific box.

	cb, err = bf.New(box)
	log.PanicIf(err)

	return cb, true, nil
}

func readBoxes(f *File, start int64, totalSize int64) (boxes Boxes, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	i := 0
	for offset := start; offset < start+totalSize; {
		fileLogger.Debugf(nil, "Reading box (%d) at offset (0x%016x).", i, offset)

		cb, known, err := readBox(f, offset)
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
