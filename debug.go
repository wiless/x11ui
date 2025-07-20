package x11ui

import "log"

var DEBUG_LEVEL = 0

func deBug(prefix string, err error) {
	if DEBUG_LEVEL != 0 && err != nil {
		log.Printf(prefix, err)
	}
}
