/* this is rpc fact */

package account

import (
	"leo/base"
)

func BuildRpcService(port *base.Port) error {
	sv, err := NewClientReqService()
	if err != nil {
		return err
	}
	port.RegisterService(sv)

	//todo:

	return nil
}