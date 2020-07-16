package mp4box

import (
	"fmt"
	"time"

	"github.com/dsoprea/go-logging"
)

// MdhdBox is a "Media Header" box.
//
// The media header declares overall information that is media-independent,
// and relevant to characteristics of the media in a track.
type MdhdBox struct {
	Box

	version           byte
	flags             uint32
	creationEpoch     uint32
	modificationEpoch uint32
	timeScale         uint32
	scaledDuration    uint32
	language          uint16
}

// Version returns the version.
func (mb *MdhdBox) Version() byte {
	return mb.version
}

// Flags returns the flags.
func (mb *MdhdBox) Flags() uint32 {
	return mb.flags
}

// CreationTime returns the creation time.
func (mb *MdhdBox) CreationTime() time.Time {
	t := EpochToTime(mb.creationEpoch)
	return t
}

// HasCreationTime returns true if the creation-time looks present.
func (mb *MdhdBox) HasCreationTime() bool {
	return mb.creationEpoch > 0
}

// ModificationTime returns the modification time.
func (mb *MdhdBox) ModificationTime() time.Time {
	t := EpochToTime(mb.modificationEpoch)
	return t
}

// HasModificationTime returns true if the modification-time looks present.
func (mb *MdhdBox) HasModificationTime() bool {
	return mb.modificationEpoch > 0
}

// TimeScale returns the time-scale.
func (mb *MdhdBox) TimeScale() uint32 {
	return mb.timeScale
}

// HasDuration returns true if the duration has a meaningful value.
func (mb *MdhdBox) HasDuration() bool {
	allOnes := ^uint32(0)
	return mb.scaledDuration < allOnes
}

// ScaledDuration returns the duration in timescale units (divide this number by
// the time-scale to get the number of seconds).
func (mb *MdhdBox) ScaledDuration() uint32 {
	if mb.HasDuration() == false {
		log.Panicf("duration not set (scaled-duration)")
	}

	return mb.scaledDuration
}

// Duration returns the duration as a `time.Duration`.
func (mb *MdhdBox) Duration() time.Duration {
	if mb.HasDuration() == false {
		log.Panicf("duration not set (duration)")
	}

	durationSeconds := float64(mb.scaledDuration) / float64(mb.timeScale)

	return time.Duration(durationSeconds * float64(time.Second))
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
	optional := ""

	if mb.HasCreationTime() == true {
		optional = fmt.Sprintf("%s CTIME=[%s]", optional, mb.CreationTime())
	}

	if mb.HasModificationTime() == true {
		optional = fmt.Sprintf("%s MTIME=[%s]", optional, mb.ModificationTime())
	}

	return fmt.Sprintf(
		"%s VER=(0x%02x) FLAGS=(0x%08x) DUR-S=[%.02f] LANG=[%s]%s",
		mb.Box.InlineString(), mb.version, mb.flags, float64(mb.Duration())/float64(time.Second), mb.Language(), optional)
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
	b.creationEpoch = defaultEndianness.Uint32(data[4:8])
	b.modificationEpoch = defaultEndianness.Uint32(data[8:12])
	b.timeScale = defaultEndianness.Uint32(data[12:16])
	b.scaledDuration = defaultEndianness.Uint32(data[16:20])
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
