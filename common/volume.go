package bmfcommon

import (
	"fmt"
)

// Volume represents the audio volume of the media.
type Volume uint16

// Decode returns the deconstructed value.
func (volume Volume) Decode() FixedPoint16 {

	// TODO(dustin): Add test

	return Uint16ToFixedPoint16(uint16(volume), 8, 8)
}

// IsFullVolume returns true if the volume is at maximum.
func (volume Volume) IsFullVolume() bool {

	// TODO(dustin): Add test

	return volume == 0x0100
}

// String returns a human representation of the volume.
func (volume Volume) String() string {

	// TODO(dustin): Add test

	if volume.IsFullVolume() == true {
		return "FULL"
	}

	fp16 := volume.Decode()

	numerator, _ := fp16.Rational()

	if numerator == 0 {
		return "OFF"
	}

	return fmt.Sprintf("%.5f", fp16.Float())
}
