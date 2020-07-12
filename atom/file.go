package atom

import (
	"io"

	"encoding/binary"

	"github.com/dsoprea/go-logging"
)

// File defines a file structure.
type File struct {
	rs io.ReadSeeker

	ftyp *FtypBox
	moov *MoovBox
	mdat *MdatBox
	size int64

	isFragmented bool
}

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

	for _, box := range boxes {
		switch box.name {
		case "ftyp":
			f.ftyp = &FtypBox{Box: box}
			f.ftyp.parse()

		case "wide":
			// fmt.Println("found wide")

		case "mdat":
			f.mdat = &MdatBox{Box: box}
			// No mdat boxes to parse

		case "moov":
			f.moov = &MoovBox{Box: box}
			f.moov.parse()

			f.isFragmented = f.moov.IsFragmented
		}
	}
	return nil
}

// readBoxAt reads a box from an offset.
func (f *File) readBoxAt(offset int64) (boxSize uint32, boxType string, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	buf, err := f.readBytesAt(BoxHeaderSize, offset)
	log.PanicIf(err)

	boxSize = binary.BigEndian.Uint32(buf[0:4])
	boxType = string(buf[4:8])

	return boxSize, boxType, nil
}

// readBytesAt reads a box at n and offset.
func (f *File) readBytesAt(n int64, offset int64) (word []byte, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	buf := make([]byte, n)

	_, err = f.rs.Seek(offset, io.SeekStart)
	log.PanicIf(err)

	_, err = f.rs.Read(buf)
	log.PanicIf(err)

	return buf, nil
}
