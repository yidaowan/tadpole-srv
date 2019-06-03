package main

import (
	"github.com/xiaonanln/goworld"
)

var (
	_SERVICE_NAMES = []string{
		"OnlineService",
		"SpaceService",
	}
)

func main() {
	goworld.RegisterSpace(&MySpace{}) // 注册自定义的Space类型

	goworld.RegisterService("OnlineService", &OnlineService{})
	goworld.RegisterService("SpaceService", &SpaceService{})
	// 注册Account类型
	goworld.RegisterEntity("Account", &Account{})
	// 注册Monster类型
	// goworld.RegisterEntity("Monster", &Monster{})
	goworld.RegisterEntity("Dewdrop", &Dewdrop{})
	// 注册Avatar类型，并定义属性
	goworld.RegisterEntity("Player", &Player{})
	// 运行游戏服务器
	goworld.Run()
}
