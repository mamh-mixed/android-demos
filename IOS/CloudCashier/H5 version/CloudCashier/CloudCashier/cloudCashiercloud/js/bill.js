var Bill = function() {

    var initBill = function(callback) {
        var password = Util.getUrlParam("password"),
            username = Util.getUrlParam("username"),
            key = Util.getUrlParam("key"),
            clientid = Util.getUrlParam("clientid"),
            now = new Date(),
            transtime = now.Format("YYYYMMddHHmmss"),
            month = now.Format("YYYYMM"),
            index = 0,
            password = hex_md5(password);

        window.localStorage.setItem("bill_index", index);


        var status = "all";
        if (document.getElementById('bill_all').className === 'active') {
            status = "all";
        } else if (document.getElementById('bill_success').className === 'active') {
            status = "success";
        } else {
            status = "fail";
        }


        var signatureStr = "clientid=" + clientid +
            "&index=" + index +
            "&month=" + month +
            "&password=" + password +
            "&status=" + status +
            "&transtime=" + transtime +
            "&username=" + username + key;


        var sign = hex_sha1(signatureStr);

        var data = "clientid=" + clientid +
            "&index=" + index +
            "&month=" + month +
            "&password=" + password +
            "&status=" + status +
            "&transtime=" + transtime +
            "&username=" + username +
            "&sign=" + sign;

        send(data, callback);


    };


    var getUpBill = function(callback) {
        var password = Util.getUrlParam("password"),
            username = Util.getUrlParam("username"),
            key = Util.getUrlParam("key"),
            clientid = Util.getUrlParam("clientid"),
            now = new Date(),
            transtime = now.Format("YYYYMMddHHmmss"),
            month = now.Format("YYYYMM"),
            index_str = window.localStorage.getItem("bill_index"),
            count_str = window.localStorage.getItem("bill_count");
        var _count = Number(count_str);
        var index = 0;
        if (index_str !== null) {
            index = Number(index_str) + 15;
            if (index > _count) {
                var orderList = [];
                callback(orderList);
                return;
            }
        }
        window.localStorage.setItem("bill_index", index);
        password = hex_md5(password);
        var status = "all";
        if (document.getElementById('bill_all').className === 'active') {
            status = "all";
        } else if (document.getElementById('bill_success').className === 'active') {
            status = "success";
        } else {
            status = "fail";
        }


        var signatureStr = "clientid=" + clientid +
            "&index=" + index +
            "&month=" + month +
            "&password=" + password +
            "&status=" + status +
            "&transtime=" + transtime +
            "&username=" + username + key;


        var sign = hex_sha1(signatureStr);

        var data = "clientid=" + clientid +
            "&index=" + index +
            "&month=" + month +
            "&password=" + password +
            "&status=" + status +
            "&transtime=" + transtime +
            "&username=" + username +
            "&sign=" + sign;


        send(data, callback);
    };

    var getDownBill = function(callback) {
        var password = Util.getUrlParam("password"),
            username = Util.getUrlParam("username"),
            key = Util.getUrlParam("key"),
            clientid = Util.getUrlParam("clientid"),
            now = new Date(),
            transtime = now.Format("YYYYMMddHHmmss"),
            month = now.Format("YYYYMM"),
            index_str = window.localStorage.getItem("bill_index"),
            count_str = window.localStorage.getItem("bill_count");
        var index = 0;

        if (index_str !== null) {
            index = Number(index_str) - 15;
            if (-15 < index < 0) {
                index = 0;
            } else {
                var orderList = [];
                callback(orderList);
                return;
            }
        }

        password = hex_md5(password);
        window.localStorage.setItem("bill_index", index);

        var status = "all";
        if (document.getElementById('bill_all').className === 'active') {
            status = "all";
        } else if (document.getElementById('bill_success').className === 'active') {
            status = "success";
        } else {
            status = "fail";
        }


        var signatureStr = "clientid=" + clientid +
            "&index=" + index +
            "&month=" + month +
            "&password=" + password +
            "&status=" + status +
            "&transtime=" + transtime +
            "&username=" + username + key;


        var sign = hex_sha1(signatureStr);

        var data = "clientid=" + clientid +
            "&index=" + index +
            "&month=" + month +
            "&password=" + password +
            "&status=" + status +
            "&transtime=" + transtime +
            "&username=" + username +
            "&sign=" + sign;


        send(data, callback);
    };

    var send = function(data, callback) {


        var XMLHttpReq, resultURL, url = "http://211.144.213.118/bill";
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
                    var count = jsonobj.count;
                    var total = jsonobj.total;
                    var refdcount = jsonobj.refdcount;
                    var refdtotal = jsonobj.refdtotal;
                    var txn = jsonobj.txn;
                    //alert(now.Format("YYYY")+'年'+now.Format("MM")+'月 交易笔数'+count+' 交易总金额：'+total+'元 退款'+refdcount+'笔（'+refdtotal+'元)')


                    var orderList = [];
                    for (var i = 0, l = txn.length; i < l; i++) {
                        var order = new Object();
                        var json = txn[i];
                        var m_request = json.m_request;
                        var busicd = m_request.busicd;

                        if (busicd === "PURC") {
                            order.type = "下单支付";
                        } else if (busicd === "PAUT") {
                            order.type = "预下单支付";
                        } else if (busicd === "REFD") {
                            order.type = "退款";
                        } else if (busicd === "VOID") {
                            order.type = "撤销";
                        } else if (busicd === "CANC") {
                            order.type = "取消";
                        }

                        var chcd = m_request.chcd;
                        if (chcd === "WXP") {
                            order.img = "img/wpay.png";
                        } else {
                            order.img = "img/apay.png";
                        }



                        var date = json.system_date;
                        order.date = date;
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


                        var money = m_request.txamt;
                        if (money != null && money.length > 0) {
                            order.money = Util.getNormalTxamt(money);
                        }
                        order.orderNum = m_request.orderNum;
                        orderList[i] = order;


                    }

                    window.localStorage.setItem("bill_count", count);
                    now = new Date();
                    document.getElementById('bill_title').innerHTML = now.Format("YYYY") + '年' + now.Format("MM") + '月 交易笔数' + count + ' 交易总金额：' + total + '元 退款' + refdcount + '笔（' + refdtotal + '元)';
                    callback(orderList);

                }
            }
        };
        XMLHttpReq.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
        XMLHttpReq.send(data);
    };
    var getOrderInfo = function(orderNum, callback) {

        var password = Util.getUrlParam("password"),
            username = Util.getUrlParam("username"),
            key = Util.getUrlParam("key"),
            clientid = Util.getUrlParam("clientid"),
            now = new Date(),
            transtime = now.Format("YYYYMMddHHmmss");

        password = hex_md5(password);

        var signatureStr = "clientid=" + clientid +
            "&orderNum=" + orderNum +
            "&password=" + password +
            "&transtime=" + transtime +
            "&username=" + username + key;


        var sign = hex_sha1(signatureStr);

        var data = "clientid=" + clientid +
            "&orderNum=" + orderNum +
            "&password=" + password +
            "&transtime=" + transtime +
            "&username=" + username +
            "&sign=" + sign;

        var XMLHttpReq, resultURL, url = @"http://211.144.213.118/getOrder";//url = "http://211.147.72.70:10003/getOrder",
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
                    var json = jsonobj.txn;
                    var order = new Object();
                    var m_request = json.m_request;
                    var busicd = m_request.busicd;
                    var response = json.response;
                    if (busicd === "PURC" || busicd === "PAUT") {
                        if (response === "00") {
                            order.refd = true;
                        }

                    }

                    var chcd = m_request.chcd;
                    if (chcd === "WXP") {
                        order.img = "img/wpay.png";
                        order.chcd = "微信扫码支付";
                    } else {
                        order.img = "img/apay.png";
                        order.chcd = "支付宝扫码支付";
                    }



                    var date = json.system_date;
                    order.transtime = date;

                    if (response === "00") {
                        order.status = "交易成功";
                        order.success = true;

                    } else if (response === "09") {
                        order.status = "未支付";
                    } else {
                        order.status = "交易失败";
                    }


                    var money = m_request.txamt;
                    if (money != null && money.length > 0) {
                        order.money = Util.getNormalTxamt(money);
                    }
                    order.orderNum = m_request.orderNum;
                    order.username = Util.getUrlParam("username");
                    order.goodsInfo = m_request.goodsInfo;


                    callback(order);

                }
            }
        };
        XMLHttpReq.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
        XMLHttpReq.send(data);

    }

    var refd = function(orderNum, total) {
        var device = Util.getUrlParam("device");
        var json = {
            "orderNum": orderNum,
            "total": total
        };
        var jsonStr = JSON.stringify(json);
        if (device == "android") {
            window.couldCashier.OnDo("refd", jsonStr);

        } else if (device == "ios") {
            var result = "couldcashier://refd:/" + jsonStr;
            document.location = result;
        }

    }
    return {
        initBill: initBill,
        getUpBill: getUpBill,
        getDownBill: getDownBill,
        getOrderInfo: getOrderInfo,
        refd: refd
    }
}();
// 交易管理
var TransManage = function() {
    var myScroll,
        pullDownEl, pullDownOffset,
        pullUpEl, pullUpOffset;
    var init = function() {
        _events();
    };
    var _events = function() {
        // 全部交易，交易成功，交易失败切换事件
        $('#transTypeBtnGroup').find('button').on(CLICK, function() {
            $(this).addClass('active')
                .siblings().removeClass('active');
            initScroll();
        });
        document.addEventListener('touchmove', function(e) {
            e.preventDefault();
        }, false);

        document.addEventListener('DOMContentLoaded', function() {
            setTimeout(loaded, 200);
        }, false);

    };
    // 使用handlebar渲染订单信息
    var _fillOrderToTemplate = function(order) {
        var source = $("#orderShortInfoTemplate").html();
        var template = Handlebars.compile(source);
        return template(order);
    };
    var pullDownAction = function() {
        Bill.getDownBill(function(orderList) {
            var htmls = [];
            for (var i = 0, l = orderList.length; i < l; i++) {
                htmls[i] = _fillOrderToTemplate(orderList[i]);
            }
            $('#thelist').prepend(htmls.join(''));
            myScroll.refresh();
        });
    };
    var pullUpAction = function() {
        Bill.getUpBill(function(orderList) {
            var htmls = [];
            for (var i = 0, l = orderList.length; i < l; i++) {
                htmls[i] = _fillOrderToTemplate(orderList[i]);
            }
            $('#thelist').prepend(htmls.join(''));
            myScroll.refresh();
        });
    };

    var initScroll = function() {
        $('#thelist').text("");
        window.setTimeout(function() {
            Bill.initBill(function(orderList) {
                var htmls = [];
                for (var i = 0, l = orderList.length; i < l; i++) {
                    htmls[i] = _fillOrderToTemplate(orderList[i]);
                }
                $('#thelist').prepend(htmls.join(''));
                myScroll.refresh();
            });
        }, 300);
    };
    var loaded = function() {
        if (myScroll) {
            myScroll.destroy();
        }

        pullDownEl = document.getElementById('pullDown');
        pullDownOffset = pullDownEl.offsetHeight;
        pullUpEl = document.getElementById('pullUp');
        pullUpOffset = pullUpEl.offsetHeight;

        // myScroll = new iScroll('wrapper', {
        //     useTransition: true,
        //     topOffset: pullDownOffset,
        //     onRefresh: function() {
        //         if (pullDownEl.className.match('loading')) {
        //             pullDownEl.className = '';
        //             pullDownEl.querySelector('.pullDownLabel').innerHTML = '下拉刷新...';
        //         } else if (pullUpEl.className.match('loading')) {
        //             pullUpEl.className = '';
        //             pullUpEl.querySelector('.pullUpLabel').innerHTML = '上拉加载更多...';
        //         }
        //     },
        //     onScrollMove: function() {
        //         if (this.y > 5 && !pullDownEl.className.match('flip')) {
        //             pullDownEl.className = 'flip';
        //             pullDownEl.querySelector('.pullDownLabel').innerHTML = '松开刷新...';
        //             this.minScrollY = 0;
        //         } else if (this.y < 5 && pullDownEl.className.match('flip')) {
        //             pullDownEl.className = '';
        //             pullDownEl.querySelector('.pullDownLabel').innerHTML = '下拉刷新...';
        //             this.minScrollY = -pullDownOffset;
        //         } else if (this.y < (this.maxScrollY - 5) && !pullUpEl.className.match('flip')) {
        //             pullUpEl.className = 'flip';
        //             pullUpEl.querySelector('.pullUpLabel').innerHTML = '松开刷新...';
        //             this.maxScrollY = this.maxScrollY;
        //         } else if (this.y > (this.maxScrollY + 5) && pullUpEl.className.match('flip')) {
        //             pullUpEl.className = '';
        //             pullUpEl.querySelector('.pullUpLabel').innerHTML = '上拉加载更多...';
        //             this.maxScrollY = pullUpOffset;
        //         }
        //     },
        //     onScrollEnd: function() {
        //         if (pullDownEl.className.match('flip')) {
        //             pullDownEl.className = 'loading';
        //             pullDownEl.querySelector('.pullDownLabel').innerHTML = 'Loading...';
        //             pullDownAction(); // Execute custom function (ajax call?)
        //         } else if (pullUpEl.className.match('flip')) {
        //             pullUpEl.className = 'loading';
        //             pullUpEl.querySelector('.pullUpLabel').innerHTML = 'Loading...';
        //             pullUpAction(); // Execute custom function (ajax call?)
        //         }
        //     }
        // });
        setTimeout(function() {
            document.getElementById('wrapper').style.left = '0';
        }, 800);
    };
    return {
        init: init,
        loaded: loaded,
        initScroll: initScroll
    };
}(window);

// 账单详情
var TransDetail = function() {
    var init = function(orderId) {
        _events();
    };
    // 载入数据
    var loadedData = function(orderId) {
        // console.log('Tran order id is ', orderId);
        $('#orderDetailInfo').html("");
        setTimeout(function() {
            Bill.getOrderInfo(orderId, function(order) {

                $('#orderDetailInfo').html(_fillOrderDetailInfoToTemplate(order));
            });
        }, 0);
    };
    var _fillOrderDetailInfoToTemplate = function(order) {
        var source = $("#orderDetailInfoTemplate").html();
        var template = Handlebars.compile(source);
        return template(order);
    };
    var _events = function() {
        $('#returnBtn').click(function() {
            $('#transManage').show()
            history.go(-1);
        });
    };
    return {
        init: init,
        loadedData: loadedData
    };
}(window);
