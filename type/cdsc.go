package bmftype

import (
	"fmt"
	"sort"
	"strings"

	"github.com/dsoprea/go-logging"

	"github.com/dsoprea/go-iso-bmf/common"
)

// CdscBox is a "Item Reference" box.
type CdscBox struct {
	bmfcommon.Box

	version    byte
	fromItemId uint32
	toItemIds  []uint32
}

// InlineString returns an undecorated string of field names and values.
func (cdsc *CdscBox) InlineString() string {
	toItemIds := make([]int, len(cdsc.toItemIds))

	for i, toItemId := range cdsc.toItemIds {
		toItemIds[i] = int(toItemId)
	}

	sort.Ints(toItemIds)

	toItemIdsPhrases := make([]string, len(toItemIds))
	for i, toItemId := range toItemIds {
		toItemIdsPhrases[i] = fmt.Sprintf("%d", toItemId)
	}

	toItemIdsPhrase := strings.Join(toItemIdsPhrases, ",")

	return fmt.Sprintf(
		"%s VER=(%d) FROM-ITEM-ID=(%d) TO-ITEM-IDS=(%d)[%v]",
		cdsc.Box.InlineString(), cdsc.version, cdsc.fromItemId, len(toItemIdsPhrases), toItemIdsPhrase)
}

type cdscBoxFactory struct {
}

// Name returns the name of the type.
func (cdscBoxFactory) Name() string {
	return "cdsc"
}

// New returns a new value instance.
//
// This contains other boxes, but the box-types are actually the reference-
// types (e.g. cdsc)..
func (cdscBoxFactory) New(box bmfcommon.Box) (cb bmfcommon.CommonBox, childBoxSeriesOffset int, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	data, err := box.ReadBoxData()
	log.PanicIf(err)

	fbi := box.Index()

	irefCommonBox, found := fbi[bmfcommon.IndexedBoxEntry{"meta.iref", 0}]
	if found == false {
		log.Panicf("CDSC box encountered before IREF box")
	}

	iref := irefCommonBox.(*IrefBox)
	irefVersion := iref.Version()

	var fromItemId uint32

	offset := 0

	if irefVersion == 0 {
		fromItemId16 := bmfcommon.DefaultEndianness.Uint16(data[offset : offset+2])
		fromItemId = uint32(fromItemId16)

		offset += 2
	} else if irefVersion == 1 {
		fromItemId = bmfcommon.DefaultEndianness.Uint32(data[offset : offset+4])
		offset += 4
	} else {
		log.Panicf("iref: version (%d) not supported (1)", irefVersion)
	}

	referenceCount := bmfcommon.DefaultEndianness.Uint16(data[offset : offset+2])

	offset += 2

	toItemIds := make([]uint32, referenceCount)

	for i := 0; i < int(referenceCount); i++ {
		var toItemId uint32

		if irefVersion == 0 {
			toItemId16 := bmfcommon.DefaultEndianness.Uint16(data[offset : offset+2])
			toItemId = uint32(toItemId16)

			offset += 2
		} else if irefVersion == 1 {
			toItemId = bmfcommon.DefaultEndianness.Uint32(data[offset : offset+4])
			offset += 4
		} else {
			log.Panicf("iref: version (%d) not supported (2)", irefVersion)
		}

		toItemIds[i] = toItemId
	}

	cdsc := &CdscBox{
		Box:        box,
		fromItemId: fromItemId,
		toItemIds:  toItemIds,
	}

	return cdsc, -1, nil
}

var (
	_ bmfcommon.BoxFactory = cdscBoxFactory{}
	_ bmfcommon.CommonBox  = &CdscBox{}
)

func init() {
	bmfcommon.RegisterBoxType(cdscBoxFactory{})
}
