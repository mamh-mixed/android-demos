// Package test 用于系统测试，而非单元测试。

package test

import (
	"github.com/CardInfoLink/quickpay/channel/cil"
	"github.com/CardInfoLink/quickpay/core"
	"github.com/CardInfoLink/quickpay/mongo"
)

func init() {

	// 连接到 MongoDB
	mongo.Connect()

	// 初始化卡 Bin 树
	core.BuildTree()

	//
	cil.Connect()
}
