package core

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"
	"github.com/CardInfoLink/log"
)

var tree trieTree

// node 节点信息
type node struct {
	flag     bool      // 标识该节点是否是结束标识
	children [10]*node // 孩子节点
}

// TrieTree 前缀树
type trieTree struct {
	root  node
	mutex sync.RWMutex
}

// findCardBin 根据卡号找卡bin
func findCardBin(cardNum string) (*model.CardBin, error) {

	bin := tree.match(cardNum)
	if bin == "" {
		return nil, fmt.Errorf("no bin match cardNum(%s)", cardNum)
	}
	log.Debugf("cardNum=%s, cardBin=%s", cardNum, bin)
	return mongo.CardBinColl.Find(bin, len(cardNum))
}

// ReBuildTree 重新初始化树
// 遇到错误，直接返回
func ReBuildTree() error {

	temp := trieTree{}
	cbs, err := mongo.CardBinColl.LoadAll()
	if err != nil {
		return fmt.Errorf("fail to load all CardBin when rebuilding the tree: (%s)", err)
	}
	for _, v := range cbs {
		// 建立前缀树
		err := temp.build(v.Bin)
		if err != nil {
			return fmt.Errorf("fail to build cardBin tree with the given Bin(%s): %s", v.Bin, err)
		}
	}

	// 测试新建的树是否能正确匹配
	s := temp.match("6222022003008481261")
	if s != "622202" {
		return fmt.Errorf("%s", "the new cardBin tree does not work correctly")
	}

	// 加上写锁，改变根节点
	tree.mutex.Lock()
	tree.root = temp.root

	//for test death lock
	//tree.match("6222022003008481261")

	tree.mutex.Unlock()
	log.Infof("rebuild cardBin tree success %+v", tree)

	return nil
}

// BuildTree 初始化前缀树
func BuildTree() {
	// 加载所有卡bin
	cbs, err := mongo.CardBinColl.LoadAll()
	if err != nil {
		log.Panicf("fail to load all CardBin : %s", err)
	}

	for _, v := range cbs {
		// 建立前缀树
		err := tree.build(v.Bin)
		if err != nil {
			log.Panicf("fail to build cardBin tree with the given Bin(%s): %s", v.Bin, err)
		}
	}
	log.Info("cardBin trieTree init success")
}

// build 建立树
func (t *trieTree) build(word string) error {
	// 根节点开始
	root := &t.root
	for i := 0; i < len(word); i++ {
		index, err := strconv.Atoi(string(word[i]))
		if err != nil {
			return err
		}
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
	return nil
}

func (t *trieTree) match(cardNum string) string {

	s, temp := "", ""

	// 加上读锁，防止写操作
	t.mutex.RLock()
	root := &t.root
	t.mutex.RUnlock()

	for i := 0; i < len(cardNum); i++ {
		index, err := strconv.Atoi(string(cardNum[i]))
		if err != nil {
			return ""
		}
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
