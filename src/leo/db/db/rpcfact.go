/* this is rpc fact */

package db

import (
	"leo/base"
)

func BuildRpcService(port *base.Port) error {
	sv, err := NewRpcService()
	if err != nil {
		return err
	}
	port.RegisterService(sv)

	//todo:

	return nil
}