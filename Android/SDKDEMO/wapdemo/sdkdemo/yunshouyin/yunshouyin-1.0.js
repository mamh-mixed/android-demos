/**
 * Created by bensonzhang on 15/2/9.
 */
'use strict';
var yunshouyin = function() {

  var redirect_uri = "http://qrcode.cardinfolink.net/sdk/payment.html";

  var startWPay = function(orderData, debug) {
    if (debug) {
      redirect_uri = "http://qrcode.cardinfolink.net/sdk/payment.html";
    } else {
      redirect_uri = "http://qrcode.cardinfolink.net/jssdk/payment.html";
    }

    pay(orderData);


  };

  var startOneKeyPay = function(orderData, debug) {
    if (debug) {
      redirect_uri = "http://qrcode.cardinfolink.net/sdk/payorder.html";
    } else {
      redirect_uri = "http://qrcode.cardinfolink.net/jssdk/payorder.html";
    }

    pay(orderData);


  };





  var pay=function(orderData){
     var now = new Date();
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
      "&chcd=WXP";
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
   
      var data ;
      if(orderData.attach!==undefined&&orderData.attach.length>0){
       data= "attach="+base64encode(utf16to8(orderData.attach))+
       "&backUrl="+orderData.backUrl;
      }else{
         data= "backUrl="+orderData.backUrl;
      }
      data=data+
      
      "&busicd=PURC"+
      "&chcd=WXP";
      if(orderData.goodsInfo!==undefined&&orderData.goodsInfo.length>0){
        data=data+
        "&goodsInfo="+base64encode(utf16to8(orderData.goodsInfo));
      }
      data=data+
      "&inscd="+orderData.inscd+
      "&mchntid="+orderData.merID+
      "&orderCurrency="+orderData.orderCurrency+
      "&orderNum="+orderData.orderNumber+
      "&txamt="+getTxamt(orderData.orderAmount)+
      "&txndir=Q"+
      "&sign="+signatureStr;
      data=data+"&frontUrl="+orderData.frontUrl;
      data=base64encode(utf16to8(data));
      redirect_uri=redirect_uri+"?data="+data;
    var url = "https://open.weixin.qq.com/connect/oauth2/authorize?appid=wx25ac886b6dac7dd2" + "&redirect_uri=" + encodeURI(redirect_uri) + "&response_type=code&scope=snsapi_base&state=123#wechat_redirect";
    window.open(url);
  }

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
    startOneKeyPay:startOneKeyPay
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
