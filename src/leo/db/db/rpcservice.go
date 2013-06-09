/* this is rpc service */

package db

import (
//	"fmt"
	"leo/base"
	"leo/common"
)

type RpcService struct {
}

func NewRpcService() (*RpcService, error) {
	sv := new(RpcService)
	err := sv.init()
	return sv, err
}

func (sv *RpcService) init() error {
	return nil
}

func (sv *RpcService) Get(q *common.DBGet, r *common.DBGetResp) error {
	if q == nil {
		return base.NewLeoError(common.LeoErrRpcNoArg, "")
	}
	rcd, e := Root.Driver.Get(q.Table, q.Key, q.Keyname)
	if e != nil && rcd != nil {
		r.Ok = true
		r.Value = rcd
	} else {
		r.Ok = false
	}
	return e
}

func (sv *RpcService) Set(q *common.DBSet, v *int) error {
	if q == nil {
		return base.NewLeoError(common.LeoErrRpcNoArg, "")
	}
	*v = 0
	return Root.Driver.Set(q.Table, q.Key, q.Keyname, q.Value)
}

func (sv *RpcService) Add(q *common.DBAdd, v *int) error {
	if q == nil {
		return base.NewLeoError(common.LeoErrRpcNoArg, "")
	}
	*v = 0
	return Root.Driver.Add(q.Table, q.Key, q.Keyname, q.Value)
}

func (sv *RpcService) Del(q *common.DBDel, v *int) error {
	if q == nil {
		return base.NewLeoError(common.LeoErrRpcNoArg, "")
	}
	*v = 0
	return Root.Driver.Del(q.Table, q.Key, q.Keyname)
}

func (sv *RpcService) Query(q *common.DBQuery, r *common.DBQueryResp) error {
	if q == nil {
		return base.NewLeoError(common.LeoErrRpcNoArg, "")
	}
	rcds, e := Root.Driver.Query(q.Sql, q.Args)
	if e == nil && rcds != nil {
		r.Ok = true
		r.Values = rcds
	} else {
		r.Ok = false
	}
	return e
}

func (sv *RpcService) Count(table string, count *int) error {
	*count = Root.Driver.Count(table)
	return nil
}