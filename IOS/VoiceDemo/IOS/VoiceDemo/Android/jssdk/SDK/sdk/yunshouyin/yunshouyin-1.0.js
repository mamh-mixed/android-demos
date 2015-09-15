/**
 * Created by bensonzhang on 15/2/9.
 */
'use strict';
var yunshouyin = function() {
  var XMLHttpReq, resultURL;
  
var startWPay = function(orderData) {
    var code = window.localStorage.getItem("vt_code");
    if (!code) {
      var orderDataStr = JSON.stringify(orderData);
      window.localStorage.setItem("vt_orderdata", orderDataStr);
      var redirect_uri=location.href.split('?')[0];
      var wx_url="https://open.weixin.qq.com/connect/oauth2/authorize?appid=wx25ac886b6dac7dd2"+"&redirect_uri=" + encodeURI(redirect_uri) + "&response_type=code&scope=snsapi_base&state=123#wechat_redirect";
      window.location.href=wx_url;
    } else {
      window.localStorage.removeItem("vt_code");
      var now = new Date();
      var location_href=location.href.split('#')[0];
      var nowStr = now.Format("YYYYMMddHHmmss");
      orderData.transTime = nowStr;
      var encryptionStr ;
      if(orderData.attach!==undefined&&orderData.attach.length>0){
       encryptionStr= "attach="+orderData.attach+
       "&backUrl="+orderData.backUrl;
      }else{
         encryptionStr= "backUrl="+orderData.backUrl;
      }
      encryptionStr=encryptionStr+
      
      "&busicd=PURC"+
      "&chcd=WXP"+
      "&code="+code;
      if(orderData.goodsInfo!==undefined&&orderData.goodsInfo.length>0){
        encryptionStr=encryptionStr+
        "&goodsInfo="+orderData.goodsInfo;
      }
      encryptionStr=encryptionStr+
      "&inscd="+orderData.inscd+
      "&mchntid="+orderData.merID+
      "&orderCurrency="+orderData.orderCurrency+
      "&orderNum="+orderData.orderNumber+
      "&txamt="+getTxamt(orderData.orderAmount)+
      "&txndir=Q"+
      orderData.secretKey ;
      var signatureStr  =new Uint8Array(encodeUTF8(encryptionStr));
       signatureStr=sha1(signatureStr);
       signatureStr=Array.prototype.map.call(signatureStr,function(e){
    return (e<16?"0":"")+e.toString(16);
  }).join("");
     // var url = "http://192.168.199.174:8081/pay";
      var url = "http://211.147.72.70:10003/pay";
     var data ;
      if(orderData.attach!==undefined&&orderData.attach.length>0){
       data= "attach="+orderData.attach+
       "&backUrl="+orderData.backUrl;
      }else{
         data= "backUrl="+orderData.backUrl;
      }
      data=data+
      
      "&busicd=PURC"+
      "&chcd=WXP"+
      "&code="+code;
      if(orderData.goodsInfo!==undefined&&orderData.goodsInfo.length>0){
        data=data+
        "&goodsInfo="+orderData.goodsInfo;
      }
      data=data+
      "&inscd="+orderData.inscd+
      "&mchntid="+orderData.merID+
      "&orderCurrency="+orderData.orderCurrency+
      "&orderNum="+orderData.orderNumber+
      "&txamt="+getTxamt(orderData.orderAmount)+
      "&txndir=Q"+
      "&sign="+signatureStr;

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
            var json = eval('(' + text + ')'); 
            var jsonobj=json.payjson;        
            var dataJson1 = jsonobj.config;
            var dataJson2 = jsonobj.chooseWXPay;
            var config_appId = dataJson1.appId;
            var config_timestamp=dataJson1.timestamp;
            var config_nonceStr=dataJson1.nonceStr;
            var confg_signature=dataJson1.signature;


            var chooseWXPay_timestamp=dataJson2.timeStamp;
            var chooseWXPay_nonceStr=dataJson2.nonceStr;
            var chooseWXPay_package=dataJson2.package;
            var chooseWXPay_signType=dataJson2.signType;
            var chooseWXPay_paySign=dataJson2.paySign;
            
            var orderNum=json.orderNum;
            var mchntid=json.mchntid;
            var inscd=json.inscd;
            var goodsinfo=json.goodsinfo;
            var orderAmount=json.txamt;
            var orderCurrency=json.orderCurrency;
            var errorDetail=json.errorDetail;
            var busicd=json.busicd;
            var attach=json.attach;

            var data=
            "attach="+attach+                    
            "&goodsinfo="+goodsinfo+
            "&orderAmount="+orderAmount+
            "&orderCurrency="+orderCurrency+
            "&orderNum="+orderNum;

             


            

            wx.config({
            debug: false, // 开启调试模式,调用的所有api的返回值会在客户端alert出来，若要查看传入的参数，可以在pc端打开，参数信息会通过log打出，仅在pc端时才会打印。
            appId: config_appId, // 必填，公众号的唯一标识
            timestamp: config_timestamp, // 必填，生成签名的时间戳
            nonceStr: config_nonceStr, // 必填，生成签名的随机串
            signature: confg_signature,// 必填，签名，见附录1
            jsApiList: [
            'checkJsApi',
            'chooseWXPay'
            ] // 必填，需要使用的JS接口列表，所有JS接口列表见附录2           
           });

            wx.ready(function(){
              wx.chooseWXPay({
              timestamp: chooseWXPay_timestamp, // 支付签名时间戳，注意微信jssdk中的所有使用timestamp字段均为小写。但最新版的支付后台生成签名使用的timeStamp字段名需大写其中的S字符
              nonceStr:  chooseWXPay_nonceStr, // 支付签名随机串，不长于 32 位
              package:   chooseWXPay_package, // 统一支付接口返回的prepay_id参数值，提交格式如：prepay_id=***）
              signType:  chooseWXPay_signType, // 签名方式，默认为'SHA1'，使用新版支付需传入'MD5'
              paySign:   chooseWXPay_paySign, // 支付签名
               success: function (res) {
               data=data+"&state=0";
               var signStr=data+orderData.secretKey;
               var sign  =new Uint8Array(encodeUTF8(signStr));
               sign=sha1(sign);
              sign=Array.prototype.map.call(sign,function(e){
    return (e<16?"0":"")+e.toString(16);
  }).join("");
              data=data+"&sign="+sign;
              window.location.replace(orderData.frontUrl+"?"+data);
            },
            fail:function(){
              data=data+"&state=1";
              var signStr=data+orderData.secretKey;
               var sign  =new Uint8Array(encodeUTF8(signStr));
               sign=sha1(sign);
              sign=Array.prototype.map.call(sign,function(e){
    return (e<16?"0":"")+e.toString(16);
  }).join("");
              data=data+"&sign="+sign;
              window.location.replace(orderData.frontUrl+"?"+data);
            },
            cancel:function(){
              data=data+"&state=-1";
             var signStr=data+orderData.secretKey;
               var sign  =new Uint8Array(encodeUTF8(signStr));
               sign=sha1(sign);
              sign=Array.prototype.map.call(sign,function(e){
    return (e<16?"0":"")+e.toString(16);
  }).join("");
              data=data+"&sign="+sign;
              window.location.replace(orderData.frontUrl+"?"+data);
            }
});
            });

           

          }
        }
      };
      XMLHttpReq.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
      XMLHttpReq.send(data);

    }

  };

  

  var onload = function() {
    var url = window.location.href;
    url = url.substring(url.indexOf('?') + 1);
    var code = getQueryString("&" + url, "code");
    if (code) {
      window.localStorage.setItem("vt_code", code);
      var orderData = JSON.parse(window.localStorage.getItem("vt_orderdata"));
      startWPay(orderData);
    }

  }



  
 
  
  var startOneKeyPay = function(orderData) {
      var orderDataStr = JSON.stringify(orderData);
      window.localStorage.setItem("vt_onekeyorderdata", orderDataStr);
         window.location.replace("onekey.html");
  };



