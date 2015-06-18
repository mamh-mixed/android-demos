package weixin

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/CardInfoLink/quickpay/tools"
)

// mockPostForm 模拟请求
func mockPostForm(url string, data url.Values) (*http.Response, error) {

	// default return value
	resp := new(http.Response)

	// get params
	service := data.Get("service")
	charSet := data.Get("_input_charset")
	partner := data.Get("partner")
	sign := data.Get("sign")
	// currency := data.Get("currency")
	signType := data.Get("sign_type")
	orderNum := data.Get("out_trade_no")
	subject := data.Get("subject")
	totalFee := data.Get("total_fee")
	productCode := data.Get("product_code")

	// 返回参数
	alpResp := &alpResponse{IsSuccess: "T"}
	alipay := alpDetail{}

	switch service {
	case "alipay.acquire.createandpay":
		// check must value
		if charSet == "" || partner == "" || sign == "" || signType == "" || orderNum == "" ||
			subject == "" || productCode == "" || totalFee == "" {

			alipay.ResultCode = "ORDER_FAIL"
			alipay.DetailErrorCode = "INVALID_PARAMETER"
			alipay.DetailErrorDes = "参数无效"
			break
		}
		// 一定规则随机返回
		txamt, _ := strconv.ParseFloat(totalFee, 64)
		if txamt < 1 {
			alipay.ResultCode = "ORDER_SUCCESS_PAY_SUCCESS"
			alipay.TradeNo = tools.Millisecond()
		} else if txamt >= 1 && txamt < 10 {
			alipay.ResultCode = "ORDER_SUCCESS_PAY_INPROCESS"
			alipay.TradeNo = tools.Millisecond()
		} else if txamt >= 10 && txamt < 100 {
			alipay.ResultCode = "ORDER_SUCCESS_PAY_FAIL"
			alipay.TradeNo = tools.Millisecond()
		} else if txamt >= 100 && txamt < 1000 {
			alipay.ResultCode = "UNKNOWN"
		} else {
			alipay.ResultCode = "ORDER_FAIL"
		}
		alipay.BuyerLogonId = "156****3236"
		alipay.BuyerUserId = "2088212959731883"

	case "alipay.acquire.precreate":

	case "alipay.acquire.refund":

	case "alipay.acquire.query":
		if charSet == "" || partner == "" || sign == "" || signType == "" || orderNum == "" {

			alipay.ResultCode = "FAIL"
			alipay.DetailErrorCode = "INVALID_PARAMETER"
			alipay.DetailErrorDes = "参数无效"
			break
		}
		alipay.ResultCode = "SUCCESS"
		alipay.TradeStatus = "TRADE_SUCCESS"
		alipay.TradeNo = tools.Millisecond()
		alipay.BuyerLogonId = "156****3236"
		alipay.BuyerUserId = "2088212959731883"

	case "alipay.acquire.cancel":

	default:
		alpResp.IsSuccess = "F"
		alpResp.Error = "ILLEGAL_SERVICE"
	}
	// 结果
	alpResp.Response.Alipay = alipay

	b, err := xml.Marshal(alpResp)
	rc := ioutil.NopCloser(bytes.NewReader(b))
	resp.Body = rc

	return resp, err
}
