package mp4box

import (
	"github.com/dsoprea/go-logging"
)

// MdiaBox is a "Media" box.
//
// The media declaration container contains all the objects that declare information
// about the media data within a track.
type MdiaBox struct {
	Box

	LoadedBoxIndex
}

func (b *MdiaBox) parse() (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	boxes, err := b.Box.readBoxes(0)
	log.PanicIf(err)

	b.LoadedBoxIndex = boxes.Index()

	return nil
}

type mdiaBoxFactory struct {
}

// Name returns the name of the type.
func (mdiaBoxFactory) Name() string {
	return "mdia"
}

// New returns a new value instance.
func (mdiaBoxFactory) New(box Box) (cb CommonBox, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	mdiaBox := &MdiaBox{
		Box: box,
	}

	err = mdiaBox.parse()
	log.PanicIf(err)

	return mdiaBox, nil
}

var (
	_ boxFactory = mdiaBoxFactory{}
	_ CommonBox  = &MdiaBox{}
)

func init() {
	registerAtom(mdiaBoxFactory{})
}
