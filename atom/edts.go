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
}

func (b *EdtsBox) parse() (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	boxes, err := readBoxes(b.File(), b.Start()+BoxHeaderSize, b.Size()-BoxHeaderSize)
	log.PanicIf(err)

	for _, box := range boxes {
		switch box.Name() {
		case "elst":
			b.Elst = &ElstBox{Box: box}

			err := b.Elst.parse()
			log.PanicIf(err)
		}
	}

	return nil
}
