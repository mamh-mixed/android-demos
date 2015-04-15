package mongo

import (
	// "bytes"
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
	"gopkg.in/mgo.v2/bson"
	"strconv"
)

// CardBinColl 卡Bin Collection
var CardBinColl = cardBinCollection{"cardBin"}

var tree TrieTree

// buildTree 初始化前缀树
func buildTree() {

	// 加载所有卡bin
	cbs, err := CardBinColl.LoadAll()
	if err != nil {
		log.Panicf("fail to load all CardBin : (%s)", err)
	}

	for _, v := range cbs {
		// 建立前缀树
		tree.build(v.Bin)
	}
}

type cardBinCollection struct {
	name string
}

// Find 根据卡长度查找卡BIN列表
func (c *cardBinCollection) Find(cardNum string) (cb *model.CardBin, err error) {
	cb = new(model.CardBin)
	// q := bson.M{
	// 	"cardLen":  len(cardNum),
	// 	"bin":      bson.M{"$lte": cardNum},
	// 	"overflow": bson.M{"$gt": cardNum},
	// }
	// err = database.C(c.name).Find(q).Sort("-bin", "overflow").Limit(1).One(&cb)
	// if err != nil {
	// 	log.Errorf("Find CardBin ERROR! error message is: %s; condition is: %+v", err.Error(), q)
	// 	return nil, err
	// }

	// 从树中取出卡bin
	cardBin := tree.match(cardNum)
	log.Debugf("cardNum : %s, cardBin : %s", cardNum, cardBin)
	q := bson.M{
		"bin":     cardBin,
		"cardLen": len(cardNum),
	}
	err = database.C(c.name).Find(q).One(cb)

	return
}

// LoadAll 加载所有卡bin
func (c *cardBinCollection) LoadAll() ([]*model.CardBin, error) {
	var cardBins []*model.CardBin
	err := database.C(c.name).Find(nil).All(&cardBins)
	return cardBins, err
}

// node 节点信息
type node struct {
	Keys     []int8
	Children *node
}

// TrieTree 前缀树
type TrieTree struct {
	Root node
}

func (t *TrieTree) build(word string) {

	root := &t.Root
	for i := 0; i < len(word); i++ {

		c, err := strconv.Atoi(string(word[i]))
		if err != nil {
			continue
		}
		keys := root.Keys
		if root.Children == nil {
			root.Children = new(node)
		}
		if keys == nil {
			keys = make([]int8, 10, 10)
		}
		keys[c] = int8(c)
		root.Keys = keys
		root = root.Children
	}
}

func (t *TrieTree) match(cardNum string) string {
	s := ""
	root := &t.Root
	for i := 0; i < len(cardNum); i++ {
		c, _ := strconv.Atoi(string(cardNum[i]))
		if root.Keys == nil {
			break
		}
		key := root.Keys[c]
		if c != 0 && key == 0 {
			break
		}
		// log.Debugf("%d", c)
		s += strconv.Itoa(c)
		root = root.Children
	}
	return s
}
