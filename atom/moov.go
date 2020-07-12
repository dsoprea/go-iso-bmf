package atom

import (
	"github.com/dsoprea/go-logging"
)

// Flag constants.
const (
	TrackFlagEnabled   = 0x0001
	TrackFlagInMovie   = 0x0002
	TrackFlagInPreview = 0x0004
)

// MoovBox - Movie Box
// Box Type: moov
// Container: File
// Mandatory: Yes
// Quantity: Exactly one.
//
// The metadata for a presentation is stored in the single Movie Box which occurs
// at the top-level of a file. Normally this box is close to the beginning or end
// of the file, though this is not required.
type MoovBox struct {
	*Box

	Mvhd  *MvhdBox
	Traks []*TrakBox

	IsFragmented bool // check for mvex box exists
}

func (b *MoovBox) parse() (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	boxes, err := readBoxes(b.File(), b.Start()+BoxHeaderSize, b.Size()-BoxHeaderSize)
	log.PanicIf(err)

	for _, box := range boxes {
		switch box.Name() {
		case "mvhd":
			b.Mvhd = &MvhdBox{Box: box}

			err := b.Mvhd.parse()
			log.PanicIf(err)

		case "iods":
			// fmt.Println("found iods")

		case "trak":
			trak := &TrakBox{Box: box}

			err := trak.parse()
			log.PanicIf(err)

			b.Traks = append(b.Traks, trak)

		case "udta":
			// fmt.Println("found udta")

		case "mvex":
			// fmt.Println("found mvex")

			// TODO(dustin): What is this?

			b.IsFragmented = true
		}

	}

	return nil
}
