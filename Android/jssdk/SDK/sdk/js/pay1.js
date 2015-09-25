var SDK=function(){
var startPay=function(){
 var code=Util._getUrlParam("code");
    if(code){
    window.localStorage.setItem("vt_code", code);
    return;
   }

	var orderData = {
                orderNumber: Util._getUrlParam("orderNumber"),
                orderAmount: Util._getUrlParam("orderAmount"),
                orderCurrency: Util._getUrlParam("orderCurrency"),
                backUrl: Util._getUrlParam("backUrl"),
                frontUrl: Util._getUrlParam("frontUrl"),
                merID: Util._getUrlParam("merID"),
                inscd:Util._getUrlParam("inscd"),
                appId: Util._getUrlParam("appId"),
                secretKey: Util._getUrlParam("secretKey"),
                goodsInfo:utf8to16(base64decode(Util._getUrlParam("goodsInfo"))),
                attach:utf8to16(base64decode(Util._getUrlParam("attach"))),
        };

      
         var type=Util._getUrlParam("type");
         if(type==="wpay"){

         	 yunshouyin.startWPay(orderData);
         }else if(type==="onekeypay"){
         	 yunshouyin.startOneKeyPay(orderData);
         }
        
}

var pay=function(){

  // var paid =window.localStorage.getItem("paid");
  //   if(paid){
  //    return;
  //  }
  //  window.localStorage.setItem("paid","true");

    var code=Util._getUrlParam("code");
    var orderdata=Util._getUrlParam("data");
    orderdata=utf8to16(base64decode(orderdata));
     if(getQueryString(orderdata,"goodsInfo")){
       orderdata=replaceParamVal(orderdata,"goodsInfo",utf8to16(base64decode(getQueryString(orderdata,"goodsInfo"))));
    }

    if(getQueryString(orderdata,"attach")){
      orderdata=replaceParamVal(orderdata,"attach",utf8to16(base64decode(getQueryString(orderdata,"attach"))));
    }
    
    var errorData=
     "attach="+getQueryString(orderdata,"attach")+                    
    "&goodsInfo="+getQueryString(orderdata,"goodsInfo")+
    "&orderAmount="+Util.getNormalTxamt(getQueryString(orderdata,"txamt"))+
    "&orderCurrency="+getQueryString(orderdata,"orderCurrency")+
    "&orderNum="+getQueryString(orderdata,"orderNum")+
    "&state=1";
 
    var frontUrl = getQueryString(orderdata,"frontUrl");

    if (!code) {
      errorData=errorData+"&errorDetail=URL拼接错误"
      window.location.replace(frontUrl + "?" + errorData);
    }

    var postData = {
      orderNum: getQueryString(orderdata,"orderNum"),
      txamt: getQueryString(orderdata,"txamt"),
      backUrl:getQueryString(orderdata,"backUrl"),
      mchntid: getQueryString(orderdata,"mchntid"),
      inscd: getQueryString(orderdata,"inscd"),
      txndir: "Q",
      goodsInfo: getQueryString(orderdata,"goodsInfo"),
      chcd: "WXP",
      busicd: "JSZF",
      needUserInfo: "NO",
      code: code,
      attach: getQueryString(orderdata,"attach"),
      currency: getQueryString(orderdata,"orderCurrency")
    };



    var url = "http://test.quick.ipay.so/scanpay/unified";


    $.ajax({
      type: 'POST',
      url: url,
      async: false,
      data: JSON.stringify(postData),
      dataType: 'json',
      success: function(data) {

        var json = data;
        var jsonobj = json.payjson;
        if (!jsonobj) {
          errorData=errorData+"&errorDetail="+json.errorDetail;
          window.location.replace(frontUrl + "?" + errorData);
          return;

        }


        var dataJson1 = jsonobj.config;
        var dataJson2 = jsonobj.chooseWXPay;
        var config_appId = dataJson1.appId;
        var config_timestamp = dataJson1.timestamp;
        var config_nonceStr = dataJson1.nonceStr;
        var confg_signature = dataJson1.signature;


        var chooseWXPay_timestamp = dataJson2.timeStamp;
        var chooseWXPay_nonceStr = dataJson2.nonceStr;
        var chooseWXPay_package = dataJson2.package;
        var chooseWXPay_signType = dataJson2.signType;
        var chooseWXPay_paySign = dataJson2.paySign;



       var data=
            "attach="+json.attach+                    
            "&goodsInfo="+getQueryString(orderdata,"goodsInfo")+
            "&orderAmount="+Util.getNormalTxamt(json.txamt)+
            "&orderCurrency="+getQueryString(orderdata,"orderCurrency")+
            "&orderNum="+json.orderNum;



        wx.config({
          debug: false, // 开启调试模式,调用的所有api的返回值会在客户端alert出来，若要查看传入的参数，可以在pc端打开，参数信息会通过log打出，仅在pc端时才会打印。
          appId: config_appId, // 必填，公众号的唯一标识
          timestamp: config_timestamp, // 必填，生成签名的时间戳
          nonceStr: config_nonceStr, // 必填，生成签名的随机串
          signature: confg_signature, // 必填，签名，见附录1
          jsApiList: [
              'checkJsApi',
              'chooseWXPay'
            ] // 必填，需要使用的JS接口列表，所有JS接口列表见附录2           
        });

        wx.ready(function() {
          wx.chooseWXPay({
            timestamp: chooseWXPay_timestamp, // 支付签名时间戳，注意微信jssdk中的所有使用timestamp字段均为小写。但最新版的支付后台生成签名使用的timeStamp字段名需大写其中的S字符
            nonceStr: chooseWXPay_nonceStr, // 支付签名随机串，不长于 32 位
            package: chooseWXPay_package, // 统一支付接口返回的prepay_id参数值，提交格式如：prepay_id=***）
            signType: chooseWXPay_signType, // 签名方式，默认为'SHA1'，使用新版支付需传入'MD5'
            paySign: chooseWXPay_paySign, // 支付签名
            success: function(res) {
              data = data + "&state=0";
              window.location.replace(frontUrl + "?" + data);
            },
            fail: function() {
              data = data + "&state=1";
              window.location.replace(frontUrl + "?" + data);
            },
            cancel: function() {
              data = data + "&state=-1";
              window.location.replace(frontUrl + "?" + data);
            }

          });
        });

      },

      error: function(message) {
        alert(JSON.stringify(message));
      }
    });


    
}


var pay1 = function() {
    var code = Util._getUrlParam("code");
    var orderdata = Util._getUrlParam("data");
    orderdata = utf8to16(base64decode(orderdata));
    orderdata = eval('(' + orderdata + ')');
    var errorData =
      "attach=" + orderdata.attach +
      "&txamt=" + orderdata.txamt +
      "&goodsInfo=" + orderdata.goodsInfo +
      "&orderCurrency=" + orderdata.orderCurrency +
      "&orderNum=" + orderdata.orderNum +
      "&state=1";
 
    var frontUrl = orderdata.frontUrl;

    if (!code) {
      errorData=errorData+"&errorDetail=URL拼接错误"
      window.location.replace(frontUrl + "?" + errorData);
    }

    var postData = {
      orderNum: orderdata.orderNum,
      txamt: orderdata.txamt,
      backUrl: orderdata.backUrl,
      mchntid: orderdata.mchntid,
      inscd: orderdata.inscd,
      txndir: "Q",
      goodsInfo: orderdata.goodsInfo,
      chcd: "WXP",
      busicd: "JSZF",
      needUserInfo: "NO",
      code: code,
      attach: orderdata.attach,
      currency: orderdata.orderCurrency
    };



    var url = "http://test.quick.ipay.so/scanpay/unified";


    $.ajax({
      type: 'POST',
      url: url,
      async: false,
      data: JSON.stringify(postData),
      dataType: 'json',
      success: function(data) {

        var json = data;
        var jsonobj = json.payjson;
        if (!jsonobj) {
          errorData=errorData+"&errorDetail="+json.errorDetail;
          window.location.replace(frontUrl + "?" + errorData);
          return;

        }


        var dataJson1 = jsonobj.config;
        var dataJson2 = jsonobj.chooseWXPay;
        var config_appId = dataJson1.appId;
        var config_timestamp = dataJson1.timestamp;
        var config_nonceStr = dataJson1.nonceStr;
        var confg_signature = dataJson1.signature;


        var chooseWXPay_timestamp = dataJson2.timeStamp;
        var chooseWXPay_nonceStr = dataJson2.nonceStr;
        var chooseWXPay_package = dataJson2.package;
        var chooseWXPay_signType = dataJson2.signType;
        var chooseWXPay_paySign = dataJson2.paySign;



        var data =
          "attach=" + json.attach +
          "&txamt=" + json.txamt +
          "&goodsInfo=" + orderdata.goodsInfo +
          "&orderCurrency=" + orderdata.orderCurrency +
          "&orderNum=" + json.orderNum;



        wx.config({
          debug: false, // 开启调试模式,调用的所有api的返回值会在客户端alert出来，若要查看传入的参数，可以在pc端打开，参数信息会通过log打出，仅在pc端时才会打印。
          appId: config_appId, // 必填，公众号的唯一标识
          timestamp: config_timestamp, // 必填，生成签名的时间戳
          nonceStr: config_nonceStr, // 必填，生成签名的随机串
          signature: confg_signature, // 必填，签名，见附录1
          jsApiList: [
              'checkJsApi',
              'chooseWXPay'
            ] // 必填，需要使用的JS接口列表，所有JS接口列表见附录2           
        });

        wx.ready(function() {
          wx.chooseWXPay({
            timestamp: chooseWXPay_timestamp, // 支付签名时间戳，注意微信jssdk中的所有使用timestamp字段均为小写。但最新版的支付后台生成签名使用的timeStamp字段名需大写其中的S字符
            nonceStr: chooseWXPay_nonceStr, // 支付签名随机串，不长于 32 位
            package: chooseWXPay_package, // 统一支付接口返回的prepay_id参数值，提交格式如：prepay_id=***）
            signType: chooseWXPay_signType, // 签名方式，默认为'SHA1'，使用新版支付需传入'MD5'
            paySign: chooseWXPay_paySign, // 支付签名
            success: function(res) {
              data = data + "&state=0";
              window.location.replace(frontUrl + "?" + data);
            },
            fail: function() {
              data = data + "&state=1";
              window.location.replace(frontUrl + "?" + data);
            },
            cancel: function() {
              data = data + "&state=-1";
              window.location.replace(frontUrl + "?" + data);
            }

          });
        });

      },

      error: function(message) {
        alert(JSON.stringify(message));
      }
    });


  }



var initPayOrder=function(){
    var orderdata=Util._getUrlParam("data");
    orderdata=utf8to16(base64decode(orderdata));
    if(getQueryString(orderdata,"goodsInfo")){
       orderdata=replaceParamVal(orderdata,"goodsInfo",utf8to16(base64decode(getQueryString(orderdata,"goodsInfo"))));
    }

    if(getQueryString(orderdata,"attach")){
      orderdata=replaceParamVal(orderdata,"attach",utf8to16(base64decode(getQueryString(orderdata,"attach"))));
    }

var moneytype;
var orderCurrency=getQueryString(orderdata,"orderCurrency");
var orderAmount=Util.getNormalTxamt(getQueryString(orderdata,"txamt"));
var goodsInfo=getQueryString(orderdata,"goodsInfo");
if(orderCurrency=="CNY"){
  moneytype="￥";
}else if(orderCurrency=="HKD"){
  moneytype="HK$"
}else if(orderCurrency=="USD"){
 moneytype="$";
}

document.getElementById("paymoney").innerHTML=moneytype+orderAmount;
document.getElementById("goodsinfo").innerHTML=goodsInfo;
 

}



 var getQueryString = function(url, name) {
    var reg = new RegExp("(^|&)" + name + "=([^&]*)(&|$)", "i");
    var r = url.match(reg);
    if (r != null) return unescape(r[2]);
    return null;
  };
  

 
  function replaceParamVal(data,paramName,replaceWith) {
    var re=eval('/('+ paramName+'=)([^&]*)/gi');
    var data = data.replace(re,paramName+'='+replaceWith);
    return data;
}

 function delQueStr(url, ref) {
            var str = url;       
            var arr = "";
            var returnurl = "";
            var setparam = "";
            if (str.indexOf('&') != -1) {
                arr = str.split('&');            
                for (i in arr) {
                    if (arr[i].split('=')[0] != ref) {
                        returnurl = returnurl + arr[i].split('=')[0] + "=" + arr[i].split('=')[1] + "&";
                    }
                }
                return  returnurl.substr(0, returnurl.length - 1);
            }
            else {
                arr = str.split('=');
                if (arr[0] == ref) {
                    return url.substr(0, url.indexOf('?'));
                }
                else {
                    return url;
                }
            }
        }







return{
	startPay:startPay,
    pay1: pay1,
    pay:pay,
    initPayOrder:initPayOrder
}

}();






