var Trade = function() {

	var init = function() {

		var merchantCode = Util._getUrlParam("merchantCode");
		if (merchantCode === null) {
			return;
		}
		now = new Date(),
			transtime = now.Format("YYYYMMddHHmmss");

		var signatureStr = "merchantCode=" + merchantCode +
			"&transtime=" + transtime + "eu1dr0c8znpa43blzy1wirzmk8jqdaon";
		var sign = hex_sha1(signatureStr);

		var data = "merchantCode=" + merchantCode +
			"&transtime=" + transtime +
			"&sign=" + sign;

		var XMLHttpReq, resultURL;;//
		var url = "http://211.147.72.70:10003/scanMerchantCodeBill";
        //var url="http://192.168.199.174:8081/scanMerchantCodeBill"
      //  var url = "http://211.144.213.118/scanMerchantCodeBill";

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
					var jsonobj = eval('(' + text + ')');
					var txn = jsonobj.txn;
					var title_one=jsonobj.title_one;
					var title_two=jsonobj.title_two;
					document.title=title_one+"-"+title_two;
                
					//alert(now.Format("YYYY")+'年'+now.Format("MM")+'月 交易笔数'+count+' 交易总金额：'+total+'元 退款'+refdcount+'笔（'+refdtotal+'元)')


					var orderList = [];
					for (var i = 0, l = txn.length; i < l; i++) {
						var order = new Object();
						var json = txn[i];
						var m_request = json.m_request;        
						var userinfo = json.payjson.userinfo;

                      
						if (userinfo) {
							order.headimgurl = userinfo.headimgurl;
							order.nickname = userinfo.nickname;
							order.code = userinfo.code;

						}


						var transtime = json.system_date;
						order.transtime = transtime.substring(0, 4) + '-' + transtime.substr(4, 2) + '-' + transtime.substr(6, 2)+ ' ' + transtime.substr(8, 2)	+ ':' + transtime.substr(10, 2)	+ ':' + transtime.substr(12, 2);			
						var response = json.response;
						if (response === "00") {
							order.status = "交易成功";
							order.success = true;

						} else if (response === "09") {
							order.status = "未支付";
							order.fail = true;
						} else {
							order.status = "交易失败";
							order.fail = true;
						}


						var amount = m_request.txamt;
						if (amount != null && amount.length > 0) {
							order.amount = Util.getNormalTxamt(amount);
						}
						order.orderNum = m_request.orderNum;
						orderList[i] = order;
					}
                     
   						
   						var htmls = [];
					for (var j = 0, l = orderList.length; j < l; j++) {
						htmls[j] = _fillOrderToTemplate(orderList[j]);					
					}
                 

					$('#thelist').html(htmls.join(''));

				}
			}
		};

		XMLHttpReq.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
		XMLHttpReq.send(data);
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