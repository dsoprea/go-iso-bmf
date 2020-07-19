package bmftype

import (
	"github.com/dsoprea/go-logging"

	"github.com/dsoprea/go-iso-bmf/common"
)

// TkhdBox is a "Track Header" box.
type TkhdBox struct {
	bmfcommon.Box

	version          byte
	flags            uint32
	creationTime     uint32
	modificationTime uint32
	trackID          uint32
	duration         uint32
	layer            uint16
	alternateGroup   uint16
	volume           bmfcommon.Volume
	matrix           []byte
	width            TkhdWidthOrHeight
	height           TkhdWidthOrHeight
}

// Version returns the version of the box
func (tb *TkhdBox) Version() byte {
	return tb.version
}

// Flags returns the flags of the box. The first byte is the version.
func (tb *TkhdBox) Flags() uint32 {
	return tb.flags
}

// CreationTime returns the creation time.
func (tb *TkhdBox) CreationTime() uint32 {

	// TODO(dustin): Finish converting this to return a time.Time .

	return tb.creationTime
}

// ModificationTime returns the modification time.
func (tb *TkhdBox) ModificationTime() uint32 {

	// TODO(dustin): Finish converting this to return a time.Time .

	return tb.modificationTime
}

// TrackID returns the track-ID.
func (tb *TkhdBox) TrackID() uint32 {
	return tb.trackID
}

// Duration returns the duration.
func (tb *TkhdBox) Duration() uint32 {
	return tb.duration
}

// Layer returns the layer.
func (tb *TkhdBox) Layer() uint16 {
	return tb.layer
}

// AlternateGroup returns the alternate group.
func (tb *TkhdBox) AlternateGroup() uint16 {
	return tb.alternateGroup
}

// Volume returns the volume.
func (tb *TkhdBox) Volume() bmfcommon.Volume {
	return tb.volume
}

// Matrix returns the matrix.
func (tb *TkhdBox) Matrix() []byte {
	return tb.matrix
}

// TkhdWidthOrHeight represents either a width or height value.
type TkhdWidthOrHeight uint32

// IsSet returns true if applicable (if non-zero).
func (twh TkhdWidthOrHeight) IsSet() bool {

	// TODO(dustin): Add test

	return twh > 0
}

// Decode returns the deconstructed value.
func (twh TkhdWidthOrHeight) Decode() bmfcommon.FixedPoint32 {

	// TODO(dustin): Add test

	return bmfcommon.Uint32ToFixedPoint32(uint32(twh), 16, 16)
}

// Width returns the width.
func (b *TkhdBox) Width() TkhdWidthOrHeight {
	return b.width
}

// Height returns the height.
func (b *TkhdBox) Height() TkhdWidthOrHeight {
	return b.height
}

func (b *TkhdBox) parse() (err error) {
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
		log.Panicf("tkhd: only version (0) is supported")
	}

	b.flags = bmfcommon.DefaultEndianness.Uint32(data[0:4])
	b.creationTime = bmfcommon.DefaultEndianness.Uint32(data[4:8])
	b.modificationTime = bmfcommon.DefaultEndianness.Uint32(data[8:12])
	b.trackID = bmfcommon.DefaultEndianness.Uint32(data[12:16])
	b.duration = bmfcommon.DefaultEndianness.Uint32(data[20:24])
	b.layer = bmfcommon.DefaultEndianness.Uint16(data[32:34])
	b.alternateGroup = bmfcommon.DefaultEndianness.Uint16(data[34:36])
	b.volume = bmfcommon.Volume(bmfcommon.DefaultEndianness.Uint16(data[36:38]))
	b.matrix = data[40:76]
	b.width = TkhdWidthOrHeight(bmfcommon.DefaultEndianness.Uint32(data[76:80]))
	b.height = TkhdWidthOrHeight(bmfcommon.DefaultEndianness.Uint32(data[80:84]))

	return nil
}

type tkhdBoxFactory struct {
}

// Name returns the name of the type.
func (tkhdBoxFactory) Name() string {
	return "tkhd"
}

// New returns a new value instance.
func (tkhdBoxFactory) New(box bmfcommon.Box) (cb bmfcommon.CommonBox, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	tkhdBox := &TkhdBox{
		Box: box,
	}

	err = tkhdBox.parse()
	log.PanicIf(err)

	return tkhdBox, nil
}

var (
	_ bmfcommon.BoxFactory = tkhdBoxFactory{}
	_ bmfcommon.CommonBox  = &TkhdBox{}
)

func init() {
	bmfcommon.RegisterBoxType(tkhdBoxFactory{})
}
