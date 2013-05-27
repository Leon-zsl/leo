/* this is session manager
*/

package gate

type SessionMgr struct {
	sessionMap map[string] *Session
}

func NewSessionMgr() (mgr *SessionMgr, err error) {
	mgr = nil
	err = nil
	//todo:
	return
}

func (mgr *SessionMgr) Start() {
	//todo:
}

func (mgr *SessionMgr) Update() {
	//todo:
}

func (mgr *SessionMgr) Close() {
	//todo:
}

func (mgr *SessionMgr) AddSession() {
	//todo:
}

func (mgr *SessionMgr) RemoveSession() {
	//todo:
}