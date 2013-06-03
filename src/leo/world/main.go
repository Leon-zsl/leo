/* this is the stage server
*/

package main

import (
	"fmt"
	"leo/world/world"
)

func startup() {
	w, err := world.NewWorld()
	if err != nil {
		fmt.Println("world server create failed")
	} else {
		w.Run()
	}
}

func main() {
	fmt.Println("start world server")
	startup()
	fmt.Println("master world close")
}
