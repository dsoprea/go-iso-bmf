package atom

import (
	"github.com/dsoprea/go-logging"
)

// TrakBox - Track Box
// Box Type: tkhd
// Container: Movie Box (moov)
// Mandatory: Yes
// Quantity: One or more.
type TrakBox struct {
	*Box

	// SamplesDuration
	// SamplesSize
	// SampleGroupsInfo

	Tkhd *TkhdBox
	Mdia *MdiaBox
	Edts *EdtsBox

	// chunks []Chunk
	// samples []Sample
}

func (b *TrakBox) parse() (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	boxes, err := b.Box.readBoxes(0)
	log.PanicIf(err)

	for _, box := range boxes {
		switch box.Name() {
		case "tkhd":
			b.Tkhd = &TkhdBox{Box: box}
			b.Tkhd.parse()

		case "mdia":
			b.Mdia = &MdiaBox{Box: box}
			b.Mdia.parse()

		case "edts":
			b.Edts = &EdtsBox{Box: box}
			b.Edts.parse()
		}
	}

	return nil
}

// func (b *TrakBox) Size() (sz int) {
//     sz += b.Tkhd.Size
// 	boxes := readBoxes(b.File, b.Start+BoxHeaderSize, b.Size-BoxHeaderSize)
// }

type trakBoxFactory struct {
}

// Name returns the name of the type.
func (trakBoxFactory) Name() string {
	return "trak"
}

// New returns a new value instance.
func (trakBoxFactory) New(box *Box) (cb CommonBox, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	trakBox := &TrakBox{
		Box: box,
	}

	err = trakBox.parse()
	log.PanicIf(err)

	return trakBox, nil
}

var (
	_ boxFactory = trakBoxFactory{}
	_ CommonBox  = TrakBox{}
)

func init() {
	registerAtom(trakBoxFactory{})
}
