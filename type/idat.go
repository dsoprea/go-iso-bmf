package bmftype

import (
	"fmt"

	"github.com/dsoprea/go-logging"

	"github.com/dsoprea/go-iso-bmf/common"
)

// IdatBox is a "Handler Reference" box.
type IdatBox struct {
	bmfcommon.Box

	data []byte
}

// InlineString returns an undecorated string of field names and values.
func (idat *IdatBox) InlineString() string {

	// TODO(dustin): Add test

	return fmt.Sprintf(
		"%s DATA-SIZE=(%d)",
		idat.Box.InlineString(), len(idat.data))
}

type idatBoxFactory struct {
}

// Name returns the name of the type.
func (idatBoxFactory) Name() string {

	// TODO(dustin): Add test

	return "idat"
}

// New returns a new value instance.
func (idatBoxFactory) New(box bmfcommon.Box) (cb bmfcommon.CommonBox, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	// TODO(dustin): Add test

	data, err := box.ReadBoxData()
	log.PanicIf(err)

	idat := &IdatBox{
		Box:  box,
		data: data,
	}

	return idat, nil
}

var (
	_ bmfcommon.BoxFactory = idatBoxFactory{}
	_ bmfcommon.CommonBox  = &IdatBox{}
)

func init() {
	bmfcommon.RegisterBoxType(idatBoxFactory{})
}