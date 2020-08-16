package bmftype

import (
	"bufio"
	"bytes"
	"fmt"
	"unicode"

	"encoding/binary"

	"github.com/dsoprea/go-logging"

	"github.com/dsoprea/go-iso-bmf/common"
)

// InfeItemType allows simple handling of item-types. It can be compared as a
// uint32 or as a string, adds a few boolean tests, and allows exporting as a
// simple string.
type InfeItemType uint32

// InfeItemTypeFromBytes converts four-bytes to a uint32. This supports testing.
func InfeItemTypeFromBytes(typeBytes [4]byte) InfeItemType {
	return InfeItemType(bmfcommon.DefaultEndianness.Uint32(typeBytes[:]))
}

// EqualsName encodes the integer back to the original byte-order, convert to a
// string, and compare. The item-type can be interpreted as both an integer and
// a string.
func (iit InfeItemType) EqualsName(s string) bool {
	return s == iit.string()
}

func (iit InfeItemType) string() string {
	b := make([]byte, 4)
	bmfcommon.DefaultEndianness.PutUint32(b, uint32(iit))

	return string(b)
}

// IsMime returns true if MIME.
func (iit InfeItemType) IsMime() bool {
	return iit.EqualsName("mime")
}

// IsUri returns true if URI.
func (iit InfeItemType) IsUri() bool {
	return iit.EqualsName("uri ")
}

// String returns the ASCII equivalent if all printable characters or else the
// hex representation.
func (iit InfeItemType) String() string {
	s := iit.string()

	for _, r := range s {
		if unicode.IsPrint(r) == false {
			return fmt.Sprintf("TYPE<0x%08x>", uint32(iit))
		}
	}

	return s
}

// InfeBox is the ItemInfoEntry box.
type InfeBox struct {
	// Box is the base inner box.
	bmfcommon.Box

	version             byte
	itemId              uint32
	itemProtectionIndex uint16
	itemName            string
	contentType         string
	contentEncoding     string

	extensionType uint32

	itemType InfeItemType

	itemUriType string
}

// ItemId returns the ID of the item.
func (infe *InfeBox) ItemId() uint32 {
	return infe.itemId
}

// ItemProtectionIndex returns the item protection index, which supports
// protecting the values.
func (infe *InfeBox) ItemProtectionIndex() uint16 {
	return infe.itemProtectionIndex
}

// ItemName returns the item name.
func (infe *InfeBox) ItemName() string {
	return infe.itemName
}

// ContentType returns the content type (if a MIME type).
func (infe *InfeBox) ContentType() string {
	return infe.contentType
}

// ContentEncoding returns the content encoding (if a MIME type).
func (infe *InfeBox) ContentEncoding() string {
	return infe.contentEncoding
}

// ExtensionType returns the extension type.
func (infe *InfeBox) ExtensionType() uint32 {
	return infe.extensionType
}

// ItemType returns the item-type.
func (infe *InfeBox) ItemType() InfeItemType {
	return infe.itemType
}

// ItemUriType returns the URI type (if a URI).
func (infe *InfeBox) ItemUriType() string {
	return infe.itemUriType
}

// InlineString returns an undecorated string of field names and values.
func (infe *InfeBox) InlineString() string {
	var extTypePhrase string

	if infe.version == 1 {
		extTypePhrase = fmt.Sprintf(" EXT-TYPE=(%d)", infe.extensionType)
	}

	var mimePhrase string
	var uriPhrase string

	if infe.version >= 2 {
		if infe.itemType.IsMime() == true {
			mimePhrase = fmt.Sprintf(
				" CONTENT-TYPE=[%s] CONTENT-ENCODING=[%s]",
				infe.contentType, infe.contentEncoding)
		}

		if infe.itemType.IsUri() == true {
			uriPhrase = fmt.Sprintf(
				" URI-TYPE=[%s]",
				infe.itemUriType)
		}
	}

	return fmt.Sprintf(
		"%s VER=(%d) ITEM-ID=(%d) PROTECTION-INDEX=(%d) NAME=[%s] ITEM-TYPE=[%s]%s%s%s",
		infe.Box.InlineString(), infe.version, infe.itemId, infe.itemProtectionIndex, infe.itemName, infe.itemType, extTypePhrase, mimePhrase, uriPhrase)
}

