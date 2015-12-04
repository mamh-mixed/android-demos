'use strict';

require('../css/index.css');
require('../css/success.css');
require('./vendor/util');
var $ = require('webpack-zepto');

var isButtonCumstome = false;
var code = Util.getUrlParam('code');
var orderAmount = Util.getUrlParam('orderAmount');
document.getElementById('code').innerHTML = code;
document.getElementById('amount').innerHTML = '付款金额:¥' + orderAmount;
document.title = window.localStorage.getItem('title_one');

(function() {
	$(document).ready(() => {
		let btnText = window.sessionStorage.getItem('successButtonText');
		if (btnText && btnText !== '') {
			isButtonCumstome = true;
		}
		let btnLink = window.sessionStorage.getItem('successButtonLink');
		let isPostAmount = window.sessionStorage.getItem('isPostAmount');

		$('#customeButton').text(btnText);
		if (isPostAmount && isPostAmount === 'true') {
			if (btnLink.indexOf('?') > 0) {
				btnLink += '&amount=' + orderAmount;
			} else {
				btnLink += '?amount=' + orderAmount;
			}
		}
		$('#customeButton').attr('href', btnLink);

		if (isButtonCumstome) {
			$('#returnBtn').addClass('hidden');
			$('#customeButton').removeClass('hidden');
		} else {
			$('#returnBtn').removeClass('hidden');
			$('#customeButton').addClass('hidden');
		}

    $('#returnBtn').on('click', () => {
			WeixinJSBridge.call('closeWindow');
		});
	});
})($);
