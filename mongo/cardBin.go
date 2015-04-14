package mongo

import (
	"github.com/CardInfoLink/quickpay/model"
	"github.com/omigo/log"
	"gopkg.in/mgo.v2/bson"
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
	log.Debugf("cardNum : %s,cardBin : %s", cardNum, cardBin)
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
	Key      byte
	Children []*node
}

// TrieTree 前缀树
type TrieTree struct {
	Root node
}

func (t *TrieTree) build(word string) {
	root := &t.Root
	chars := []byte(word)
	for i := 0; i < len(chars); i++ {

		c := chars[i]

		nodes := root.Children
		s, flag := t.isContain(nodes, c)
		// 没包含，添加到子节点
		if !flag {
			s = &node{Key: c}
			nodes = append(nodes, s)
			root.Children = nodes
			root = s
			continue
		}
		// 已包含
		root = s
	}
}

func (t *TrieTree) match(cardNum string) string {
	var result []byte
	// 根节点不包含信息
	var root = t.Root.Children
	chars := []byte(cardNum)
	for i := 0; i < len(chars); i++ {
		c := chars[i]

		n, flag := t.isContain(root, c)
		// 没找到退出循环
		if !flag {
			break
		}
		result = append(result, n.Key)
		// 继续下个节点
		root = n.Children
	}
	return string(result)
}

func (t *TrieTree) isContain(nodes []*node, key byte) (*node, bool) {
	for _, v := range nodes {
		if v.Key == key {
			return v, true
		}
	}
	return nil, false
}
