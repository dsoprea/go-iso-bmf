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

	index map[string]*Box
}

// TODO(dustin): This are bridging accessors to keep the tests working while we transition to using indexing.
//>>>>

func (f *File) Ftyp() *FtypBox {
	return f.ftyp
}

func (f *File) Moov() *MoovBox {
	return f.moov
}

func (f *File) Mdat() *MdatBox {
	return f.mdat
}

// GetChildBox returns the given child box or panics uncontrollably.
func (f *File) GetChildBox(name string) CommonBox {
	cb, found := f.index[name]
	if found == false {
		log.Panicf("child box not found: [%s]", name)
	}

	return cb
}

// <<<<

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

			err := f.ftyp.parse()
			log.PanicIf(err)

		case "wide":
			// fmt.Println("found wide")

		case "mdat":
			f.mdat = &MdatBox{Box: box}
			// No mdat boxes to parse

		case "moov":
			f.moov = &MoovBox{Box: box}

			err := f.moov.parse()
			log.PanicIf(err)

			f.isFragmented = f.moov.IsFragmented
		}
	}

	f.index = boxes.Index()

	return nil
}

// readBoxAt reads a box from an offset.
func (f *File) readBoxAt(offset int64) (boxSize uint32, boxType string, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	buf, err := f.readBytesAt(boxHeaderSize, offset)
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
