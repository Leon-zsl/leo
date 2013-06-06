/* this is rpc define */

package common

import (
	"leo/base"
)

type RpcSendTo struct {
	Sid string
	Pkt *base.Packet
}

type RpcSendToAll struct {
	Pkt *base.Packet
}

type RpcBroadcast struct {
	Sids []string
	Pkt *base.Packet
}

type RpcMoveWorld struct {
	Sid string
	PortID int
}