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

type NewSessionListener interface {
	HandleNewSession(ssn *Session)
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
	listener.Owner.AddNewSession(sess)
}

type SessionMgr struct {
	lock sync.Mutex
	sessionMap map[string] *Session

	ssnConnListener session_connectlister

	newSsnListeners []NewSessionListener
	newssnlock sync.Mutex
	newSsns []*Session
}

func NewSessionMgr() (mgr *SessionMgr, err error) {
	mgr = new(SessionMgr)
	err = mgr.init()
	return
}

func (mgr *SessionMgr) init() error {
	mgr.sessionMap = make(map[string] *Session)
	mgr.ssnConnListener = session_connectlister{Owner:mgr}
	mgr.newSsnListeners = make([]NewSessionListener, 0)
	mgr.newSsns = make([]*Session, 0)
	return nil
}

func (mgr *SessionMgr) Start() {
	Root.Acceptor.RegisterConnListener(mgr.ssnConnListener)
}

func (mgr *SessionMgr) Update() {
	mgr.handle_newsession()

	for _, ssn := range(mgr.sessionMap) {
		ssn.Update()
	}
}

func (mgr *SessionMgr) Close() {
	Root.Acceptor.UnRegisterConnListener(mgr.ssnConnListener)

	mgr.lock.Lock()
	defer mgr.lock.Unlock()
	for _, v := range(mgr.sessionMap) {
		v.Close()
	}
}

func (mgr *SessionMgr) RegisterNewSessionListener(l NewSessionListener) {
	for _, v := range(mgr.newSsnListeners) {
		if v == l {
			return
		}
	}
	mgr.newSsnListeners = append(mgr.newSsnListeners, l)
}

func (mgr *SessionMgr) UnRegisterNewSessionListener(l NewSessionListener) {
	for i, v := range(mgr.newSsnListeners) {
		if v == l {
			mgr.newSsnListeners = append(mgr.newSsnListeners[:i],
				mgr.newSsnListeners[i+1:]...)
			break
		}
	}
}

func (mgr *SessionMgr) Session(sid string) (sess *Session, ok bool) {
	sess, ok = mgr.sessionMap[sid]
	return
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

func (mgr *SessionMgr) AddNewSession(sess *Session) {
	if sess == nil {
		return
	}
	
	mgr.newssnlock.Lock()
	defer mgr.newssnlock.Unlock()
	mgr.newSsns = append(mgr.newSsns, sess)
}

func (mgr *SessionMgr) handle_newsession() {
	if len(mgr.newSsns) == 0 {
		return
	}

	mgr.newssnlock.Lock()
	defer mgr.newssnlock.Unlock()

	if len(mgr.newSsnListeners) == 0 {
		for _, ssn := range(mgr.newSsns) {
			ssn.Close()
		}
	} else {
		for _, ssn := range(mgr.newSsns) {
			for _, l := range(mgr.newSsnListeners) {
				l.HandleNewSession(ssn)
			}
			mgr.AddSession(ssn)
			ssn.Start()
		}
	}

	mgr.newSsns = make([]*Session, 0)
}