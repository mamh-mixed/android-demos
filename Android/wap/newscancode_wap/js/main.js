var Main=function(){
var merID,secretKey,inscd;

var init = function() {


	var merchantCode = Util._getUrlParam("merchantCode");
	if(merchantCode===null){
		return;
	}
	now = new Date(),
		transtime = now.Format("YYYYMMddHHmmss");

	var signatureStr = "merchantCode=" + merchantCode +
		"&transtime=" + transtime + "eu1dr0c8znpa43blzy1wirzmk8jqdaon";
	var sign = hex_sha1(signatureStr);

	var data = "merchantCode=" + merchantCode +
		"&transtime=" + transtime +
		"&sign=" + sign;

	var XMLHttpReq, resultURL;
	var url = "http://211.147.72.70:10003/scanMerchantCode";
	//var url = "http://211.144.213.118/scanMerchantCode";
	

	try {
        XMLHttpReq = new ActiveXObject("Msxml2.XMLHTTP");
      } catch (E) {
        try {
          XMLHttpReq = new ActiveXObject("Microsoft.XMLHTTP");
        } catch (E) {
          XMLHttpReq = new XMLHttpRequest();
        }
      }


      XMLHttpReq.open("post", url, false);
      XMLHttpReq.onreadystatechange = function() {
        if (XMLHttpReq.readyState == 4) {
          if (XMLHttpReq.status == 200) {
            var text = XMLHttpReq.responseText;  
            var jsonobj = eval('(' + text + ')'); 
            var merchant=jsonobj.merchant;
            var headMerchant=merchant.headMerchant
            var title_one =headMerchant.title_one;
            var title_two=headMerchant.title_two; 
            merID=merchant.clientid;
            secretKey=merchant.md5;
            inscd=merchant.inscd;
            
               	
            document.getElementById("title_one").innerHTML="欢迎光临"+title_one;
            document.getElementById("title_two").innerHTML="本店名称-"+title_two;

            document.title=title_one;
            window.localStorage.setItem("title_one", title_one);
           
          }
        }
      };

      XMLHttpReq.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
      XMLHttpReq.send(data);



	
}

var back=function(index){
	WeixinJSBridge.call('closeWindow');

}

var getCode=function(){

	  var merchantCode = Util._getUrlParam("merchantCode");
	  var redirect_uri = "http://qrcode.cardinfolink.net/agent/pay.html?merchantCode="+merchantCode+"&showwxpaytitle=1";
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
			orderNumber: orderNum,
			orderAmount: "" + orderAmount,
			orderCurrency: currency,
			backUrl: "http://wx.vtpayment.com/wapdemo/agent/onekeypay/result",
			frontUrl: "payresult.html",
			merID: merID,
			inscd: inscd,
			appId: "wx25ac886b6dac7dd2",
			secretKey: secretKey,
			goodsInfo: "云收银wap客户端",
			attach: "用户附加数据原样返回",
			veriCode:veriCode

		};
		yunshouyin.startWPay(orderData,code);
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
        window.location.replace("result.html?code="+Util._getUrlParam("code")+"&orderAmount="+Util.getNormalTxamt(Util._getUrlParam("orderAmount")));
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