package boxtype

import (
	"fmt"

	"github.com/dsoprea/go-logging"
)

// MdhdBox is a "Media Header" box.
//
// The media header declares overall information that is media-independent,
// and relevant to characteristics of the media in a track.
type MdhdBox struct {
	Box
	standard32TimeSupport

	version  byte
	flags    uint32
	language uint16
}

// Version returns the version.
func (mb *MdhdBox) Version() byte {
	return mb.version
}

// Flags returns the flags.
func (mb *MdhdBox) Flags() uint32 {
	return mb.flags
}

// Language returns the stringified language.
func (mb *MdhdBox) Language() string {
	languageString := mb.getLanguageString(mb.language)

	return languageString
}

// String returns a descriptive string.
func (mb *MdhdBox) String() string {
	return fmt.Sprintf("mdhd<%s>", mb.InlineString())
}

// InlineString returns an undecorated string of field names and values.
func (mb *MdhdBox) InlineString() string {
	return fmt.Sprintf(
		"%s VER=(0x%02x) FLAGS=(0x%08x) LANG=[%s] %s",
		mb.Box.InlineString(), mb.version, mb.flags, mb.Language(),
		mb.standard32TimeSupport.InlineString())
}

func (b *MdhdBox) parse() (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	data, err := b.readBoxData()
	log.PanicIf(err)

	// TODO(dustin): In the past, we've made some changes regarding the overlap between the "version" byte and whatever comes behind it. However, this box decodes pairly in several respects if we "fix" it.

	b.version = data[0]
	b.flags = defaultEndianness.Uint32(data[0:4])

	creationEpoch := defaultEndianness.Uint32(data[4:8])
	modificationEpoch := defaultEndianness.Uint32(data[8:12])
	timeScale := defaultEndianness.Uint32(data[12:16])
	duration := defaultEndianness.Uint32(data[16:20])

	b.standard32TimeSupport = newStandard32TimeSupport(
		creationEpoch,
		modificationEpoch,
		duration,
		timeScale)

	b.language = defaultEndianness.Uint16(data[20:22])

	return nil
}

func (b *MdhdBox) getLanguageString(language uint16) string {
	var l [3]uint8

	l[0] = uint8((language >> 10) & 0b11111)
	l[1] = uint8((language >> 5) & 0b11111)
	l[2] = uint8((language) & 0b11111)

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
	RegisterBoxType(mdhdBoxFactory{})
}
