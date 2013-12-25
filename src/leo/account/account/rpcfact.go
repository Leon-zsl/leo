/* this is rpc fact */

package account

import (
	"leo/base"
)

func BuildRpcService(port *base.Port) error {
	_, err := NewClientReqDispatcher(port)
	if err != nil {
		return err
	}

	//todo:

	return nil
}
