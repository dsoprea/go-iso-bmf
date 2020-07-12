package atom

import (
	"io"

	"encoding/binary"

	"github.com/dsoprea/go-logging"
)

const (
	// BoxHeaderSize Size of box header.
	BoxHeaderSize = int64(8)
)

// File defines a file structure.
type File struct {
	rs io.ReadSeeker

	// TODO(dustin): Stop exporting.
	Ftyp *FtypBox
	Moov *MoovBox
	Mdat *MdatBox
	Size int64

	IsFragmented bool
}

func NewFile(rs io.ReadSeeker, size int64) *File {
	return &File{
		rs:   rs,
		Size: size,
	}
}

// Parse reads an MP4 file for atom boxes.
func (f *File) Parse() (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	boxes, err := readBoxes(f, int64(0), f.Size)
	log.PanicIf(err)

	for _, box := range boxes {
		switch box.Name {
		case "ftyp":
			f.Ftyp = &FtypBox{Box: box}
			f.Ftyp.parse()
		case "wide":
			// fmt.Println("found wide")
		case "mdat":
			f.Mdat = &MdatBox{Box: box}
			// No mdat boxes to parse
		case "moov":
			f.Moov = &MoovBox{Box: box}
			f.Moov.parse()

			f.IsFragmented = f.Moov.IsFragmented
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

// Box defines an Atom Box structure.
type Box struct {
	Name        string
	Size, Start int64
	File        *File
}

type CommonBox interface {
	// TODO(dustin): Rename to Data()
	readBoxData() (data []byte, err error)
}

type boxFactory interface {
	New(box *Box) (cb CommonBox, err error)
	Name() string
}

// ReadBoxData reads the box data from an atom box.
func (b *Box) readBoxData() (data []byte, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	if b.Size <= BoxHeaderSize {
		return nil, nil
	}

	data, err = b.File.readBytesAt(b.Size-BoxHeaderSize, b.Start+BoxHeaderSize)
	log.PanicIf(err)

	return data, nil
}

func readBoxes(f *File, start int64, n int64) (l []*Box, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	// TODO(dustin): Make this a common method?

	for offset := start; offset < start+n; {
		size, name, err := f.readBoxAt(offset)
		log.PanicIf(err)

		b := &Box{
			Name:  string(name),
			Size:  int64(size),
			File:  f,
			Start: offset,
		}

		l = append(l, b)
		offset += int64(size)
	}

	return l, nil
}