var getTxamt=function(txamt){
    var str=txamt;
    var i=parseFloat(str);
    var j=i.toFixed(2);
    j=j*100;
    var num=j;
    str=""+num;
    var  k=12-str.length;
    var sum="";
    for(var l=0;l<k;l++){
      sum=sum+"0";
    }
    sum=sum+str;
    return sum;
  }
 


  var getQueryString = function(url, name) {
    var reg = new RegExp("(^|&)" + name + "=([^&]*)(&|$)", "i");
    var r = url.match(reg);
    if (r != null) return unescape(r[2]);
    return null;
  };
  
function sha1(data){
  /**************************************************
  Author：次碳酸钴（admin@web-tinker.com）
  Input：Uint8Array
  Output：Uint8Array
  **************************************************/
  var i,j,t;
  var l=((data.length+8)>>>6<<4)+16,s=new Uint8Array(l<<2);
  s.set(new Uint8Array(data.buffer)),s=new Uint32Array(s.buffer);
  for(t=new DataView(s.buffer),i=0;i<l;i++)s[i]=t.getUint32(i<<2);
  s[data.length>>2]|=0x80<<(24-(data.length&3)*8);
  s[l-1]=data.length<<3;
  var w=[],f=[
    function(){return m[1]&m[2]|~m[1]&m[3];},
    function(){return m[1]^m[2]^m[3];},
    function(){return m[1]&m[2]|m[1]&m[3]|m[2]&m[3];},
    function(){return m[1]^m[2]^m[3];}
  ],rol=function(n,c){return n<<c|n>>>(32-c);},
  k=[1518500249,1859775393,-1894007588,-899497514],
  m=[1732584193,-271733879,null,null,-1009589776];
  m[2]=~m[0],m[3]=~m[1];
  for(i=0;i<s.length;i+=16){
    var o=m.slice(0);
    for(j=0;j<80;j++)
      w[j]=j<16?s[i+j]:rol(w[j-3]^w[j-8]^w[j-14]^w[j-16],1),
      t=rol(m[0],5)+f[j/20|0]()+m[4]+w[j]+k[j/20|0]|0,
      m[1]=rol(m[1],30),m.pop(),m.unshift(t);
    for(j=0;j<5;j++)m[j]=m[j]+o[j]|0;
  };
  t=new DataView(new Uint32Array(m).buffer);
  for(var i=0;i<5;i++)m[i]=t.getUint32(i<<2);
  return new Uint8Array(new Uint32Array(m).buffer);
};


