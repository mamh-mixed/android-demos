'use strict';

require('./vendor/util');
var $ = require('webpack-zepto');

const base64EncodeChars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/';
const base64DecodeChars = new Array(　　-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 　　-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 　　-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 62, -1, -1, -1, 63, 　　52, 53, 54, 55, 56, 57, 58, 59, 60, 61, -1, -1, -1, -1, -1, -1, 　　-1, 　0, 　1, 　2, 　3, 4, 　5, 　6, 　7, 　8, 　9, 10, 11, 12, 13, 14, 　　15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, -1, -1, -1, -1, -1, 　　-1, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 　　41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, -1, -1, -1, -1, -1);

function base64decode(str) {　　
	var c1, c2, c3, c4;　　
	var i, len, out;　　
	len = str.length;　　
	i = 0;　　
	out = '';　　
	while (i < len) {
		/* c1 */
		do {　　
			c1 = base64DecodeChars[str.charCodeAt(i++) & 0xff];
		} while (i < len && c1 == -1);
		if (c1 == -1)　　 break;
		/* c2 */
		do {　　
			c2 = base64DecodeChars[str.charCodeAt(i++) & 0xff];
		} while (i < len && c2 == -1);
		if (c2 == -1)　　 break;
		out += String.fromCharCode((c1 << 2) | ((c2 & 0x30) >> 4));
		/* c3 */
		do {　　
			c3 = str.charCodeAt(i++) & 0xff;　　
			if (c3 == 61)　 return out;　　
			c3 = base64DecodeChars[c3];
		} while (i < len && c3 == -1);
		if (c3 == -1)　　 break;
		out += String.fromCharCode(((c2 & 0XF) << 4) | ((c3 & 0x3C) >> 2));
		/* c4 */
		do {　　
			c4 = str.charCodeAt(i++) & 0xff;　　
			if (c4 == 61)　 return out;　　
			c4 = base64DecodeChars[c4];
		} while (i < len && c4 == -1);
		if (c4 == -1)　　 break;
		out += String.fromCharCode(((c3 & 0x03) << 6) | c4);　　
	}　　
	return out;
}

function utf8to16(str) {　　
	var out, i, len, c;　　
	var char2, char3;　　
	out = '';　　
	len = str.length;　　
	i = 0;　　
	while (i < len) {
		c = str.charCodeAt(i++);
		switch (c >> 4) {　
			case 0:
			case 1:
			case 2:
			case 3:
			case 4:
			case 5:
			case 6:
			case 7:
				　　 // 0xxxxxxx
				　　out += str.charAt(i - 1);　　
				break;　
			case 12:
			case 13:
				　　 // 110x xxxx　 10xx xxxx
				　　char2 = str.charCodeAt(i++);　　
				out += String.fromCharCode(((c & 0x1F) << 6) | (char2 & 0x3F));　　
				break;　
			case 14:
				　　 // 1110 xxxx　10xx xxxx　10xx xxxx
				　　char2 = str.charCodeAt(i++);　　
				char3 = str.charCodeAt(i++);　　
				out += String.fromCharCode(((c & 0x0F) << 12) | 　　　　((char2 & 0x3F) << 6) | 　　　　((char3 & 0x3F) << 0));　　
				break;
		}　　
	}　　
	return out;
}

