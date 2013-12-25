/* this is rpc define */

package common

import (
	"leo/base"
)

//db service args
type DBQuery struct {
	Sql  string
	Args []interface{}
}

type DBQueryResp struct {
	Ok     bool
	Values []*base.Record
}

type DBGet struct {
	Table   string
	Key     int
	Keyname string
}

type DBGetResp struct {
	Ok    bool
	Value *base.Record
}

type DBSet struct {
	DBGet
	Value *base.Record
}

type DBAdd struct {
	DBSet
}

type DBDel struct {
	DBGet
}

//client router args
type RpcClientRequest struct {
	Sid string
	Pkt *base.Packet
}

type RpcSendTo struct {
	Sid string
	Pkt *base.Packet
}

type RpcSendToAll struct {
	Pkt *base.Packet
}

type RpcBroadcast struct {
	Sids []string
	Pkt  *base.Packet
}

//gate service args
type RpcMoveWorld struct {
	Sid    string
	PortID int
}
