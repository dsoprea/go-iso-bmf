package bmftype

import (
	"github.com/dsoprea/go-iso-bmf/common"
)

func writeIlocExtentVersion032bitBytes(data *[]byte, extentOffset, extentLength uint32) {
	bmfcommon.PushBytes(data, extentOffset)
	bmfcommon.PushBytes(data, extentLength)
}

func writeIlocExtentVersion1OrVersion232bitBytes(data *[]byte, extentIndex uint32, extentOffset, extentLength uint32) {
	bmfcommon.PushBytes(data, extentIndex)
	bmfcommon.PushBytes(data, extentOffset)
	bmfcommon.PushBytes(data, extentLength)
}

func writeIlocItemVersion032bitBytes(data *[]byte, itemId uint16, dataReferenceIndex uint16, baseOffset uint32, extents []IlocExtent) {
	bmfcommon.PushBytes(data, itemId)
	bmfcommon.PushBytes(data, dataReferenceIndex)
	bmfcommon.PushBytes(data, baseOffset)

	extentCount := uint16(len(extents))
	bmfcommon.PushBytes(data, extentCount)

	for _, ie := range extents {
		writeIlocExtentVersion032bitBytes(
			data,
			uint32(ie.extentOffset),
			uint32(ie.extentLength))
	}
}

func writeIlocItemVersion132bitBytes(data *[]byte, itemId uint16, constructionMethod uint16, dataReferenceIndex uint16, baseOffset uint32, extents []IlocExtent) {
	bmfcommon.PushBytes(data, itemId)
	bmfcommon.PushBytes(data, constructionMethod)
	bmfcommon.PushBytes(data, dataReferenceIndex)
	bmfcommon.PushBytes(data, baseOffset)

	extentCount := uint16(len(extents))
	bmfcommon.PushBytes(data, extentCount)

	for _, ie := range extents {
		writeIlocExtentVersion1OrVersion232bitBytes(
			data,
			uint32(ie.extentIndex),
			uint32(ie.extentOffset),
			uint32(ie.extentLength))
	}
}

func writeIlocItemVersion232bitBytes(data *[]byte, itemId uint32, constructionMethod uint16, dataReferenceIndex uint16, baseOffset uint32, extents []IlocExtent) {
	bmfcommon.PushBytes(data, itemId)
	bmfcommon.PushBytes(data, constructionMethod)
	bmfcommon.PushBytes(data, dataReferenceIndex)
	bmfcommon.PushBytes(data, baseOffset)

	extentCount := uint16(len(extents))
	bmfcommon.PushBytes(data, extentCount)

	for _, ie := range extents {
		writeIlocExtentVersion1OrVersion232bitBytes(
			data,
			uint32(ie.extentIndex),
			uint32(ie.extentOffset),
			uint32(ie.extentLength))
	}
}
