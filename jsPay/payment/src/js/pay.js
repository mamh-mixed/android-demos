'use strict';

require('../css/index.css');
require('./vendor/util');
require('./vendor/yunshouyin-1.0')
var $ = require('webpack-zepto');
var logImg = document.querySelector('#logo');
logImg.src = require('../img/logo.png');

var merID, inscd;
var submit = false;

(function($) {
	$(function() {
		init();

		$('#wxpay').on('click', ()=>{
			wpay();
		});
	});
})($);

/**
 * 支付
 */
function wpay() {
	if (submit === true) {
		return;
	}

	let code = Util.getUrlParam('code'),
		money = $('#money').val();
	if (check(money)) {
		submit = true;
		let orderAmount = parseFloat(money);
		orderAmount = orderAmount.toFixed(2);
		var currency = 'CNY';
		var now = new Date();
		var orderNum = now.Format('YYMMddHHmmss');
		var veriCode = ''
		for (var i = 0; i < 5; i++) {
			orderNum = orderNum + Math.floor(Math.random() * 10);
		}

		for (var j = 0; j < 4; j++) {
			veriCode = veriCode + Math.floor(Math.random() * 10);
		}

		var orderData = {
			orderNum: orderNum,
			txamt: '' + Util.getTxamt(orderAmount),
			orderCurrency: currency,
			backUrl: '',
			frontUrl: 'payresult.html',
			mchntid: merID,
			inscd: inscd,
			goodsInfo: '云收银wap客户端',
			attach: '用户附加数据原样返回',
			veriCode: veriCode

		};

		yunshouyin.startWPayWithCode(orderData, code);
	}
}

/**
 * 检查金额是否OK
 */
function check(v) {
	if (v.length === 0) {
		window.alert('金额不能为空');
		return false;
	}
	var a = /^[0-9]*(\.[0-9]{1,2})?$/;
	if (!a.test(v)) {
		window.alert('金额不正确');
		return false;
	}
	return true;

}

/**
 * 页面初始化方法
 */
function init() {
	var merchantCode = Util.getUrlParam('merchantCode');
	if (merchantCode === null) {
		return;
	}
	var data = 'merchantCode=' + merchantCode
	var url = Util.getServer() + '/scanpay/fixed/merInfo';
	$.ajax({
		type: 'POST',
		url: url,
		async: true,
		data: data,
		dataType: 'json',
		success: (data) => {

			if (data.response === '00') {
				merID = data.merID;
				inscd = data.inscd;
				let titleOne = data['title_one'];
				let titleTwo = data['title_two'];
				$('#titleOne').text('欢迎光临' + titleOne);
				$('#titleTwo').text('本店名称-' + titleTwo);
				document.title = titleOne;
				window.localStorage.setItem('title_one', titleOne);
				if (data.succBtnTxt && data.succBtnTxt !== '') {
					window.sessionStorage.setItem('successButtonText', data.succBtnTxt.trim());
				}
				if (data.succBtnLink && data.succBtnLink !== '') {
					window.sessionStorage.setItem('successButtonLink', data.succBtnLink.trim());
				}
				if (data.isPostAmount) {
					window.sessionStorage.setItem('isPostAmount', data.isPostAmount);
				}
			} else {
				window.alert(data.errorDetail);
				WeixinJSBridge.call('closeWindow');
			}
		},

		error: (message) => {
			window.alert(JSON.stringify(message));
			WeixinJSBridge.call('closeWindow');
		}
	});
}
