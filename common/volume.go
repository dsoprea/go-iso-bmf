package bmfcommon

import (
	"fmt"
)

// Volume represents the audio volume of the media.
type Volume uint16

// Decode returns the deconstructed value.
func (volume Volume) Decode() FixedPoint16 {
	return Uint16ToFixedPoint16(uint16(volume), 8, 8)
}

// IsFullVolume returns true if the volume is at maximum.
func (volume Volume) IsFullVolume() bool {
	return volume == 0x0100
}

// String returns a human representation of the volume.
func (volume Volume) String() string {

	// TODO(dustin): Add test once we better understand how to represent a fractional value.

	if volume.IsFullVolume() == true {
		return "FULL"
	} else if volume == 0 {
		return "OFF"
	}

	fp16 := volume.Decode()
	numerator, denominator := fp16.Rational()

	return fmt.Sprintf("%d/%d (%.1f%%)", numerator, denominator, fp16.Float())
}
