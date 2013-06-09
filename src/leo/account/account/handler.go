/* this is handler interface */

package account

import pblib "code.google.com/p/goprotobuf/proto"

import (
	"fmt"
	"leo/base"
	"leo/proto"
	"leo/common"
)

type ClientReqHandler interface {
	Handle(replay *common.RpcClientRequest)
}

func SendRegisterResponse(sid string, code int32, msg string) {
	pb := &proto.RegisterResp{ ErrorCode : pblib.Int32(code), ErrorMsg : pblib.String(msg)}
	val, _ := pblib.Marshal(pb)
	pkt := base.NewPacket(proto.REGISTERRESP, val)
	resp := &common.RpcSendTo{sid, pkt}
	Root.Port.SendAsync(AccountServiceIns.GateServer(), "RpcService.SendTo", resp)
}

func SendLoginResponse(sid string, code int32, msg string) {
	pb := &proto.LoginResp{ ErrorCode : pblib.Int32(code), ErrorMsg : pblib.String(msg)}
	val, _ := pblib.Marshal(pb)
	pkt := base.NewPacket(proto.LOGINRESP, val)
	resp := &common.RpcSendTo{sid, pkt}
	Root.Port.SendAsync(AccountServiceIns.GateServer(), "RpcService.SendTo", resp)
}

type RegisterHandler struct {
}

func (h *RegisterHandler) Handle(reply *common.RpcClientRequest) {
	if reply.Pkt.Args == nil {
		SendRegisterResponse(reply.Sid, proto.EC_NO_ARG, "")
		return
	}

	dt := &proto.Register{}
	err := pblib.Unmarshal(reply.Pkt.Args, dt)
	if err != nil {
		SendRegisterResponse(reply.Sid, proto.EC_ARG_BAD_FORMAT, "")
		return
	}

	q := new(common.DBQuery)
	q.Sql = "SELECT * FROM MAIN_USER WHERE ACCOUNT=?"
	sl := make([]interface{}, 0)
	q.Args = append(sl, dt.GetName())
	r := new(common.DBQueryResp)
	r.Ok = false
	r.Values = make([]*base.Record, 0)
	err = Root.Port.Call(AccountServiceIns.DBServer(), "RpcService.Query", q, r)
	if err != nil {
		base.LoggerIns.Error("query account err:", err, dt.GetName())
		SendRegisterResponse(reply.Sid, proto.EC_SERV_ERR, "")
		return
	}
	
	if r.Ok && r.Values != nil && len(r.Values) > 0 {
		SendRegisterResponse(reply.Sid, proto.EC_ACCOUNT_EXISTS, "")
		return
	}

	var count int = 0
	err = Root.Port.Call(AccountServiceIns.DBServer(), "RpcService.Count", "main_user", &count)
	if err != nil {
		base.LoggerIns.Error("query account item count err:", err, dt.GetName())
		SendRegisterResponse(reply.Sid, proto.EC_SERV_ERR, "")
		return
	}

	rcd, _ := base.NewRecord()
	rcd.SetValue("uid", count + 1)
	rcd.SetValue("account", dt.GetName())
	rcd.SetValue("pwd", dt.GetPwd())
	u := new(common.DBAdd)
	u.Table = "main_user"
	u.Value = rcd
	val := 0
	err = Root.Port.Call(AccountServiceIns.DBServer(), "RpcService.Add", u, &val)
	if err != nil {
		base.LoggerIns.Error("update account err:", err, dt.GetName())
		SendRegisterResponse(reply.Sid, proto.EC_SERV_ERR, "")
		return
	}

	SendRegisterResponse(reply.Sid, proto.EC_OK, "")
}

type LoginHandler struct {
}

func (h *LoginHandler) Handle(reply *common.RpcClientRequest) {
	if reply.Pkt.Args == nil {
		SendLoginResponse(reply.Sid, proto.EC_NO_ARG, "")
		return
	}

	dt := &proto.Login{}
	err := pblib.Unmarshal(reply.Pkt.Args, dt)
	if err != nil {
		SendLoginResponse(reply.Sid, proto.EC_ARG_BAD_FORMAT, "")
		return
	}
	q := new(common.DBQuery)
	q.Sql = "SELECT * FROM MAIN_USER WHERE ACCOUNT=?"
	sl := make([]interface{}, 0)
	q.Args = append(sl, dt.GetName())
	r := new(common.DBQueryResp)
	r.Ok = false
	r.Values = make([]*base.Record, 0)
	err = Root.Port.Call(AccountServiceIns.DBServer(), "RpcService.Query", q, r)
	if err != nil {
		base.LoggerIns.Error("query account err:", err, dt.GetName())
		SendLoginResponse(reply.Sid, proto.EC_SERV_ERR, "")
		return
	}

	if !r.Ok || r.Values == nil || len(r.Values) == 0 {
		SendLoginResponse(reply.Sid, proto.EC_ACCOUNT_NOT_EXISTS, "")
		return
	}

	pwd := ""
	switch r.Values[0].Value("pwd").(type) {
	case string:
		pwd = r.Values[0].Value("pwd").(string)
	case []byte:
		pwd = string(r.Values[0].Value("pwd").([]byte))
	default:
		base.LoggerIns.Error("pwd format can not to string")
	}
	if pwd != dt.GetPwd() {
		SendLoginResponse(reply.Sid, proto.EC_PWD_WRONG, "")
		return
	}
	SendLoginResponse(reply.Sid, proto.EC_OK, "")
}