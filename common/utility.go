package bmfcommon

import (
	"fmt"
	"log"
	"reflect"
)

// bmfcommon.DumpBytes prints a list of hex-encoded bytes.
func DumpBytes(data []byte) {

	// TODO(dustin): Add test

	fmt.Printf("DUMP: ")
	for _, x := range data {
		fmt.Printf("%02x ", x)
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

	DefaultEndianness.PutUint32(
		(*buffer)[start:start+4],
		uint32(len(data))+uint32(BoxHeaderSize))

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