(function($) {
	$(function() {

		let [code, orderData, url] = [Util.getUrlParam('code'), Util.getUrlParam('data'), Util.getServer() + '/scanpay/unified'];

		if (orderData === null) {
			window.alert('没有订单信息');
			return;
		}

		orderData = utf8to16(base64decode(orderData));
		var order = null;

		try {
			order = JSON.parse(orderData);
		} catch (e) {
			window.alert('不合法的JSON字符串');
			return;
		}

		let [frontUrl, errorData] = [
			order.frontUrl, [
				'attach=' + order.attach,
				'txamt=' + order.txamt,
				'goodsInfo=' + order.goodsInfo,
				'orderCurrency=' + order.orderCurrency,
				'orderNum=' + order.orderNum,
				'state=1'
			].join('&')
		];


		if (!code) {
			errorData = errorData + '&errorDetail=URL拼接错误';
			window.location.replace(frontUrl + '?' + errorData);
			return;
		}

		let postData = {
			orderNum: order.orderNum,
			sign: order.sign,
			txamt: order.txamt,
			backUrl: order.backUrl,
			mchntid: order.mchntid,
			inscd: order.inscd,
			txndir: 'Q',
			goodsInfo: order.goodsInfo,
			chcd: 'WXP',
			busicd: 'JSZF',
			needUserInfo: 'NO',
			code: code,
			attach: order.attach,
			currency: order.orderCurrency
		};

		$.ajax({
			type: 'POST',
			url: url,
			async: false,
			data: JSON.stringify(postData),
			dataType: 'json',
			success: (data) => {
				let payJson = data.payjson;

				if (!payJson) {
					errorData += '&errorDetail=' + encodeURI(encodeURI(data.errorDetail));
					window.location.replace(frontUrl + '?' + errorData);
					return;
				}

				let {
					config: aaa,
					chooseWXPay: bbb
				} = payJson;

				let {
					appId: configAppId,
					timestamp: configTimestamp,
					nonceStr: configNonceStr,
					signature: configSignature
				} = aaa;

				let {
					timeStamp: chooseWxPayTimestamp,
					nonceStr: chooseWxPayNonceStr,
					package: chooseWxPayPackage,
					signType: chooseWxPaySignType,
					paySign: chooseWxPayPaySign
				} = bbb;

				let post = [
					'attach=' + data.attach,
					'txamt=' + data.txamt,
					'goodsInfo=' + orderData.goodsInfo,
					'orderCurrency=' + orderData.orderCurrency,
					'orderNum=' + data.orderNum
				].join('&');

				wx.config({
					debug: false, // 开启调试模式,调用的所有api的返回值会在客户端alert出来，若要查看传入的参数，可以在pc端打开，参数信息会通过log打出，仅在pc端时才会打印。
					appId: configAppId, // 必填，公众号的唯一标识
					timestamp: configTimestamp, // 必填，生成签名的时间戳
					nonceStr: configNonceStr, // 必填，生成签名的随机串
					signature: configSignature, // 必填，签名，见附录1
					jsApiList: [
							'checkJsApi',
							'chooseWXPay'
						] // 必填，需要使用的JS接口列表，所有JS接口列表见附录2
				});

				wx.checkJsApi({
					jsApiList: ['chooseWXPay'], // 需要检测的JS接口列表，所有JS接口列表见附录2,
					success: function(res) {
						alert(JSON.stringify(res));
						// 以键值对的形式返回，可用的api值true，不可用为false
						// 如：{"checkResult":{"chooseImage":true},"errMsg":"checkJsApi:ok"}
					}
				});

				wx.ready(() => {
					wx.chooseWXPay({
						timestamp: chooseWxPayTimestamp, // 支付签名时间戳，注意微信jssdk中的所有使用timestamp字段均为小写。但最新版的支付后台生成签名使用的timeStamp字段名需大写其中的S字符
						nonceStr: chooseWxPayNonceStr, // 支付签名随机串，不长于 32 位
						package: chooseWxPayPackage, // 统一支付接口返回的prepay_id参数值，提交格式如：prepay_id=***）
						signType: chooseWxPaySignType, // 签名方式，默认为'SHA1'，使用新版支付需传入'MD5'
						paySign: chooseWxPayPaySign, // 支付签名
						success: (res) => {
							console.log('success');
							post = post + '&state=0';
							window.location.replace(frontUrl + '?' + post);
						},
						fail: () => {
							console.log('fail');
							post = post + '&state=1';
							window.location.replace(frontUrl + '?' + post);
						},
						cancel: () => {
							console.log('cancel');
							post = post + '&state=-1';
							window.location.replace(frontUrl + '?' + post);
						}
					});
				});
			},
			error: (message) => {
				console.log('error', JSON.stringify(message));
				window.alert(JSON.stringify(message));
			}
		});
	});
})($);
