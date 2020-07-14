package atom

import (
	"fmt"
	"testing"
)

func TestGetDurationString(t *testing.T) {
	s := GetDurationString(123456789, 12345)

	if s != "02:46:40:549" {
		t.Fatalf("Duration string not correct.")
	}
}

func TestFixed16(t *testing.T) {
	a := []byte{0x00, 0x00}
	b := []byte{0x01, 0x00}

	a1 := fixed16(a)
	b1 := fixed16(b)

	if a1 != 0 {
		t.Fatalf("al not correct.")
	}

	if b1 != 256 {
		fmt.Println(uint16(b1))
		t.Fatalf("bl not correct.")
	}

	if uint16(b1) != uint16(defaultEndianness.Uint16(b)) {
		t.Fatalf("bl not correct.")
	}

}

func TestFixed32(t *testing.T) {
	a := []byte{0x00, 0x01, 0x00, 0x00}

	a1 := fixed32(a)

	if a1 != 65536 {
		t.Fatalf("fixed32 not correct.")
	}

	if uint32(a1) != uint32(defaultEndianness.Uint32(a)) {
		t.Fatalf("uint32 not correct.")
	}
}
