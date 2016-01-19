package domestic

import (
	"encoding/xml"
	"testing"

	"github.com/CardInfoLink/log"
)

var data = `
<?xml version="1.0" encoding="utf-8"?>
<alipay>
   <is_success>T</is_success>
   <request>
		<param name="body">声波支付-分账-sky</param>
		<param name="operator_id">55</param>
		<param name="subject">声波支付-分账-sky</param>
		<param name="sign_type">MD5</param>
		<param name="out_trade_no">7085502131376415</param>
		<param name="dynamic_id">kfylrwezsbeqhh553e</param>
		<param name="royalty_parameters">
		[{"serialNo":"1","transOut":"2088101126765726","transIn":"208810112670840 2","amount":"1.00","desc":"分账测试 1"}]
		</param>
        <param name="royalty_type">ROYALTY</param>
        <param name="total_fee">10</param>
        <param name="partner">2088101106499364</param>
        <param name="quantity">10</param>
		<param name="dynamic_id_type">soundwave</param>
		<param name="alipay_ca_request">2</param>
		<param name="sign">a1cb41a4019351965d4418c9cb933f0f</param> 
		<param name="_input_charset">UTF-8</param>
		<param name="price">1</param>
		<param name="it_b_pay">1d</param>
		<param name="product_code">SOUNDWAVE_PAY_OFFLINE</param> 
		<param name="service">alipay.acquire.createandpay</param> 
		<param name="seller_id">2088101106499364</param>
   </request>
   <response>
		<alipay>
			<buyer_logon_id>138****0011</buyer_logon_id>
			<buyer_user_id>2088102105236945</buyer_user_id>
			<out_trade_no>7085502131376415</out_trade_no>
			<result_code>ORDER_SUCCESS_PAY_SUCCESS</result_code>
			<trade_no>2013112311001004940000384027</trade_no>
			<fund_bill_list>
				<TradeFundBill>
			       <amount>70.00</amount>
			       <fund_channel>10</fund_channel>
			   </TradeFundBill>
			   <TradeFundBill>
			       <amount>20.00</amount>
			       <fund_channel>00</fund_channel>
			   </TradeFundBill>
			   <TradeFundBill>
			       <amount>10.00</amount>
			       <fund_channel>30</fund_channel>
			   </TradeFundBill>
			</fund_bill_list>
		</alipay>
   </response>
   <sign>ea489fc31da63253bab52ed77fb45eb7</sign>
   <sign_type>MD5</sign_type>
</alipay>`

// TestDecodeXml test if unmarshal response success
func TestDecodeXml(t *testing.T) {

	v := &alpResponse{}
	err := xml.Unmarshal([]byte(data), v)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	mer, chcd := v.Response.Alipay.DisCount()
	log.Debugf("%+v,%s,%s", v, mer, chcd)
}

func TestEncodeXml(t *testing.T) {
	v := &alpResponse{
		Request: []Param{
			{
				Name:  "operator_id",
				Value: "55",
			},
			{
				Name:  "subject",
				Value: "test",
			},
		},
	}
	b, _ := xml.Marshal(v)
	log.Debug(string(b))
}
