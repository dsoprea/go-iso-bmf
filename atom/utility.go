package atom

import (
	"fmt"
)

// DumpBytes prints a list of hex-encoded bytes.
func DumpBytes(data []byte) {
	fmt.Printf("DUMP: ")
	for _, x := range data {
		fmt.Printf("%02x ", x)
	}

	fmt.Printf("\n")
}
