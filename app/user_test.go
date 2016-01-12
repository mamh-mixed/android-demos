package app

import (
	"testing"
	"time"

	"github.com/CardInfoLink/quickpay/model"
)

func TestGetQiniuToken(t *testing.T) {
	result := User.getQiniuToken(&reqParams{
		UserName: "fnghwsj@qq.com",
		Password: "83d90a0f21db74e4cb78d6f2cbccb387",
	})

	t.Logf("result is %+v", result)
}

func TestMonthParse(t *testing.T) {
	m := "201611"

	tm, err := time.ParseInLocation("200601", m, time.Local)
	if err != nil {
		t.Errorf("result is %s", err)
	}

	endTime := tm.AddDate(0, 1, 0).Add(-time.Second)
	t.Logf("time is %s; end time is %s", tm.Format("2006-01-02 15:04:05"), endTime.Format("2006-01-02 15:04:05"))
}

func TestUserV3GetUserBills(t *testing.T) {
	result := UserV3.getUserBills(&reqParams{
		UserName:  "842712881@Qq.com",
		Password:  "e10adc3949ba59abbe56e057f20f883e",
		Transtime: "201512",
		Month:     "201601",
		Status:    "all",
		Index:     "1",
		Size:      "size",
		TransType: model.PayTrans,
	})

	t.Logf("result is %+v", result)
}

func TestUserV3TotalSummaryHandler(t *testing.T) {
	result := UserV3.getDaySummary(&reqParams{
		UserName:     "842712881@Qq.com",
		Password:     "e10adc3949ba59abbe56e057f20f883e",
		Transtime:    "201512",
		Date:         "20160105",
		BusinessType: "1",
	})

	t.Logf("result is %+v", result)
}

func TestUserV3CouponsHandler(t *testing.T) {
	result := User.couponsHandler(&reqParams{
		UserName:  "453481716@qq.com",
		Password:  "670b14728ad9902aecba32e22fa4f6bd",
		Transtime: "201512",
		// Month:        "",
		Index:    "0",
		ClientId: "99911888",
	})

	t.Logf("result is %+v", result)
}

func TestForgetPasswordSendEmail(t *testing.T) {
	result := User.forgetPassword(&reqParams{
		UserName: "fnghwsj@qq.com",
	})

	t.Logf("result is %+v", result)
}
