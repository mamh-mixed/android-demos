'use strict';
require('../css/index.css');
require('../css/success.css');
require('./vendor/util');
var $ = require('webpack-zepto');


(function($) {
	$(function() {
		var state = Util.getUrlParam('state');
		switch (state) {
			case '0': // 成功
        $('#abc').text('支付成功');
				break;
			case '1': // 支付失败
        $('#abc').text(decodeURI(Util.getUrlParam('errorDetail')));
				break;
			case '-1': // 用户取消：未支付
        $('#abc').text('用户取消');
				break;
			default:
		}

		$('#returnBtn').on('click', () => {
			WeixinJSBridge.call('closeWindow');
		});
	});
})($);
