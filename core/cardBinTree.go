package core

import (
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/omigo/log"
	"strconv"
)

var tree trieTree

// node 节点信息
type node struct {
	flag     bool      // 标识该节点是否是结束标识
	children [10]*node // 孩子节点
}

// TrieTree 前缀树
type trieTree struct {
	root node
}

// buildTree 初始化前缀树
func init() {

	// 加载所有卡bin
	cbs, err := mongo.CardBinColl.LoadAll()
	if err != nil {
		log.Panicf("fail to load all CardBin : (%s)", err)
	}

	for _, v := range cbs {
		// 建立前缀树
		tree.build(v.Bin)
	}
	log.Infof("cardBin trieTree init success %+v", tree)
}

// build 建立树
func (t *trieTree) build(word string) {
	// 根节点开始
	root := &t.root
	for i := 0; i < len(word); i++ {
		index, _ := strconv.Atoi(string(word[i]))
		k := root.children[index]
		if k == nil {
			k = new(node)
		}
		// 结束时加标志位
		if i == len(word)-1 {
			k.flag = true
		}
		root.children[index] = k
		root = k
	}
}

func (t *trieTree) match(cardNum string) string {
	s := ""
	temp := ""
	root := &t.root
	for i := 0; i < len(cardNum); i++ {
		index, _ := strconv.Atoi(string(cardNum[i]))
		k := root.children[index]
		if k == nil {
			break
		}
		s += strconv.Itoa(index)
		// 判断是否为卡bin结束位
		// 是的话赋值给temp
		if k.flag {
			// log.Debugf("%d", index)
			temp = s
		}
		root = root.children[index]
	}

	return temp
}
