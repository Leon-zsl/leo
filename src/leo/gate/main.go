/* this is the gate server
 */

package main

import (
	"fmt"
	"leo/gate/gate"
)

var (
	gate *Gate
)

func startup() {
	gate, err = gate.Create()

	if err != nil {
		fmt.Println("gate server start up failed")
	} else {
		gate.Run()
	}
}

func main() {
	fmt.Println("start gate server")
	startup()
	fmt.Println("gate server close")
}
