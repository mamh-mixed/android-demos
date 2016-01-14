package mongo

import (
	"github.com/CardInfoLink/log"
	"github.com/CardInfoLink/quickpay/cache"
	"github.com/CardInfoLink/quickpay/model"
	"gopkg.in/mgo.v2/bson"
)

var ScanPayRespCol = &scanPayRespCollection{"respCode.sp",
	&model.ScanPayRespCode{"", "", "58", "未知", false, "UNKNOWN"},
}

type scanPayRespCollection struct {
	name        string
	DefaultResp *model.ScanPayRespCode // 未知应答码时，使用渠道应答
}

var spRespCache = cache.New(model.Cache_ScanPayResp)

// Get 根据传入的errorCode类型得到Resp对象
// 屏蔽8583与6位应答码的差别
func (c *scanPayRespCollection) Get(errorCode string) (resp *model.ScanPayRespCode) {

	o, found := spRespCache.Get(errorCode)
	if found {
		resp = o.(*model.ScanPayRespCode)
		return resp
	}

	resp = &model.ScanPayRespCode{}
	err := database.C(c.name).Find(bson.M{"errorCode": errorCode}).One(resp)
	if err != nil {
		log.Errorf("can not find scanPayResp for %s: %s", errorCode, err)
		// 没找到对应应答码，返回默认应答
		return c.DefaultResp
	}

	// save cache
	spRespCache.Set(errorCode, resp, cache.NoExpiration)

	return resp
}

// Get8583CodeAndMsg 8583应答
func (c *scanPayRespCollection) Get8583CodeAndMsg(errorCode string) (code, ch_msg, msg string) {
	spResp := c.Get(errorCode)
	return spResp.ISO8583Code, spResp.ISO8583Msg, errorCode
}

// GetByAlp 由支付宝应答得到Resp对象
func (c *scanPayRespCollection) GetByAlp(code, busicd string) (resp *model.ScanPayRespCode) {
	resp = &model.ScanPayRespCode{}
	q := bson.M{
		"alp": bson.M{
			"$elemMatch": bson.M{
				"code":   code,
				"busicd": busicd,
			},
		},
	}
	err := database.C(c.name).Find(q).One(resp)
	if err != nil {
		log.Errorf("can not find scanPayResp for (code:%s,busicd:%s): %s", code, busicd, err)
		return c.DefaultResp
	}

	return resp
}

// GetByWxp 由微信应答得到Resp对象
func (c *scanPayRespCollection) GetByWxp(code, busicd string) (resp *model.ScanPayRespCode) {
	resp = &model.ScanPayRespCode{}

	q := bson.M{
		"wxp": bson.M{
			"$elemMatch": bson.M{
				"code":   code,
				"busicd": busicd,
			},
		},
	}
	err := database.C(c.name).Find(q).One(resp)
	if err != nil {
		log.Errorf("can not find scanPayResp for (code:%s,busicd:%s): %s", code, busicd, err)
		return c.DefaultResp
	}
	return resp
}

/* only use for import respCode */

func (c *scanPayRespCollection) Add(r *model.ScanPayCSV) error {
	err := database.C(c.name).Insert(r)
	return err
}

func (c *scanPayRespCollection) FindOne(code string) (*model.ScanPayCSV, error) {
	q := new(model.ScanPayCSV)
	err := database.C(c.name).Find(bson.M{"ISO8583Code": code}).One(q)
	return q, err
}

func (c *scanPayRespCollection) Update(r *model.ScanPayCSV) error {
	err := database.C(c.name).Update(bson.M{"ISO8583Code": r.ISO8583Code}, r)
	return err
}
