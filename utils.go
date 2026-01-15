package main

import (
	"fmt"
)

func usage() string {
	return fmt.Sprintf(`vj - JSON viewer version %v

Usage: vj [file]
   or: curl ... | vj

Arguments:
   -h, --help            print help
   -v, --version         print version

Key bindings:
   h, ←                  fold JSON object or array
   l, →                  unfold JSON object or array
   j, ↓                  move cursor down
   k, ↑                  move cursor up
   5j                    move cursor 5 lines down from current position
   5k                    move cursor 5 lines up from current position
   {                     move cursor to previous sibling
   }                     move cursor to next sibling
   g                     move cursor to the first line of the document
   G                     move cursor to the last line of the document
   :                     switch to command mode
   :.                    find path in JSON, for example :.users[0].email
   :q                    quit`, version,
	)
}
