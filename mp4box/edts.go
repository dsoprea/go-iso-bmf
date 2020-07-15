package mp4box

import (
	"github.com/dsoprea/go-logging"
)

// EdtsBox is an "Edit" box.
type EdtsBox struct {
	Box

	LoadedBoxIndex
}

func (b *EdtsBox) parse() (err error) {
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

type edtsBoxFactory struct {
}

// Name returns the name of the type.
func (edtsBoxFactory) Name() string {
	return "edts"
}

// New returns a new value instance.
func (edtsBoxFactory) New(box Box) (cb CommonBox, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	edtsBox := &EdtsBox{
		Box: box,
	}

	err = edtsBox.parse()
	log.PanicIf(err)

	return edtsBox, nil
}

var (
	_ boxFactory = edtsBoxFactory{}
	_ CommonBox  = &EdtsBox{}
)

func init() {
	RegisterBoxType(edtsBoxFactory{})
}
