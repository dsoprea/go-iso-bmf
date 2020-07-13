package atom

import (
	"github.com/dsoprea/go-logging"
)

var (
	boxLogger = log.NewLogger("mp4.atom.box")
)

const (
	// boxHeaderSize is the size of box header.
	boxHeaderSize = int64(8)
)

// Box defines an Atom Box structure.
type Box struct {
	name  string
	size  int64
	start int64
	file  *File
}

// Name returns the box name.
func (box *Box) Name() string {
	return box.name
}

// Size returns the box size.
func (box *Box) Size() int64 {
	return box.size
}

// Start returns the box start offset.
func (box *Box) Start() int64 {
	return box.start
}

func (box *Box) readBoxes(startDisplace int) (boxes Boxes, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	start := box.Start() + boxHeaderSize + int64(startDisplace)
	stop := box.Size() - boxHeaderSize

	boxes, err = readBoxes(box.file, start, stop)

	log.PanicIf(err)

	return boxes, err
}

// ReadBoxData reads the box data from an atom box.
func (b *Box) readBoxData() (data []byte, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	if b.size <= boxHeaderSize {
		return nil, nil
	}

	data, err = b.file.readBytesAt(b.size-boxHeaderSize, b.start+boxHeaderSize)
	log.PanicIf(err)

	return data, nil
}

// LoadedBoxIndex provides a GetChildBoxes() method that returns a child box if
// present or panics with a descriptive error.
type LoadedBoxIndex map[string][]CommonBox

// GetChildBox returns the given child box or panics. If box does not support
// children this should return ErrNoChildren.
func (lbi LoadedBoxIndex) GetChildBoxes(name string) (boxes []CommonBox, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	boxes, found := lbi[name]
	if found == false {
		log.Panicf("child box not found: [%s]", name)
	}

	return boxes, nil
}

// Boxes is a slice of boxes that have been parsed and are ready to be acted on.
type Boxes []CommonBox

// Index returns a dictionary of boxes, keyed by name.
func (boxes Boxes) Index() (index LoadedBoxIndex) {

	// TODO(dustin): !! Can there be duplicates (read: sequences of boxes that may have more than one of the same type)?

	index = make(LoadedBoxIndex)

	for _, box := range boxes {
		name := box.Name()
		if existing, found := index[name]; found == true {
			index[name] = append(existing, box)
		} else {
			index[name] = []CommonBox{box}
		}
	}

	return index
}

// UnsupportedBoxIndex provides a GetChildBoxes() method that always panics due
// to lack of support.
type UnsupportedBoxIndex struct{}

// GetChildBox returns the given child box or panics uncontrollably.
func (UnsupportedBoxIndex) GetChildBoxes(name string) (CommonBox, error) {
	return nil, ErrNoChildren
}

func readBoxes(f *File, start int64, n int64) (boxes Boxes, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	// TODO(dustin): Can make this a common method?

	for offset := start; offset < start+n; {
		size, name, err := f.readBoxAt(offset)
		log.PanicIf(err)

		b := &Box{
			name:  name,
			size:  int64(size),
			file:  f,
			start: offset,
		}

		bf := GetFactory(name)

		if bf != nil {
			// Construct the type-specific box.
			c, err := bf.New(b)
			log.PanicIf(err)

			boxes = append(boxes, c)
		} else {
			boxLogger.Warningf(nil, "No factory registered for box-type [%s].", name)
		}

		offset += int64(size)
	}

	return boxes, nil
}
