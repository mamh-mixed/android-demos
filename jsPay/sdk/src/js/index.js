'use strict';

require('../css/index.css');
require('./vendor/base64');
// var Crypto = require('cryptojs').Crypto;
var sha = require('sha1');
var logImg = document.querySelector('#logo');
logImg.src = require('../img/logo.png');

const base64EncodeChars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/';
const base64DecodeChars = new Array(　　-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 　　-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 　　-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 62, -1, -1, -1, 63, 　　52, 53, 54, 55, 56, 57, 58, 59, 60, 61, -1, -1, -1, -1, -1, -1, 　　-1, 　0, 　1, 　2, 　3, 4, 　5, 　6, 　7, 　8, 　9, 10, 11, 12, 13, 14, 　　15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, -1, -1, -1, -1, -1, 　　-1, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 　　41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, -1, -1, -1, -1, -1);

/////////////////////////////////
// 讯联测试商户（受理商模式）
const signKey = 'zsdfyreuoyamdphhaweyrjbvzkgfdycs';
const appId = 'wx8854422b20240ed2';
const merId = '100000000010001';
const insCd = '1284577401';
/////////////////////////////////
function base64encode(str) {　　
	var out, i, len;　　
	var c1, c2, c3;　　
	len = str.length;　　
	i = 0;　　
	out = '';　　
	while (i < len) {
		c1 = str.charCodeAt(i++) & 0xff;
		if (i == len) {　　
			out += base64EncodeChars.charAt(c1 >> 2);　　
			out += base64EncodeChars.charAt((c1 & 0x3) << 4);　　
			out += '==';　　
			break;
		}
		c2 = str.charCodeAt(i++);
		if (i == len) {　　
			out += base64EncodeChars.charAt(c1 >> 2);　　
			out += base64EncodeChars.charAt(((c1 & 0x3) << 4) | ((c2 & 0xF0) >> 4));　　
			out += base64EncodeChars.charAt((c2 & 0xF) << 2);　　
			out += '=';　　
			break;
		}
		c3 = str.charCodeAt(i++);
		out += base64EncodeChars.charAt(c1 >> 2);
		out += base64EncodeChars.charAt(((c1 & 0x3) << 4) | ((c2 & 0xF0) >> 4));
		out += base64EncodeChars.charAt(((c2 & 0xF) << 2) | ((c3 & 0xC0) >> 6));
		out += base64EncodeChars.charAt(c3 & 0x3F);　　
	}　　
	return out;
}

function utf16to8(str) {　　
	var out, i, len, c;　　
	out = '';　　
	len = str.length;　　
	for (i = 0; i < len; i++) {
		c = str.charCodeAt(i);
		if ((c >= 0x0001) && (c <= 0x007F)) {　　
			out += str.charAt(i);
		} else if (c > 0x07FF) {　　
			out += String.fromCharCode(0xE0 | ((c >> 12) & 0x0F));　　
			out += String.fromCharCode(0x80 | ((c >> 　6) & 0x3F));　　
			out += String.fromCharCode(0x80 | ((c >> 　0) & 0x3F));
		} else {　　
			out += String.fromCharCode(0xC0 | ((c >> 　6) & 0x1F));　　
			out += String.fromCharCode(0x80 | ((c >> 　0) & 0x3F));
		}　　
	}　　
	return out;
}

let order = {
	'orderNum': (new Date()).getTime().toString(),
	'txamt': '000000000001',
	'orderCurrency': 'CNY',
	'backUrl': 'status.html',
	'frontUrl': 'status.html',
	'mchntid': merId,
	'inscd': insCd,
	'goodsInfo': '讯联测试,0.01,1',
	'attach': '讯联数据'
};

document.getElementById('jsonContainer').innerHTML = '<pre>' + JSON.stringify(order, null, 4) + '</pre>';

let signObj = {
	orderNum: order.orderNum,
	txamt: order.txamt,
	backUrl: order.backUrl,
	mchntid: order.mchntid
};

let plainTxts = [];
Object.keys(signObj).sort().forEach((key) => {
	plainTxts.push(key + '=' + signObj[key] || '');
});

let plainTxt = plainTxts.join('&') + signKey;
let sign = sha(plainTxt);
// order.sign = sign;
document.getElementById('sign').value = sign;


document.getElementById('wxpay').onclick = () => {
	order.sign = document.getElementById('sign').value;
	let base64Txt = base64encode(JSON.stringify(order))
	let redirectUrl = 'http://qrcode.cardinfolink.net/sdk/wxpayment.html?data=' + base64Txt;
	let postUrl = 'https://open.weixin.qq.com/connect/oauth2/authorize?appid=' + appId + '&redirect_uri=' + encodeURIComponent(redirectUrl) + '&response_type=code&scope=snsapi_base&state=123#wechat_redirect';
	window.location.replace(postUrl);
};
