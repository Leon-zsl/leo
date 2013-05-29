/* this is server main
*/

package main

import (
	"fmt"
	"leo/server/server"
)

func main() {
	app, err := server.NewApp()
	if err != nil {
		fmt.Println("new app err: ", err.Error())
		return
	}

	app.Startup()

}