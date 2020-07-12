package atom

import (
	"github.com/dsoprea/go-logging"
)

const (
	// TODO(dustin): Stop exporting BoxHeaderSize

	// BoxHeaderSize Size of box header.
	BoxHeaderSize = int64(8)
)

// Box defines an Atom Box structure.
type Box struct {
	name  string
	size  int64
	start int64
	file  *File
}

func (box Box) Name() string {
	return box.name
}

func (box Box) Size() int64 {
	return box.size
}

func (box Box) Start() int64 {
	return box.start
}

func (box Box) File() *File {
	return box.file
}

func (box Box) readBoxes(startDisplace int) (boxes Boxes, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	boxes, err = readBoxes(box.File(), box.Start()+BoxHeaderSize+int64(startDisplace), box.Size()-BoxHeaderSize)
	log.PanicIf(err)

	return boxes, err
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

	if b.size <= BoxHeaderSize {
		return nil, nil
	}

	data, err = b.file.readBytesAt(b.size-BoxHeaderSize, b.start+BoxHeaderSize)
	log.PanicIf(err)

	return data, nil
}

type Boxes []*Box

func (boxes Boxes) Index() (index map[string]*Box) {
	index = make(map[string]*Box)

	for _, box := range boxes {
		index[box.name] = box
	}

	return index
}

func readBoxes(f *File, start int64, n int64) (boxes Boxes, err error) {
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
			name:  string(name),
			size:  int64(size),
			file:  f,
			start: offset,
		}

		boxes = append(boxes, b)
		offset += int64(size)
	}

	return boxes, nil
}
