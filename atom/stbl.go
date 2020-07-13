package atom

import (
	"github.com/dsoprea/go-logging"
)

// StblBox is a "Sample Table" box.
type StblBox struct {
	Box

	LoadedBoxIndex
}

func (b *StblBox) parse() (err error) {
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

type stblBoxFactory struct {
}

// Name returns the name of the type.
func (stblBoxFactory) Name() string {
	return "stbl"
}

// New returns a new value instance.
func (stblBoxFactory) New(box Box) (cb CommonBox, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	stblBox := &StblBox{
		Box: box,
	}

	err = stblBox.parse()
	log.PanicIf(err)

	return stblBox, nil
}

var (
	_ boxFactory = stblBoxFactory{}
	_ CommonBox  = &StblBox{}
)

func init() {
	registerAtom(stblBoxFactory{})
}
