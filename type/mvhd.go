package bmftype

import (
	"fmt"

	"github.com/dsoprea/go-logging"

	"github.com/dsoprea/go-iso-bmf/common"
)

// MvhdRate represents the playback speed as a proportion of normal
// speed.
type MvhdRate uint32

// Decode returns the deconstructed value.
func (rate MvhdRate) Decode() bmfcommon.FixedPoint32 {
	return bmfcommon.Uint32ToFixedPoint32(uint32(rate), 16, 16)
}

// IsFullSpeed returns true if playback is running at normal speed.
func (rate MvhdRate) IsFullSpeed() bool {
	return rate == 0x00010000
}

// String returns a human representation of the volume.
func (rate MvhdRate) String() string {
	if rate.IsFullSpeed() == true {
		return "NORMAL"
	}

	return fmt.Sprintf("%.1f%%", rate.Decode().Float())
}

// MvhdBox is the "Movie Header" box.
//
// This box defines overall information which is media-independent,
// and relevant to the entire presentationconsidered as a whole.
type MvhdBox struct {
	bmfcommon.Box
	bmfcommon.Standard32TimeSupport

	flags   uint32
	version uint8
	rate    MvhdRate
	volume  bmfcommon.Volume
}

// Flags returns the flags of the box. The first byte is the version.
func (mb *MvhdBox) Flags() uint32 {
	return mb.flags
}

// Version returns the version of the box
func (mb *MvhdBox) Version() uint8 {
	return mb.version
}

// Rate returns the playback rate.
func (mb *MvhdBox) Rate() MvhdRate {
	return mb.rate
}

// Volume returns the audio volume.
func (mb *MvhdBox) Volume() bmfcommon.Volume {
	return mb.volume
}

// InlineString returns an undecorated string of field names and values.
func (mb *MvhdBox) InlineString() string {
	return fmt.Sprintf(
		"%s VER=(0x%02x) FLAGS=(0x%08x) RATE=(%d]) VOLUME=[%s] %s",
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

	var creationEpoch uint64
	var modificationEpoch uint64
	var timeScale uint64
	var duration uint64

	if b.version == 0 {
		creationEpoch32 := bmfcommon.DefaultEndianness.Uint32(data[4:8])
		creationEpoch = uint64(creationEpoch32)

		modificationEpoch32 := bmfcommon.DefaultEndianness.Uint32(data[8:12])
		modificationEpoch = uint64(modificationEpoch32)

		timeScale32 := bmfcommon.DefaultEndianness.Uint32(data[12:16])
		timeScale = uint64(timeScale32)

		duration32 := bmfcommon.DefaultEndianness.Uint32(data[16:20])
		duration = uint64(duration32)

		b.rate = MvhdRate(bmfcommon.DefaultEndianness.Uint32(data[20:24]))
		b.volume = bmfcommon.Volume(bmfcommon.DefaultEndianness.Uint16(data[24:26]))
	} else if b.version == 1 {
		creationEpoch = bmfcommon.DefaultEndianness.Uint64(data[4:12])
		modificationEpoch = bmfcommon.DefaultEndianness.Uint64(data[12:20])
		timeScale = bmfcommon.DefaultEndianness.Uint64(data[20:28])
		duration = bmfcommon.DefaultEndianness.Uint64(data[28:36])

		b.rate = MvhdRate(bmfcommon.DefaultEndianness.Uint32(data[36:40]))
		b.volume = bmfcommon.Volume(bmfcommon.DefaultEndianness.Uint16(data[40:42]))
	} else {
		log.Panicf("mvhd: version (%d) not supported", b.version)
	}

	b.Standard32TimeSupport = bmfcommon.NewStandard32TimeSupport(
		creationEpoch,
		modificationEpoch,
		uint32(duration),
		uint32(timeScale))

	return nil
}

type mvhdBoxFactory struct {
}

// Name returns the name of the type.
func (mvhdBoxFactory) Name() string {
	return "mvhd"
}

// New returns a new value instance.
func (mvhdBoxFactory) New(box bmfcommon.Box) (cb bmfcommon.CommonBox, childBoxSeriesOffset int, err error) {
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

	return mvhdBox, -1, nil
}

var (
	_ bmfcommon.BoxFactory = mvhdBoxFactory{}
	_ bmfcommon.CommonBox  = &MvhdBox{}
)

func init() {
	bmfcommon.RegisterBoxType(mvhdBoxFactory{})
}
