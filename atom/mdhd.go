package atom

import (
	"encoding/binary"

	"github.com/dsoprea/go-logging"
)

// MdhdBox is a "Media Header" box.
//
// The media header declares overall information that is media-independent,
// and relevant to characteristics of the media in a track.
type MdhdBox struct {
	Box

	Version          byte
	Flags            uint32
	CreationTime     uint32
	ModificationTime uint32
	Timescale        uint32
	Duration         uint32
	Language         uint16
	LanguageString   string
}

func (b *MdhdBox) parse() (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	data, err := b.readBoxData()
	log.PanicIf(err)

	b.Version = data[0]
	b.Flags = binary.BigEndian.Uint32(data[0:4])
	b.CreationTime = binary.BigEndian.Uint32(data[4:8])
	b.ModificationTime = binary.BigEndian.Uint32(data[8:12])
	b.Timescale = binary.BigEndian.Uint32(data[12:16])
	b.Duration = binary.BigEndian.Uint32(data[16:20])
	b.Language = binary.BigEndian.Uint16(data[20:22])
	b.LanguageString = getLanguageString(b.Language)

	return nil
}

func getLanguageString(language uint16) string {

	// TODO(dustin): Make this a method?

	var lang [3]uint8

	lang[0] = uint8((language >> 10) & 0x1F)
	lang[1] = uint8((language >> 5) & 0x1F)
	lang[2] = uint8((language) & 0x1F)

	return string([]byte{lang[0] + 0x60, lang[1] + 0x60, lang[2] + 0x60})
}

type mdhdBoxFactory struct {
}

// Name returns the name of the type.
func (mdhdBoxFactory) Name() string {
	return "mdhd"
}

// New returns a new value instance.
func (mdhdBoxFactory) New(box Box) (cb CommonBox, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	mdhdBox := &MdhdBox{
		Box: box,
	}

	err = mdhdBox.parse()
	log.PanicIf(err)

	return mdhdBox, nil
}

var (
	_ boxFactory = mdhdBoxFactory{}
	_ CommonBox  = &MdhdBox{}
)

func init() {
	registerAtom(mdhdBoxFactory{})
}
