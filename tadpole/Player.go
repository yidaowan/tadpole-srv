package main

import (
	"regexp"
	"strings"

	"github.com/xiaonanln/goworld"
	"github.com/xiaonanln/goworld/engine/common"
	"github.com/xiaonanln/goworld/engine/consts"
	"github.com/xiaonanln/goworld/engine/entity"
	"github.com/xiaonanln/goworld/engine/gwlog"
)

// Player 对象代表一名玩家
type Player struct {
	entity.Entity
}

func (a *Player) DescribeEntityType(desc *entity.EntityTypeDesc) {
	desc.SetPersistent(true).SetUseAOI(true, 2000)
	desc.DefineAttr("name", "AllClients", "Persistent")
	desc.DefineAttr("lv", "AllClients", "Persistent")
	desc.DefineAttr("hp", "AllClients")
	desc.DefineAttr("hpmax", "AllClients")
	desc.DefineAttr("action", "AllClients")
	desc.DefineAttr("spaceKind", "Persistent")
	desc.DefineAttr("dewdrop", "Persistent")
}

// OnCreated 在Player对象创建后被调用
func (a *Player) OnCreated() {
	a.Entity.OnCreated()
	a.setDefaultAttrs()
}

// setDefaultAttrs 设置玩家的一些默认属性
func (a *Player) setDefaultAttrs() {
	a.Attrs.SetDefaultInt("spaceKind", 1)
	a.Attrs.SetDefaultStr("name", "noname")
	a.Attrs.SetDefaultInt("lv", 1)
	a.Attrs.SetDefaultInt("hp", 100)
	a.Attrs.SetDefaultInt("hpmax", 100)
	a.Attrs.SetDefaultStr("action", "idle")
	a.Attrs.SetDefaultInt("dewdrop", 0)

	a.SetClientSyncing(true)
}

// GetSpaceID 获得玩家的场景ID并发给调用者
func (a *Player) GetSpaceID(callerID common.EntityID) {
	a.Call(callerID, "OnGetPlayerSpaceID", a.ID, a.Space.ID)
}

func (p *Player) enterSpace(spaceKind int) {
	if p.Space.Kind == spaceKind {
		return
	}
	if consts.DEBUG_SPACES {
		gwlog.Infof("%s enter space from %d => %d", p, p.Space.Kind, spaceKind)
	}
	goworld.CallService("SpaceService", "EnterSpace", p.ID, spaceKind)
}

// OnClientConnected is called when client is connected
func (a *Player) OnClientConnected() {
	gwlog.Infof("%s client connected", a)
	a.enterSpace(int(a.GetInt("spaceKind")))
}

// OnClientDisconnected is called when client is lost
func (a *Player) OnClientDisconnected() {
	gwlog.Infof("%s client disconnected", a)
	a.Destroy()
}

// EnterSpace_Client is enter space RPC for client
func (a *Player) EnterSpace_Client(kind int) {
	a.enterSpace(kind)
}

// DoEnterSpace is called by SpaceService to notify avatar entering specified space
func (a *Player) DoEnterSpace(kind int, spaceID common.EntityID) {
	// let the avatar enter space with spaceID
	a.EnterSpace(spaceID, entity.Vector3{})
}

//func (a *Player) randomPosition() entity.Vector3 {
//	minCoord, maxCoord := -400, 400
//	return entity.Vector3{
//		X: entity.Coord(minCoord + rand.Intn(maxCoord-minCoord)),
//		Y: 0,
//		Z: entity.Coord(minCoord + rand.Intn(maxCoord-minCoord)),
//	}
//}

// OnEnterSpace is called when avatar enters a space
func (a *Player) OnEnterSpace() {
	gwlog.Infof("%s ENTER SPACE %s", a, a.Space)
	a.SetClientSyncing(true)
}

var spaceSep *regexp.Regexp = regexp.MustCompile("\\s")

func (a *Player) SendChat_Client(text string) {
	text = strings.TrimSpace(text)
	if text[0] == '/' {
		// this is a command
		cmd := spaceSep.Split(text[1:], -1)
		if cmd[0] == "join" {
			// a.enterRoom(cmd[1])
		} else {
			a.CallClient("ShowError", "无法识别的命令："+cmd[0])
		}
	} else {
		a.CallAllClients("OnRecvChat", a.GetStr("name"), text)
	}
}

func (a *Player) CollectDewdrop_Client(dewdropID common.EntityID) {
	drop := a.Space.GetEntity(dewdropID)
	if drop == nil {
		gwlog.Warnf("CollectDewdrop %s, but dewdrop not found", dewdropID)
		return
	}

	dewdrop := drop.I.(*Dewdrop)
	dewdrop.Collect()
	a.Attrs.SetInt("dewdrop", dewdrop.getAmount())
}

func (a *Player) SetAction_Client(action string) {
	if a.GetInt("hp") <= 0 { // dead already
		return
	}
	a.Attrs.SetStr("action", action)
}

func (a *Player) ShootMiss_Client() {
	a.CallAllClients("Shoot")
}

func (a *Player) ShootHit_Client(victimID common.EntityID) {
	a.CallAllClients("Shoot")
	victim := a.Space.GetEntity(victimID)
	if victim == nil {
		gwlog.Warnf("Shoot %s, but monster not found", victimID)
		return
	}

	if victim.Attrs.GetInt("hp") <= 0 {
		return
	}

	monster := victim.I.(*Monster)
	monster.TakeDamage(50)
}

func (player *Player) TakeDamage(damage int64) {
	hp := player.GetInt("hp")
	if hp <= 0 {
		return
	}

	hp = hp - damage
	if hp < 0 {
		hp = 0
	}

	player.Attrs.SetInt("hp", hp)

	if hp <= 0 {
		// now player dead ...
		player.Attrs.SetStr("action", "death")
		player.SetClientSyncing(false)
	}
}
