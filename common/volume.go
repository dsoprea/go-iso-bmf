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
	if volume.IsFullVolume() == true {
		return "FULL"
	} else if volume == 0 {
		return "OFF"
	}

	return fmt.Sprintf("%.1f%%", volume.Decode().Float())
}
