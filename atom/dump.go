package atom

import (
	"fmt"
	"strings"

	"github.com/dsoprea/go-logging"
)

func dump(box interface{}, level int) {
	switch t := box.(type) {
	case BoxChildIndexer:
		names := t.ChildrenTypes()
		for _, name := range names {
			boxes, err := t.GetChildBoxes(name)
			log.PanicIf(err)

			for _, currentBox := range boxes {
				dump(currentBox, level+1)
			}
		}
	case CommonBox:
		// We never print on level (0), so we decrement by one to avoid
		// unnecessary indentation.
		indent := strings.Repeat("  ", level-1)

		fmt.Printf("%s%s %s\n", indent, t.Name(), t.InlineString())
	}
}

func Dump(box interface{}) {
	dump(box, 0)
}