type infeBoxFactory struct {
}

// Name returns the name of the type.
func (infeBoxFactory) Name() string {
	return "infe"
}

// New returns a new value instance.
func (infeBoxFactory) New(box bmfcommon.Box) (cb bmfcommon.CommonBox, childBoxSeriesOffset int, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	data, err := box.ReadBoxData()
	log.PanicIf(err)

	infe := &InfeBox{
		Box:     box,
		version: data[0],
	}

	if infe.version > 3 {
		// NOTE(dustin): The spec implies that we'll maintain a lot of the structure below in future versions, but, obviously, new versions will carry changes and we are cynical that what we'd have would still work.
		log.Panicf("infe: version (%d) not supported", infe.version)
	}

	b := bytes.NewBuffer(data[4:])
	br := bufio.NewReader(b)

	if infe.version == 0 || infe.version == 1 {
		// itemId

		var itemId16 uint16

		err = binary.Read(br, bmfcommon.DefaultEndianness, &itemId16)
		log.PanicIf(err)

		infe.itemId = uint32(itemId16)

		// itemProtectionIndex

		err = binary.Read(br, bmfcommon.DefaultEndianness, &infe.itemProtectionIndex)
		log.PanicIf(err)

		if infe.itemProtectionIndex != 0 {
			log.Panicf("infe: protection not currently supported; please create an issue: (%d)", infe.itemProtectionIndex)
		}

		// itemName

		itemNameRaw, err := br.ReadString(0)
		log.PanicIf(err)

		infe.itemName = itemNameRaw[:len(itemNameRaw)-1]

		// contentType

		contentTypeRaw, err := br.ReadString(0)
		log.PanicIf(err)

		infe.contentType = contentTypeRaw[:len(contentTypeRaw)-1]

		// contentEncoding

		contentEncodingRaw, err := br.ReadString(0)
		log.PanicIf(err)

		infe.contentEncoding = contentEncodingRaw[:len(contentEncodingRaw)-1]

		if infe.version == 1 {
			err := binary.Read(br, bmfcommon.DefaultEndianness, &infe.extensionType)
			log.PanicIf(err)
		}
	}

	if infe.version >= 2 {
		// itemId

		if infe.version == 2 {
			var itemId16 uint16

			err = binary.Read(br, bmfcommon.DefaultEndianness, &itemId16)
			log.PanicIf(err)

			infe.itemId = uint32(itemId16)
		} else if infe.version == 3 {
			err := binary.Read(br, bmfcommon.DefaultEndianness, &infe.itemId)
			log.PanicIf(err)
		} else {
			log.Panicf("infe: version (%d) not supported", infe.version)
		}

		// itemProtectionIndex

		err = binary.Read(br, bmfcommon.DefaultEndianness, &infe.itemProtectionIndex)
		log.PanicIf(err)

		// itemType

		err = binary.Read(br, bmfcommon.DefaultEndianness, &infe.itemType)
		log.PanicIf(err)

		// itemName

		itemNameRaw, err := br.ReadString(0)
		log.PanicIf(err)

		infe.itemName = itemNameRaw[:len(itemNameRaw)-1]

		if infe.itemType.IsMime() == true {
			// contentType

			contentTypeRaw, err := br.ReadString(0)
			log.PanicIf(err)

			infe.contentType = contentTypeRaw[:len(contentTypeRaw)-1]

			// contentEncoding

			contentEncodingRaw, err := br.ReadString(0)
			log.PanicIf(err)

			infe.contentEncoding = contentEncodingRaw[:len(contentEncodingRaw)-1]
		} else if infe.itemType.IsUri() == true {
			// itemUriType

			itemUriTypeRaw, err := br.ReadString(0)
			log.PanicIf(err)

			infe.itemUriType = itemUriTypeRaw[:len(itemUriTypeRaw)-1]
		}
	}

	// Load new struct into the item index owned by the parent IINF struct.

	iinfCommon := box.Parent()

	iinf := iinfCommon.(*IinfBox)
	iinf.loadItem(infe)

	return infe, -1, nil
}

var (
	_ bmfcommon.BoxFactory = infeBoxFactory{}
	_ bmfcommon.CommonBox  = &InfeBox{}
)

func init() {
	bmfcommon.RegisterBoxType(infeBoxFactory{})
}
