/* this is session manager
*/

package gate

import (
	"sync"
	"net"
)

type session_connectlister struct {
	Owner *SessionMgr
}

func (listener session_connectlister) Handle(conn *net.TCPConn) {
	if listener.Owner == nil {
		return
	}
	if conn == nil {
		return
	}

	sess, err := NewSession(conn)
	if err != nil {
		Root.Logger.Error("create new session failed: " + conn.RemoteAddr().String())
		return
	}
	listener.Owner.AddSession(sess)
}

type SessionMgr struct {
	lock sync.Mutex
	sessionMap map[string] *Session
	connListener session_connectlister
}

func NewSessionMgr() (mgr *SessionMgr, err error) {
	mgr = new(SessionMgr)
	err = mgr.init()
	return
}

func (mgr *SessionMgr) init() error {
	mgr.sessionMap = make(map[string] *Session)
	mgr.connListener = session_connectlister{Owner:mgr}
	return nil
}

func (mgr *SessionMgr) Start() {
	Root.Acceptor.RegisterConnListener(mgr.connListener)
}

func (mgr *SessionMgr) Update() {
	//todo:
}

func (mgr *SessionMgr) Close() {
	Root.Acceptor.UnRegisterConnListener(mgr.connListener)

	mgr.lock.Lock()
	defer mgr.lock.Unlock()
	for _, v := range(mgr.sessionMap) {
		v.Close()
	}
}

func (mgr *SessionMgr) AddSession(sess *Session) {
	if sess == nil {
		return
	}
	sid := sess.SID()

	mgr.lock.Lock()
	defer mgr.lock.Unlock()
	_, ok := mgr.sessionMap[sid]
	if ok {
		Root.Logger.Warn("duplicate session find: " + sid)
		return
	}
	mgr.sessionMap[sid] = sess
}

func (mgr *SessionMgr) RemoveSession(sess *Session) {
	if sess == nil {
		return
	}
	sid := sess.SID()

	mgr.lock.Lock()
	defer mgr.lock.Unlock()
	delete(mgr.sessionMap, sid)
}