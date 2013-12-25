/* this is test client
 */

package main

import (
	"fmt"
	"leo/robot/robot"
)

func main() {
	app, err := robot.NewApp()
	if err != nil {
		fmt.Println("new app err: ", err.Error())
		return
	}

	app.Startup()
}
