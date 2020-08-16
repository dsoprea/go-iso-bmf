package bmfcommon

import (
	"testing"
)

func TestVolume_Decode(t *testing.T) {
	numerator := uint16(11)
	denominator := uint16(22)

	fp := RationalUint16ToFixedPoint16(numerator, denominator, 8, 8)

	v := Volume(fp.rawValue)

	recoveredFp := v.Decode()

	if recoveredFp != fp {
		t.Fatalf("FP16 not correct.")
	}
}

func TestVolume_IsFullVolume(t *testing.T) {
	if Volume(0x0100).IsFullVolume() != true {
		t.Fatalf("IsFullVolume() was not correct.")
	}
}

func TestVolume_String_Fractional(t *testing.T) {
	v := Volume(0x1234)

	if v.String() != "0.3%" {
		t.Fatalf("String() was not correct: [%s]", v.String())
	}
}

func TestVolume_String_Full(t *testing.T) {
	v := Volume(0x0100)

	if v.String() != "FULL" {
		t.Fatalf("String() was not correct: [%s]", v.String())
	}
}

func TestVolume_String_Off(t *testing.T) {
	v := Volume(0x0)

	if v.String() != "OFF" {
		t.Fatalf("String() was not correct: [%s]", v.String())
	}
}
