'use strict';

require('../css/index.css');
require('../css/success.css');
require('./vendor/util');
var $ = require('webpack-zepto');


let code = Util.getUrlParam('code');
let orderAmount = Util.getUrlParam('orderAmount');
document.getElementById('code').innerHTML = code;
document.getElementById('amount').innerHTML = '付款金额:¥' + orderAmount;
document.title = window.localStorage.getItem('title_one');

(function() {
	$(document).ready(() => {
    $('#returnBtn').on('click', () => {
			WeixinJSBridge.call('closeWindow');
		});
	});
})($);
