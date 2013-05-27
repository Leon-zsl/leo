/* this is the gate server
 */

package main

import (
	"fmt"
	"runtime"
	"leo/gate/gate"
)

func startup() {
	cpu := runtime.NumCPU()
	runtime.GOMAXPROCS(cpu)
	fmt.Println("number if cpu: ", cpu)

	gate, err := gate.NewGate()
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
