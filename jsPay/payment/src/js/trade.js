'use strict';

require('../css/trade.css');
require('./vendor/util');
var orderTemplate = require('../template/order.handlebars');

var $ = require('webpack-zepto');

let init = function() {
	let merchantCode = Util.getUrlParam('merchantCode');
	if (merchantCode === null) {
		return;
	}
	let data = 'merchantCode=' + merchantCode,
		url = Util.getServer() + '/scanpay/fixed/orderInfo';

	$.ajax({
		type: 'POST',
		url: url,
		data: data,
		dataType: 'json',
		success: (data) => {
			if (data.response !== '00') {
				window.alert(data.errorDetail);
				return;
			}

			document.title = data['title_one'] + '-' + data['title_two'];

			let orderList = [],
				htmls = [],
				txn = data.data || [];

			for (var i = 0, l = txn.length; i < l; i++) {
				var tx = txn[i],
					order = {};
				order.headimgurl = tx.headimgurl;
				order.nickname = tx.nickname;
				order.code = tx.veriCode;
				order.transtime = tx.transtime;
				order.status = '交易成功';
				order.amount = tx.amount;
				order.orderNum = tx.orderNum;
				orderList.push(order);
				htmls.push(orderTemplate(order));
			}

			if (txn.length === 0) {
				$('#thelist').html('<div style="text-align:center;padding-top:15px;">没有账单信息</div>');
				return;
			}
			$('#thelist').html(htmls.join(''));
		},
		error: (message) => {
			window.alert(JSON.stringify(message));
			WeixinJSBridge.call('closeWindow');
		}
	});
};

$(function() {
  init();
});
