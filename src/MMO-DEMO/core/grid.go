package core

import (
	"fmt"
	"sync"
)

type Grid struct {
	GID int
	// 单元格ID
	MinX int
	MinY int
	MaxX int
	MaxY int
	XId  int
	YId  int
	// 格子内玩家或物品的集合
	PlayersID map[int]bool
	pIDLock   sync.RWMutex
}

func NewGrid(gid, minx, miny, maxx, maxy, xid, yid int) *Grid {
	return &Grid{
		GID:       gid,
		MinX:      minx,
		MinY:      miny,
		MaxX:      maxx,
		MaxY:      maxy,
		XId:       xid,
		YId:       yid,
		PlayersID: make(map[int]bool),
	}
}

// 新增玩家
func (g *Grid) AddPlayer(playerID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	g.PlayersID[playerID] = true
}

// 删除玩家
func (g *Grid) RemovePlayer(playerID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	delete(g.PlayersID, playerID)
}

// 获得当前格子内的所有玩家
func (g *Grid) GetPlayers() (respIDs []int) {
	g.pIDLock.RLock()
	defer g.pIDLock.RUnlock()

	for k, _ := range g.PlayersID {
		respIDs = append(respIDs, k)
	}
	return
}

// DEBUG
// 获取格子的基本信息
func (g *Grid) String() string {
	return fmt.Sprintf(
		"GID: %d, MinX: %d, MinY: %d, MaxX: %d, MaxY: %d, PlayersID: %v\n",
		g.GID, g.MinX, g.MinY, g.MaxX, g.MaxY, g.PlayersID)
}
