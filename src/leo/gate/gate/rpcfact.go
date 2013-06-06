/* this is rpc service fact */

package gate

import (
	"leo/base"
)

func BuildRpcService(port *base.Port) error {
	sv, err := NewRpcService()
	if err != nil {
		return err
	}
	port.RegisterService(sv)

	return nil
}