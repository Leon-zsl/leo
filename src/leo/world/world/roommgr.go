/* this is the room manager
*/

package world

type RoomMgr struct {
	room_map map[string] *Room
}

func NewRoomMgr() (mgr *RoomMgr, err error) {
	mgr = new(RoomMgr)
	err = mgr.init()
	return
}

func (mgr *RoomMgr) init() error {
	mgr.room_map = make(map[string] *Room)
	return nil
}

func (mgr *RoomMgr) Start() error {
	//do nothing now
	return nil
}

func (mgr *RoomMgr) Close() error {
	for _, room := range(mgr.room_map) {
		room.Close()
	}
	mgr.room_map = make(map[string]*Room)
	return nil
}

func (mgr *RoomMgr) Tick() error {
	//do nothing now. 
	//one room one goroutine
	return nil
}

func (mgr *RoomMgr) Room(id string) (*Room, bool) {
	rm, ok := mgr.room_map[id]
	return rm, ok
}

func (mgr *RoomMgr) AddRoom(rm *Room) {
	if rm == nil {
		return
	}
	mgr.room_map[rm.ID()] = rm
}

func (mgr *RoomMgr) DelRoom(rm *Room) {
	if rm == nil {
		return
	}
	delete(mgr.room_map, rm.ID())
}

func (mgr *RoomMgr) DelRoomByID(id string) {
	delete(mgr.room_map, id)
}