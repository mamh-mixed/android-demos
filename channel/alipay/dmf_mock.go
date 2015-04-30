package alipay

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"
)

// mockPostForm 模拟请求
func mockPostForm(url string, data url.Values) (*http.Response, error) {

	resp := new(http.Response)
	resultCode := "ORDER_SUCCESS_PAY_SUCCESS"

	// TODO check must value
	// service := data.Get("service")
	// charSet := data.Get("_input_charset")
	// partner := data.Get("partner")
	// currency := data.Get("currency")
	// signType := data.Get("sign_type")
	// sign := data.Get("sign")
	// orderNum := data.Get("out_trade_no")
	// subject := data.Get("subject")
	// totalFee := data.Get("total_fee")
	// productCode := data.Get("product_code")
	// TODO random return resultCode

	// 返回参数
	alpResp := &alpResponse{
		IsSuccess: "T",
		Response: alpBody{
			Alipay: alpDetail{
				ResultCode: resultCode,
			},
		},
	}

	b, err := xml.Marshal(alpResp)
	rc := ioutil.NopCloser(bytes.NewReader(b))
	resp.Body = rc

	return resp, err
}