function encodeUTF8(s){
  var i,r=[],c,x;
  for(i=0;i<s.length;i++)
    if((c=s.charCodeAt(i))<0x80)r.push(c);
    else if(c<0x800)r.push(0xC0+(c>>6&0x1F),0x80+(c&0x3F));
    else {
      if((x=c^0xD800)>>10==0) //对四字节UTF-16转换为Unicode
        c=(x<<10)+(s.charCodeAt(++i)^0xDC00)+0x10000,
        r.push(0xF0+(c>>18&0x7),0x80+(c>>12&0x3F));
      else r.push(0xE0+(c>>12&0xF));
      r.push(0x80+(c>>6&0x3F),0x80+(c&0x3F));
    };
  return r;
};




  return {
    startWPay: startWPay,
    startOneKeyPay:startOneKeyPay,
    onload:onload
  }


}();


Date.prototype.Format = function(formatStr) {
  var str = formatStr;
  str = str.replace(/yyyy|YYYY/, this.getFullYear());
  str = str.replace(/yy|YY/, (this.getYear() % 100) > 9 ? (this.getYear() % 100).toString() : '0' + (this.getYear() % 100));
  var month = this.getMonth() + 1;
  str = str.replace(/MM/, month > 9 ? month.toString() : '0' + month);
  str = str.replace(/M/g, month);
  str = str.replace(/dd|DD/, this.getDate() > 9 ? this.getDate().toString() : '0' + this.getDate());
  str = str.replace(/d|D/g, this.getDate());

  str = str.replace(/hh|HH/, this.getHours() > 9 ? this.getHours().toString() : '0' + this.getHours());
  str = str.replace(/h|H/g, this.getHours());
  str = str.replace(/mm/, this.getMinutes() > 9 ? this.getMinutes().toString() : '0' + this.getMinutes());
  str = str.replace(/m/g, this.getMinutes());

  str = str.replace(/ss|SS/, this.getSeconds() > 9 ? this.getSeconds().toString() : '0' + this.getSeconds());
  str = str.replace(/s|S/g, this.getSeconds());
  return str;
}



