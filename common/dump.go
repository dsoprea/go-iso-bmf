package bmfcommon

import (
	"fmt"
	"strings"

	"github.com/dsoprea/go-logging"
)

func dump(box interface{}, level int) {

	// TODO(dustin): Add test

	// TODO(dustin): Not sure if this makes sense. Since all Box structs are a BoxChildIndexer, we're not sure if the CommonBox condition is being hit. On the other hand, things should fail if we try to get children on a non-children box. However, nothing iscurrently failing. However, the dump output is not correct (based on the parents that we print from the New() functions). Plus, we get some weird disordering *whenever we print that parent* from the New() functions.
	// -> However, maybe we were just predeterming the children, storing into an inex initial, and then just enumerating this from the dump function every time. This would imply that this is all simply a bug in the indexing routine.

	indent := strings.Repeat("  ", level)

	switch t := box.(type) {
	case BoxChildIndexer:
		_, isRoot := t.(*File)

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

func Dump(box interface{}) {

	// TODO(dustin): Add test

	dump(box, 0)
}
