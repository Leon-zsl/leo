/* this is test client
*/

package main

import (
	"fmt"
	"leo/client/client"
)

func main () {
	app, err := client.NewApp()
	if err != nil {
		fmt.Println("new app err: ", err.Error())
		return
	}

	app.Startup()
}