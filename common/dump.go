package bmfcommon

import (
	"fmt"
	"strings"

	"github.com/dsoprea/go-logging"
)

func dump(box interface{}, level int) {
	indent := strings.Repeat("  ", level)

	switch t := box.(type) {
	case BoxChildIndexer:
		_, isRoot := t.(*BmfResource)

		if isRoot == true {
			fmt.Printf("%s> [ROOT]\n", indent)
		} else {
			cb := t.(CommonBox)
			fmt.Printf("%s> %s  %s\n", indent, cb.Name(), cb.InlineString())
		}

		names := t.ChildrenTypes()
		for _, name := range names {
			boxes, err := t.GetChildBoxes(name)
			log.PanicIf(err)

			for _, currentBox := range boxes {
				dump(currentBox, level+1)
			}
		}
	case CommonBox:
		fmt.Printf("%s> %s  %s\n", indent, t.Name(), t.InlineString())
	}
}

// Dump prints the box hierarchy.
func Dump(box interface{}) {
	dump(box, 0)
}
