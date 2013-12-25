/* this is the database server
 */

package main

import (
	"fmt"
	"leo/db/db"
)

func startup() {
	db, err := db.NewDB()
	if err != nil {
		fmt.Println("db create failed")
	} else {
		db.Run()
	}
}

func main() {
	fmt.Println("start db server")
	startup()
	fmt.Println("db server close")
}
