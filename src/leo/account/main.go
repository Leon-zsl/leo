/* this is the master server
*/

package main

import (
	"fmt"
	"leo/account/account"
	)

func startup() {
	a, err := account.NewAccount()
	if err != nil {
		fmt.Println("account server create failed")
	} else {
		a.Run()
	}
}

func main() {
	fmt.Println("start account server")
	startup()
	fmt.Println("account server close")
}