var Main=function(){
var merID,inscd;

var init = function() {
	var merchantCode = Util._getUrlParam("merchantCode");
	if(merchantCode===null){
		return;
	}
	var data = "merchantCode=" + merchantCode 
	var url = "http://test.quick.ipay.so/scanpay/fixed/merInfo";
	//var url = "http://211.144.213.118/scanMerchantCode";
	 $.ajax({
      type: 'POST',
      url: url,
      async: true,
      data: data,
      dataType: 'json',
      success: function(data) {
 
         if(data.response=="00"){
           merID=data.merID;
           inscd=data.inscd;
           var title_one =data.title_one;
           var title_two=data.title_two; 
           document.getElementById("title_one").innerHTML="欢迎光临"+title_one;
		   document.getElementById("title_two").innerHTML="本店名称-"+title_two;
		   document.title=title_one;
		   window.localStorage.setItem("title_one", title_one);
         }else{
           alert(data.errorDetail);
           back();
         } 
      },

      error: function(message) {
        alert(JSON.stringify(message));
         back();
      }
    });

}

var back=function(index){
	WeixinJSBridge.call('closeWindow');

}

var getCode=function(){

	  var merchantCode = Util._getUrlParam("merchantCode");
	  var redirect_uri = "http://qrcode.cardinfolink.net/agent/pay1.html?merchantCode="+merchantCode+"&showwxpaytitle=1";
      window.location.href="https://open.weixin.qq.com/connect/oauth2/authorize?appid=wx25ac886b6dac7dd2&redirect_uri=" + encodeURIComponent(redirect_uri) + "&response_type=code&scope=snsapi_userinfo&state=123#wechat_redirect";
}

var wpay = function() {
	var code = Util._getUrlParam("code");
	var money = document.getElementById('money').value;
	if (check(money)) {
		var orderAmount = parseFloat(money);
		orderAmount = orderAmount.toFixed(2);
		var currency = "CNY"
		now = new Date();
		var orderNum = now.Format("YYMMddHHmmss");
		var veriCode=""
		for (var i = 0; i < 5; i++) {
			orderNum = orderNum + Math.floor(Math.random() * 10);
		}

		for (var j = 0; j < 4; j++) {
			veriCode = veriCode + Math.floor(Math.random() * 10);
		}

		var orderData = {
			orderNum: orderNum,
			txamt: "" + Util.getTxamt(orderAmount),
			orderCurrency: currency,
			backUrl: "http://wx.vtpayment.com/wapdemo/agent/onekeypay/result",
			frontUrl: "payresult1.html",
			mchntid: merID,
			inscd: inscd,
			goodsInfo: "云收银wap客户端",
			attach: "用户附加数据原样返回",
			veriCode:veriCode

		};

		yunshouyin.startWPay1(orderData,code);
	}


}


function check(v) {
	if (v.length === 0) {
		alert("金额不能为空");
		return false;
	}
	var a = /^[0-9]*(\.[0-9]{1,2})?$/;
	if (!a.test(v)) {
		alert("金额不正确");
		return false;
	} else {
		return true;
	}
}

var getResult=function(){
	var state=Util._getUrlParam("state");
	if(state==="0"){
        window.location.replace("result.html?code="+Util._getUrlParam("code")+"&orderAmount="+Util.getNormalTxamt(Util._getUrlParam("txamt")));
	}else{
		window.location.replace("fail.html");
	}
}

var initResult=function(){
  
	var code=Util._getUrlParam("code");
	var orderAmount=Util._getUrlParam("orderAmount");
	document.getElementById("code").innerHTML=code;
	document.getElementById("amount").innerHTML="付款金额:¥"+orderAmount;
    document.title=window.localStorage.getItem("title_one");
}

var initFailResult=function(){
	document.title=window.localStorage.getItem("title_one");
}


return{
	init:init,
	wpay:wpay,
	getResult:getResult,
	back:back,
	initResult:initResult,
	initFailResult:initFailResult,
	getCode:getCode
}

}();