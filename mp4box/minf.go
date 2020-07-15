package mp4box

import (
	"github.com/dsoprea/go-logging"
)

// MinfBox is a "Media Information" box.
//
// This box contains all the objects that declare characteristics information of
// the media in the track.
type MinfBox struct {
	Box

	LoadedBoxIndex
}

func (b *MinfBox) parse() (err error) {
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

type minfBoxFactory struct {
}

// Name returns the name of the type.
func (minfBoxFactory) Name() string {
	return "minf"
}

// New returns a new value instance.
func (minfBoxFactory) New(box Box) (cb CommonBox, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	minfBox := &MinfBox{
		Box: box,
	}

	err = minfBox.parse()
	log.PanicIf(err)

	return minfBox, nil
}

var (
	_ boxFactory = minfBoxFactory{}
	_ CommonBox  = &MinfBox{}
)

func init() {
	registerAtom(minfBoxFactory{})
}
