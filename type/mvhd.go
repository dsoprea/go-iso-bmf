package bmftype

import (
	"fmt"

	"github.com/dsoprea/go-logging"

	"github.com/dsoprea/go-iso-bmf/common"
)

// MvhdBox is a "Movie Header" box.
//
// This box defines overall information which is media-independent,
// and relevant to the entire presentationconsidered as a whole.
type MvhdBox struct {
	bmfcommon.Box
	bmfcommon.Standard32TimeSupport

	flags   uint32
	version uint8
	rate    uint32
	volume  uint16
}

func (mb *MvhdBox) Flags() uint32 {
	return mb.flags
}

func (mb *MvhdBox) Version() uint8 {
	return mb.version
}

func (mb *MvhdBox) Rate() uint32 {
	return mb.rate
}

func (mb *MvhdBox) Volume() uint16 {
	return mb.volume
}

// InlineString returns an undecorated string of field names and values.
func (mb *MvhdBox) InlineString() string {
	return fmt.Sprintf(
		"%s VER=(0x%02x) FLAGS=(0x%08x) RATE=(%d]) VOLUME=(%d) %s",
		mb.Box.InlineString(), mb.version, mb.flags, mb.rate, mb.volume,
		mb.Standard32TimeSupport.InlineString())
}

func (b *MvhdBox) parse() (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	data, err := b.ReadBoxData()
	log.PanicIf(err)

	b.version = data[0]

	// TODO(dustin): Version 1 is 64-bit. Come back to this.
	if b.version != 0 {
		log.Panicf("mvhd: only version (0) is supported")
	}

	// TODO(dustin): !! Note that there is a discrepancy of three bytes here. The first four bytes are probably, technically, the flags as in the other boxes (with the version being the first byte, also as in the other boxes).

	creationEpoch := bmfcommon.DefaultEndianness.Uint32(data[4:8])
	modificationEpoch := bmfcommon.DefaultEndianness.Uint32(data[8:12])
	timeScale := bmfcommon.DefaultEndianness.Uint32(data[12:16])
	duration := bmfcommon.DefaultEndianness.Uint32(data[16:20])

	b.Standard32TimeSupport = bmfcommon.NewStandard32TimeSupport(
		creationEpoch,
		modificationEpoch,
		duration,
		timeScale)

	b.rate = bmfcommon.DefaultEndianness.Uint32(data[20:24])
	b.volume = bmfcommon.DefaultEndianness.Uint16(data[24:26])

	return nil
}

type mvhdBoxFactory struct {
}

// Name returns the name of the type.
func (mvhdBoxFactory) Name() string {
	return "mvhd"
}

// New returns a new value instance.
func (mvhdBoxFactory) New(box bmfcommon.Box) (cb bmfcommon.CommonBox, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	mvhdBox := &MvhdBox{
		Box: box,
	}

	err = mvhdBox.parse()
	log.PanicIf(err)

	return mvhdBox, nil
}

var (
	_ bmfcommon.BoxFactory = mvhdBoxFactory{}
	_ bmfcommon.CommonBox  = &MvhdBox{}
)

func init() {
	bmfcommon.RegisterBoxType(mvhdBoxFactory{})
}
