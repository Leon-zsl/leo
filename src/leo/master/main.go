/* this is the master server
 */

package main

import (
	"fmt"
	"leo/master/master"
)

func startup() {
	m, err := master.NewMaster()
	if err != nil {
		fmt.Println("master server create failed")
	} else {
		m.Run()
	}
}

func main() {
	fmt.Println("start master server")
	startup()
	fmt.Println("master server close")
}
