package bmftype

import (
	"fmt"

	"github.com/dsoprea/go-logging"

	"github.com/dsoprea/go-iso-bmf/common"
)

// HdlrBox is the "Handler Reference" box.
type HdlrBox struct {
	bmfcommon.Box

	version byte
	flags   uint32
	handler string

	hdlrName string
}

// Version is the box version.
func (hb *HdlrBox) Version() byte {
	return hb.version
}

// Flags are flags.
func (hb *HdlrBox) Flags() uint32 {
	return hb.flags
}

// Handler is the type of media.
func (hb *HdlrBox) Handler() string {
	return hb.handler
}

// HdlrName is an optional description for debugging.
func (hb *HdlrBox) HdlrName() string {
	return hb.hdlrName
}

// String returns a descriptive string.
func (hb *HdlrBox) String() string {
	return fmt.Sprintf("hdlr<%s>", hb.InlineString())
}

// InlineString returns an undecorated string of field names and values.
func (hb *HdlrBox) InlineString() string {
	return fmt.Sprintf(
		"%s VER=(0x%02x) FLAGS=(0x%08x) HANDLER=[%s] HDLR-NAME=(%d)[%s]",
		hb.Box.InlineString(), hb.version, hb.flags, hb.handler, len(hb.hdlrName), hb.hdlrName)
}

func (b *HdlrBox) parse() (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	data, err := b.Data()
	log.PanicIf(err)

	b.version = data[0]
	b.flags = bmfcommon.DefaultEndianness.Uint32(data[0:4])

	// Bytes 4:8 are for "pre_defined", which is not further described in the
	// specification and is assumed to be analogous to reserved bytes.

	b.handler = string(data[8:12])

	// Skip twelve bytes of reserved data, here.

	var hdlrNameBytes []byte
	for i := 24; i < len(data); i++ {
		if data[i] == 0 {
			hdlrNameBytes = data[24:i]
			break
		}
	}

	if hdlrNameBytes == nil {
		log.Panicf("hdlrName is unterminated")
	}

	b.hdlrName = string(hdlrNameBytes)

	return nil
}

type hdlrBoxFactory struct {
}

// Name returns the name of the type.
func (hdlrBoxFactory) Name() string {
	return "hdlr"
}

// New returns a new value instance.
func (hdlrBoxFactory) New(box bmfcommon.Box) (cb bmfcommon.CommonBox, childBoxSeriesOffset int, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	hdlrBox := &HdlrBox{
		Box: box,
	}

	err = hdlrBox.parse()
	log.PanicIf(err)

	return hdlrBox, -1, nil
}

var (
	_ bmfcommon.BoxFactory = hdlrBoxFactory{}
	_ bmfcommon.CommonBox  = &HdlrBox{}
)

func init() {
	bmfcommon.RegisterBoxType(hdlrBoxFactory{})
}
