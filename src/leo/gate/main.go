/* this is the gate server
 */
package main

import (
	"fmt"
	"leo/gate/gate"
)

func startup() {
	gate, err := gate.NewGate()
	if err != nil {
		fmt.Println("gate server create failed")
	} else {
		gate.Start()
		gate.Run()
	}
}

func main() {
	fmt.Println("start gate server")
	startup()
	fmt.Println("gate server close")
}
