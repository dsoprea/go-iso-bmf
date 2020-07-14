package atom

import (
	"io"

	"github.com/dsoprea/go-logging"
)

var (
	fileLogger = log.NewLogger("mp4.atom.file")
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

// Parse reads an MP4 file for atom boxes.
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

	buf, err := f.readBytesAt(offset, boxHeaderSize)
	log.PanicIf(err)

	boxSize = defaultEndianness.Uint32(buf[0:4])
	boxType = string(buf[4:8])

	return boxSize, boxType, nil
}

func readBoxes(f *File, start int64, n int64) (boxes Boxes, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	// TODO(dustin): Can make this a common method?

	i := 0
	for offset := start; offset < start+n; {
		fileLogger.Debugf(nil, "Reading box (%d) at offset (0x%016x).", i, offset)

		size, name, err := f.readBoxAt(offset)
		log.PanicIf(err)

		bf := GetFactory(name)

		if bf != nil {
			// Construct the type-specific box.

			box := newBox(name, offset, int64(size), f)

			c, err := bf.New(box)
			log.PanicIf(err)

			boxes = append(boxes, c)
		} else {
			boxLogger.Warningf(nil, "No factory registered for box-type [%s].", name)
		}

		offset += int64(size)
		i++
	}

	return boxes, nil
}
