var Trade = function() {

	var init = function() {

		var merchantCode = Util._getUrlParam("merchantCode");
		if (merchantCode === null) {
			return;
		}

		var data = "merchantCode=" + merchantCode;

		var XMLHttpReq, resultURL;; //
		var url = "http://test.quick.ipay.so/scanpay/fixed/orderInfo";

		$.ajax({
			type: 'POST',
			url: url,
			async: true,
			data: data,
			dataType: 'json',
			success: function(data) {

				if (data.response == "00") {
					var jsonobj = data;
					var txn = jsonobj.data;
					var title_one = jsonobj.title_one;
					var title_two = jsonobj.title_two;
					document.title = title_one + "-" + title_two;

					//alert(now.Format("YYYY")+'年'+now.Format("MM")+'月 交易笔数'+count+' 交易总金额：'+total+'元 退款'+refdcount+'笔（'+refdtotal+'元)')


					var orderList = [];
					for (var i = 0, l = txn.length; i < l; i++) {
						var order = new Object();
						var json = txn[i];
						order.headimgurl = json.headimgurl;
						order.nickname = json.nickname;
						order.code = json.veriCode;
						order.transtime = json.transtime;
						order.status = "交易成功";									
						order.amount = json.amount;
						order.orderNum = json.orderNum;
						orderList[i] = order;
					}


					var htmls = [];
					for (var j = 0, l = orderList.length; j < l; j++) {
						htmls[j] = _fillOrderToTemplate(orderList[j]);
					}


					$('#thelist').html(htmls.join(''));

				} else {
					alert(data.errorDetail);
					back();
				}
			},

			error: function(message) {
				alert(JSON.stringify(message));
				back();
			}
		});

	}

var _fillOrderToTemplate = function(order) {
	var source = $("#orderShortInfoTemplate").html();
	var template = Handlebars.compile(source);
	return template(order);
};

return {
	init: init
}
}();