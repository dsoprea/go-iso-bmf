package bmfcommon

import (
	"fmt"
	"log"
	"reflect"
	"unicode"
)

// DumpBytes prints a list of hex-encoded bytes.
func DumpBytes(data []byte) {

	// TODO(dustin): Add test

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

func PushBox(buffer *[]byte, name string, data []byte) {
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

	extension := make([]byte, 8+len(data))
	*buffer = append(*buffer, extension...)

	// We'll just push 32-bit box-sizes (4-byte size + 4-byte type, but no
	// follow-up 8-byte size).
	boxHeaderSize := 8

	DefaultEndianness.PutUint32(
		(*buffer)[start:start+4],
		uint32(len(data))+uint32(boxHeaderSize))

	copy((*buffer)[start+4:], []byte(name))
	copy((*buffer)[start+8:], data)
}

func PushBytes(buffer *[]byte, x interface{}) {
	var encoded []byte

	if u16, ok := x.(uint16); ok == true {
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
