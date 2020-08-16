package bmfcommon

import (
	"fmt"
	"reflect"
	"unicode"

	"github.com/dsoprea/go-logging"
)

// DumpBytes prints a list of hex-encoded bytes.
func DumpBytes(data []byte) {
	fmt.Printf("DUMP:\n")
	mod := 20
	gappedat := 5
	column := 0
	buffer := make([]byte, mod)

	flushLine := func() {
		// If this is the last line, add spacing to align with the rest.
		for i := column; i < mod; i++ {
			fmt.Printf("   ")
		}

		for i, r := range buffer {
			if unicode.IsPrint(rune(r)) == true {
				fmt.Printf("%c", r)
			} else {
				fmt.Printf(".")
			}

			if i > 0 && (i+1)%gappedat == 0 {
				fmt.Printf(" ")
			}
		}

		column = 0
		buffer = make([]byte, mod)
		fmt.Printf("\n")
	}

	for i, r := range data {
		if column == 0 {
			fmt.Printf("0x%08x ", i)
		}

		buffer[column] = r
		fmt.Printf("%02x ", r)

		if (column+1)%gappedat == 0 {
			fmt.Printf(" ")
		}

		column++

		if column >= mod {
			flushLine()
		}
	}

	if column > 0 {
		flushLine()
	}

	fmt.Printf("\n")
}

// Data64BitDescribed indicates that the data should be written with a 32-bit
// size.
type Data64BitDescribed []byte

// PushBox pushes a box to the given byte-slice pointer.
func PushBox(buffer *[]byte, name string, data interface{}) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err := errRaw.(error)
			log.Panic(err)
		}
	}()

	start := len(*buffer)

	if data == nil {
		data = make([]byte, 0)
	}

	var asBytes []byte

	boxHeaderSize := 0
	is64Bit := false

	if d64bd, ok := data.(Data64BitDescribed); ok == true {
		// 64-bit
		//
		// (4-byte size + 4-byte type + 8-byte size).
		boxHeaderSize = 16
		is64Bit = true

		asBytes = []byte(d64bd)
	} else {
		// 32-bit
		//
		// (4-byte size + 4-byte type)
		boxHeaderSize = 8

		asBytes = data.([]byte)
	}

	extension := make([]byte, boxHeaderSize+len(asBytes))
	*buffer = append(*buffer, extension...)

	if is64Bit == true {
		DefaultEndianness.PutUint32(
			(*buffer)[start:start+4],
			1)

		copy((*buffer)[start+4:start+8], []byte(name))

		DefaultEndianness.PutUint64(
			(*buffer)[start+8:start+16],
			uint64(len(asBytes))+uint64(boxHeaderSize))

		copy((*buffer)[start+16:], asBytes)
	} else {
		DefaultEndianness.PutUint32(
			(*buffer)[start:start+4],
			uint32(len(asBytes))+uint32(boxHeaderSize))

		copy((*buffer)[start+4:], []byte(name))
		copy((*buffer)[start+8:], asBytes)
	}
}

// PushBytes encodes the given integer and pushes to the byte-slice pointer.
func PushBytes(buffer *[]byte, x interface{}) {
	var encoded []byte

	if u8, ok := x.(uint8); ok == true {
		*buffer = append(*buffer, u8)
	} else if u16, ok := x.(uint16); ok == true {
		encoded = make([]byte, 2)

		DefaultEndianness.PutUint16(
			encoded,
			u16)
	} else if u32, ok := x.(uint32); ok == true {
		encoded = make([]byte, 4)

		DefaultEndianness.PutUint32(
			encoded,
			u32)
	} else if u64, ok := x.(uint64); ok == true {
		encoded = make([]byte, 8)

		DefaultEndianness.PutUint64(
			encoded,
			u64)
	} else if bs, ok := x.([]byte); ok == true {
		*buffer = append(*buffer, bs...)
	} else {
		log.Panicf("can not encode [%v] [%v]", reflect.TypeOf(x), x)
	}

	*buffer = append(*buffer, encoded...)
}
