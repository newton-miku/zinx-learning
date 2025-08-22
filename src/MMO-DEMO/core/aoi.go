package core

import "fmt"

type AOIManager struct {
	// 区域左边界坐标
	MinX int
	// 区域右边界坐标
	MaxX int
	// 区域上边界坐标
	MinY int
	// 区域下边界坐标
	MaxY int
	// x轴上的格子数量
	Cntx int
	// y轴上的格子数量
	Cnty int
	// 区域内的格子序号，Key为GID，Value为格子对象
	grids map[int]*Grid
}

func NewAOIManager(minx, maxx, cntx, miny, maxy, cnty int) *AOIManager {
	aoiMgr := &AOIManager{
		MinX:  minx,
		MinY:  miny,
		MaxX:  maxx,
		MaxY:  maxy,
		Cntx:  cntx,
		Cnty:  cnty,
		grids: make(map[int]*Grid),
	}

	// 对所有格子进行初始化
	for y := 0; y < cnty; y++ {
		for x := 0; x < cntx; x++ {
			// 计算gid
			gid := y*cntx + x

			aoiMgr.grids[gid] = NewGrid(gid,
				aoiMgr.MinX+x*aoiMgr.getWidth(),
				aoiMgr.MinY+y*aoiMgr.getHeight(),
				aoiMgr.MinX+(x+1)*aoiMgr.getWidth(),
				aoiMgr.MinY+(y+1)*aoiMgr.getHeight(),
				x,
				y,
			)
		}
	}

	return aoiMgr
}

func (m *AOIManager) getWidth() int {
	return (m.MaxX - m.MinX) / m.Cntx
}

func (m *AOIManager) getHeight() int {
	return (m.MaxY - m.MinY) / m.Cnty
}

func (m *AOIManager) String() string {
	str := fmt.Sprintf("AOIManager: \nMinx:%d, MinY %d, MaxX %d, MaxY %d, Cntx %d, Cnty %d\n",
		m.MinX, m.MinY, m.MaxX, m.MaxY, m.Cntx, m.Cnty)
	for i, _ := range m.grids {
		str += fmt.Sprint(m.grids[i])
	}
	return str
}

func (m *AOIManager) GetSideGridsByGid(gid int) (resp []*Grid) {
	if _, ok := m.grids[gid]; !ok {
		return
	}
	// 将当前gid放入结果集中
	resp = append(resp, m.grids[gid])

	// 判断左边是否存在格子
	if m.grids[gid].XId > 0 {
		resp = append(resp, m.grids[gid-1])
	}
	// 判断右边是否存在格子
	if m.grids[gid].XId < m.Cntx-1 {
		resp = append(resp, m.grids[gid+1])
	}
	return
}

func (m *AOIManager) GetSurroundGridsByGid(gid int) (resp []*Grid) {
	// 先判断是否是aoi内的格子
	if _, ok := m.grids[gid]; !ok {
		return
	}

	thisGrid := m.GetSideGridsByGid(gid)
	resp = append(resp, thisGrid...)

	// 获取上面格子gid
	if upGridGid := gid - m.Cntx; upGridGid >= 0 {
		upGrid := m.GetSideGridsByGid(upGridGid)
		resp = append(resp, upGrid...)
	}
	// 获取下面格子gids
	if downGridGid := gid + m.Cntx; downGridGid < m.Cntx*m.Cnty {
		downGrid := m.GetSideGridsByGid(downGridGid)
		resp = append(resp, downGrid...)
	}
	return
}

// 通过坐标获取格子gid
func (m *AOIManager) GetGidByPos(x, y float32) (gid int) {
	xid := (int(x) - m.MinX) / m.getWidth()
	yid := (int(y) - m.MinY) / m.getHeight()

	gid = yid*m.Cntx + xid

	return
}

// 通过坐标获取周围的玩家id
func (m *AOIManager) GetPidsByPos(x, y float32) (respIDs []int) {
	gid := m.GetGidByPos(x, y)

	for _, grid := range m.GetSurroundGridsByGid(gid) {
		respIDs = append(respIDs, grid.GetPlayers()...)
	}

	return
}

func (m *AOIManager) AddPidToGrid(gid, pid int) {
	m.grids[gid].AddPlayer(pid)
}

func (m *AOIManager) RemovePidFromGrid(gid, pid int) {
	m.grids[gid].RemovePlayer(pid)
}

func (m *AOIManager) GetPidsByGid(gid int) []int {
	return m.grids[gid].GetPlayers()
}

func (m *AOIManager) AddPlayerByPos(x, y float32, pid int) {
	gid := m.GetGidByPos(x, y)
	m.AddPidToGrid(gid, pid)
}

func (m *AOIManager) RemovePlayerByPos(x, y float32, pid int) {
	gid := m.GetGidByPos(x, y)
	m.RemovePidFromGrid(gid, pid)
}
