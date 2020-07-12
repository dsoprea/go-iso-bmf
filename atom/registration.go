package atom

import (
	"github.com/dsoprea/go-logging"
)

var (
	boxMapping = make(map[string]boxFactory)
)

func registerAtom(bf boxFactory) {
	name := bf.Name()

	if _, found := boxMapping[name]; found == true {
		log.Panicf("box-factory already registered: [%s]", name)
	}

	boxMapping[name] = bf
}
