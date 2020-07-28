package bmftype

import (
	"fmt"

	"github.com/dsoprea/go-logging"

	"github.com/dsoprea/go-iso-bmf/common"
)

// TkhdBox is the "Track Header" box.
type TkhdBox struct {
	bmfcommon.Box
	bmfcommon.Standard32TimeSupport

	version        byte
	flags          uint32
	trackId        uint32
	layer          uint16
	alternateGroup uint16
	volume         bmfcommon.Volume
	matrix         []byte
	width          uint32
	height         uint32
}

// Version returns the version of the box
func (tb *TkhdBox) Version() byte {
	return tb.version
}

// Flags returns the flags of the box. The first byte is the version.
func (tb *TkhdBox) Flags() uint32 {
	return tb.flags
}

// TrackId returns the track-ID.
func (tb *TkhdBox) TrackId() uint32 {
	return tb.trackId
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

// Width returns the width.
func (tb *TkhdBox) Width() uint32 {
	return tb.width
}

// Height returns the height.
func (tb *TkhdBox) Height() uint32 {
	return tb.height
}

// InlineString returns an undecorated string of field names and values.
func (tb *TkhdBox) InlineString() string {
	return fmt.Sprintf(
		"%s VER=(0x%02x) FLAGS=(0x%08x) TRACK-ID=(%d) LAYER=(%d) ALT-GROUP=(%d) VOLUME=[%s] MATRIX=(%d) W=(%d) H=(%d) %s",
		tb.Box.InlineString(), tb.version, tb.flags, tb.trackId, tb.layer,
		tb.AlternateGroup(), tb.volume, len(tb.matrix), tb.width, tb.height,
		tb.Standard32TimeSupport.InlineString())
}

func (b *TkhdBox) parse(timeScale uint32) (err error) {
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

	creationEpoch := bmfcommon.DefaultEndianness.Uint32(data[4:8])
	modificationEpoch := bmfcommon.DefaultEndianness.Uint32(data[8:12])

	b.trackId = bmfcommon.DefaultEndianness.Uint32(data[12:16])

	// There are four reserved bytes here.

	duration := bmfcommon.DefaultEndianness.Uint32(data[20:24])

	b.layer = bmfcommon.DefaultEndianness.Uint16(data[32:34])
	b.alternateGroup = bmfcommon.DefaultEndianness.Uint16(data[34:36])
	b.volume = bmfcommon.Volume(bmfcommon.DefaultEndianness.Uint16(data[36:38]))
	b.matrix = data[40:76]

	widthRaw := bmfcommon.DefaultEndianness.Uint32(data[76:80])
	widthFp32 := bmfcommon.Uint32ToFixedPoint32(widthRaw, 16, 16)

	// The numerator is the width. The denominator is often (always?) zero.
	b.width, _ = widthFp32.Rational()

	heightRaw := bmfcommon.DefaultEndianness.Uint32(data[80:84])
	heightFp32 := bmfcommon.Uint32ToFixedPoint32(heightRaw, 16, 16)

	// The numerator is the width. The denominator is often (always?) zero.
	b.height, _ = heightFp32.Rational()

	b.Standard32TimeSupport = bmfcommon.NewStandard32TimeSupport(
		creationEpoch,
		modificationEpoch,
		duration,
		timeScale)

	return nil
}

type tkhdBoxFactory struct {
}

// Name returns the name of the type.
func (tkhdBoxFactory) Name() string {
	return "tkhd"
}

// New returns a new value instance.
func (tkhdBoxFactory) New(box bmfcommon.Box) (cb bmfcommon.CommonBox, childBoxSeriesOffset int, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	fbi := box.Index()

	mvhdCommonBox, found := fbi[bmfcommon.IndexedBoxEntry{"moov.mvhd", 0}]
	if found == false {
		log.Panicf("TKHD box encountered before MVHD box")
	}

	mvhd := mvhdCommonBox.(*MvhdBox)
	timeScale := mvhd.TimeScale()

	tkhdBox := &TkhdBox{
		Box: box,
	}

	err = tkhdBox.parse(timeScale)
	log.PanicIf(err)

	return tkhdBox, -1, nil
}

var (
	_ bmfcommon.BoxFactory = tkhdBoxFactory{}
	_ bmfcommon.CommonBox  = &TkhdBox{}
)

func init() {
	bmfcommon.RegisterBoxType(tkhdBoxFactory{})
}
