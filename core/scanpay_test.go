package core

import (
	"net/url"
	"testing"
)

func TestAlpAsyncNotify(t *testing.T) {

	params := url.Values{}

	params.Add("out_trade_no", "8ef217e208da40cd66f7d86a2c6d476e")
	params.Add("paytools_pay_amount", `[{"MCOUPON":"7.94"},{"TMPOINT":"1.69"},{"MDISCOUNT":" 5.55"}]`)

	AlpAsyncNotify(params)
}
