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
	// NOTE(dustin): !! We don't know how to write this test, but we put this here to at least have something for the coverage. The "full volume" value is 0x0100, which equates to a rational of (1/0), and the "off" value is (0/0). So, assuming that a fractional value is between the two, we're not quite sure what to do.

	v := Volume(uint16(0x2550))

	if v.String() != "37/80 (0.5%)" {
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
