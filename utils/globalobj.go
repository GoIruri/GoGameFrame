package utils

import (
	"encoding/json"
	"io/ioutil"
	"zinx/giface"
)

// GlobalObj 存储一切关于game框架的参数，一些参数可以通过game.json由用户配置
type GlobalObj struct {
	// Server
	TcpServer giface.Iserver //全局Server对象
	Host      string
	TcpPort   int
	Name      string
	// Game
	Version        string
	MaxConn        int    //当前服务器主机允许的最大链接数
	MaxPackageSize uint32 //当前框架数据包的最大值
}

// GlobalObject 定义一个全局的对外Global对象
var GlobalObject *GlobalObj

func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("mydemo/gamev0.4/conf/game.json")
	if err != nil {
		panic(err)
	}
	// 将json文件数据解析到struct中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

func init() {
	// 如果配置文件没有加载，默认的值
	GlobalObject = &GlobalObj{
		Name:           "GameServerApp",
		Version:        "V0.4",
		TcpPort:        8999,
		Host:           "0.0.0.0",
		MaxConn:        1000,
		MaxPackageSize: 4096,
	}

	// 尝试从utils/game.json中加载一些用户自定义的参数
	GlobalObject.Reload()
}
