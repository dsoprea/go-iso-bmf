package atom

import (
	"github.com/dsoprea/go-logging"
)

// MdhdBox is a "Media Header" box.
//
// The media header declares overall information that is media-independent,
// and relevant to characteristics of the media in a track.
type MdhdBox struct {
	Box

	version          byte
	flags            uint32
	creationTime     uint32
	modificationTime uint32
	timescale        uint32
	duration         uint32
	language         uint16
	languageString   string
}

func (mb *MdhdBox) Version() byte {
	return mb.version
}

func (mb *MdhdBox) Flags() uint32 {
	return mb.flags
}

func (mb *MdhdBox) CreationTime() uint32 {
	return mb.creationTime
}

func (mb *MdhdBox) ModificationTime() uint32 {
	return mb.modificationTime
}

func (mb *MdhdBox) Timescale() uint32 {
	return mb.timescale
}

func (mb *MdhdBox) Duration() uint32 {
	return mb.duration
}

func (mb *MdhdBox) Language() uint16 {
	return mb.language
}

func (mb *MdhdBox) LanguageString() string {
	return mb.languageString
}

func (b *MdhdBox) parse() (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	data, err := b.readBoxData()
	log.PanicIf(err)

	b.version = data[0]
	b.flags = defaultEndianness.Uint32(data[0:4])
	b.creationTime = defaultEndianness.Uint32(data[4:8])
	b.modificationTime = defaultEndianness.Uint32(data[8:12])
	b.timescale = defaultEndianness.Uint32(data[12:16])
	b.duration = defaultEndianness.Uint32(data[16:20])
	b.language = defaultEndianness.Uint16(data[20:22])
	b.languageString = b.getLanguageString()

	return nil
}

func (b *MdhdBox) getLanguageString() string {
	var l [3]uint8

	l[0] = uint8((b.language >> 10) & 0x1F)
	l[1] = uint8((b.language >> 5) & 0x1F)
	l[2] = uint8((b.language) & 0x1F)

	return string([]byte{l[0] + 0x60, l[1] + 0x60, l[2] + 0x60})
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
