package atom

import (
	"github.com/dsoprea/go-logging"
)

// EdtsBox - Edit Box
// Box Type: edts
// Container: Track Box (trak)
// Mandatory: No
// Quantity: Zero or one
type EdtsBox struct {
	*Box

	Elst *ElstBox

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

	for _, box := range boxes {
		switch box.Name() {
		case "elst":
			b.Elst = &ElstBox{Box: box}

			err := b.Elst.parse()
			log.PanicIf(err)
		}
	}

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
func (edtsBoxFactory) New(box *Box) (cb CommonBox, err error) {
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
	_ CommonBox  = EdtsBox{}
)

func init() {
	registerAtom(edtsBoxFactory{})
}
