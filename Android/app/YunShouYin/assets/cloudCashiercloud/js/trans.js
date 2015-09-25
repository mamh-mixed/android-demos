var TransManage = function() {
    var username, password, clientid, status, index = 0;
    var threshold = 40;
    var isLoading = false;
    var pageSize = 15;
    var MIN_HEIGHT = 32;
    var isAndroid = (/android/gi).test(navigator.appVersion);
    var init = function() {
        _events();
    };
    var _events = function() {
        // 全部交易，交易成功，交易失败切换事件
        $('#transTypeBtnGroup').find('button').on(CLICK, function() {
            $(this).addClass('active')
                .siblings().removeClass('active');
            status = $(this).data('status');
            index = 0;
            $('#wrapper').scrollTop(0);
            _freshData();
        });

        $('#thelist').on(CLICK, '.collection-item', function() {
            var orderId = $(this).data('oid');
            window.location.hash = '#/tranDetail/' + orderId;
        });

        // 账单列表触摸事件
        var theList = document.getElementById('thelist'),
            originX = 0,
            originY = 0,
            offsetX = 0,
            offsetY = 0,
            deltaX = 0,
            deltaY = 0,
            direction = '';
        var scrollTop = 0,
            wrapSpand = 0,
            listHeight = 0;
        theList.addEventListener('touchstart', function(e, a, b, c) {
            originY = e.touches[0].pageY;
            originX = e.touches[0].pageX;
            direction = '';
        });
        theList.addEventListener('touchmove', function(e, a, b, c) {
            if (isLoading) {
                return;
            }
            offsetY = e.touches[0].pageY;
            offsetX = e.touches[0].pageX;
            deltaX = offsetX - originX;
            deltaY = offsetY - originY;
            if (direction === '') {
                if (Math.abs(deltaX) > Math.abs(deltaY)) {
                    // 左右划，不处理
                    return;
                }
                if (deltaY > 0) {
                    direction = 'DOWN';
                } else {
                    direction = 'UP';
                }

            }

            scrollTop = $('#wrapper').scrollTop();
            // 滚动到最后顶部并且向下滑动
            if (direction === 'DOWN' && scrollTop <= 5) {
                if (isAndroid) {
                    e.preventDefault();
                }
                $('#topFresh').show().height(Math.abs(deltaY));
                if (Math.abs(deltaY) > MIN_HEIGHT) {
                    $('#freshContainer').show();
                } else {
                    $('#freshContainer').hide().find('i').removeClass('fa-spin');
                }
            }

            wrapSpand = $('#wrapper').height() + scrollTop;
            listHeight = $('#thelist').height();
            // 滚动到最后底部并且向上滑动
            if (direction === 'UP' && wrapSpand >= (listHeight - 10)) {
                if (isAndroid) {
                    e.preventDefault();
                }
                // 大于临界值
                if (Math.abs(deltaY) > threshold) {
                    $('#bottomFresh').show();
                } else {
                    $('#bottomFresh').hide();
                }
            }

        });
        theList.addEventListener('touchend', function() {
            if (direction === '' || isLoading) {
                return;
            }
            $('#topFresh').height(0).hide();
            // $('#topFresh').height(MIN_HEIGHT);
            $('#bottomFresh').height(MIN_HEIGHT);
            if (Math.abs(deltaY) < threshold) {
                return;
            }
            // 超过阈值了
            // 滚动到最后顶部并且向下滑动
            if (direction === 'DOWN' && scrollTop <= 5) {
                $('#freshContainer').find('i').addClass('fa-spin');
                $('#wrapper').scrollTop(5);
                _freshData();
                isLoading = true;
            }
            // 滚动到最后底部并且向上滑动
            if (direction === 'UP' && wrapSpand >= (listHeight - 10)) {
                // 加载更多数据
                _loadingMoreData();
                isLoading = true;
            }
        });
    };
    // 下拉刷新数据
    var _freshData = function() {
        index = 0;
        // 刷新数据
        $('#freshContainer').show().find('i').addClass('fa-spin');
        _requestForNewData(function(orders) {
            var htmls = [];
            for (var i = 0, l = orders.length; i < l; i++) {
                htmls[htmls.length] = _fillOrderToTemplate(orders[i]);
            }
            if (htmls.length > 0) {
                $('#thelist').html(htmls.join(''));
            } else {
                $('#thelist').html('<div style="text-align: center;margin: 100px;">没有数据</div>');
            }
            $('#freshContainer').hide().find('i').removeClass('fa-spin');
            isLoading = false;
        });
    };
    // 上拉加载更多数据
    var _loadingMoreData = function() {
        // 加载更多数据
        index += pageSize;
        _requestForNewData(function(orders) {
            var htmls = [];
            for (var i = 0, l = orders.length; i < l; i++) {
                htmls[htmls.length] = _fillOrderToTemplate(orders[i]);
            }
            $('#thelist').append(htmls.join(''));
            $('#bottomFresh').hide();
            isLoading = false;
        });
    };
    // ajax请求最新数据
    var _requestForNewData = function(callback) {
        var now = new Date();
        var dataObj = {
            clientid: Util.getClientId(),
            index: index || 0,
            month: now.Format('YYYYMM'),
            password: hex_md5(Util.getPassword()),
            status: status || 'all',
            transtime: now.Format('YYYYMMddHHmmss'),
            username: Util.getUsername()
        };
        var requestData = _generateRequestData(dataObj);
        $.ajax({
            type: 'POST',
            url: UrlMap.bill,
            data: requestData,
            dataType: 'json',
            success: function(data) {
                var now = new Date(),
                    count = data.count,
                    total = data.total,
                    refdcount = data.refdcount,
                    refdtotal = data.refdtotal,
                    orders = _markupData(data.txn);
               document.getElementById('bill_title').innerHTML = now.Format("YYYY") + '年' + now.Format("MM") + '月 交易笔数' + count + ' 交易总金额：' + total + '元 退款' + refdcount + '笔（' + refdtotal + '元)';
                callback(orders);
            }
        });

    };
    // 整理响应回来的订单列表
    var _markupData = function(data) {
        if (!data) {
            return [];
        }

        var busicdMap = {
            'PURC': '支付',
            'PAUT': '支付',
            'REFD': '退款',
            'VOID': '撤销',
            'CANC': '取消'
        };
        var tradeFromMap = {
            'app': 'APP',
            'wap': 'WAP'
        };
        var orders = [];
        for (var i = 0, l = data.length; i < l; i++) {
            var order = {};
            var rawOrder = data[i];
            var m_request = rawOrder.m_request;
            var busicd = m_request.busicd;
            var tradeFrom = m_request.tradeFrom;

            order.type = [tradeFromMap[tradeFrom] || 'PC', busicdMap[busicd]].join(' ');

            var chcd = m_request.chcd;
            if (chcd === 'WXP') {
                order.img = 'img/wechat@2x.png';
            } else {
                order.img = 'img/alipay@2x.png';
            }

            var transtime = rawOrder.system_date;
            order.date =  transtime.substring(0, 4) + '-' + transtime.substr(4, 2) + '-' + transtime.substr(6, 2)+ ' ' + transtime.substr(8, 2)	+ ':' + transtime.substr(10, 2)	+ ':' + transtime.substr(12, 2);
            var response = rawOrder.response;

            switch (response) {
                case '00':
                    order.status = '交易成功';
                    order.success = true;
                    break;
                case '09':
                    order.status = '未支付';
                    order.fail = true;
                    break;
                default:
                    order.status = '交易失败';
                    order.fail = true;
            }

            var money = m_request.txamt;
            if (money && money.length > 0) {
                order.money = Util.getNormalTxamt(money);
            }
            order.orderNum = m_request.orderNum;

            orders[orders.length] = order;
        }

        return orders;
    };
    // 对查询数据计算签名
    var _generateRequestData = function(obj) {
        var temp = [];
        Object.keys(obj).sort().forEach(function(v, i) {
            temp.push(v + '=' + obj[v]);
        });

        var key = Util.getKey();
        var plainText = temp.join('&') + key;

        var sign = hex_sha1(plainText);

        return temp.join('&') + '&sign=' + sign;
    };
    // 使用handlebar渲染订单信息
    var _fillOrderToTemplate = function(order) {
        var source = $("#orderShortInfoTemplate").html();
        var template = Handlebars.compile(source);
        return template(order);
    };
    return {
        init: init,
        freshData: _freshData
    };
}(window);

var TransDetail = function() {
    var username, password, key, clientid;
    var orderNum = '';
    var refundMoney = 0;
    var init = function() {
        initData();
        _events();
    };
    var _events = function() {
        $('#returnBtn')[CLICK](function() {
            window.location.hash = '#/transManage';
        });
        $('#orderDetailInfo').on(CLICK, '#refundBtn', function() {
            refd();
        });
    };
    // 退款
    var refd = function() {
        var json = {
            "orderNum": orderNum,
            "total": refundMoney
        };
        var jsonStr = JSON.stringify(json);
        Util.sendRequestToDevice('refd', jsonStr);
    };
    var initData = function() {
        password = Util.getPassword();
        username = Util.getUsername();
        key = Util.getKey();
        clientid = Util.getClientId();
    };
    // 载入数据
    var loadedData = function(orderId) {
        orderNum = orderId;
        initData();
        $('#orderDetailInfo').html('');
        var url = UrlMap.getOrder,
            now = new Date(),
            transtime = now.Format('YYYYMMddHHmmss'),
            dataArray = [];

        password = hex_md5(password);

        dataArray.push('clientid=' + clientid);
        dataArray.push('orderNum=' + orderId);
        dataArray.push('password=' + password);
        dataArray.push('transtime=' + transtime);
        dataArray.push('username=' + username);

        var signatureStr = dataArray.join('&') + key,
            sign = hex_sha1(signatureStr);

        dataArray.push('sign=' + sign);
        var data = dataArray.join('&');

        $.ajax({
            type: 'POST',
            url: url,
            data: data,
            dataType: 'json',
            success: function(data) {
                var json = data.txn;
                if (!json) {
                    window.alert('获取订单详情失败');
                    return;
                }
                var order = _makeupOrder(json);
                refundMoney = order.money;
                $('#orderDetailInfo').html(_fillOrderDetailInfoToTemplate(order));
            }
        });
    };
    // 整理订单数据
    var _makeupOrder = function(json) {
        var order = {},
            m_request = json.m_request,
            busicd = m_request.busicd, // 交易类型
            response = json.response;

        if (busicd === 'PURC' || busicd === 'PAUT') {
            if (response === '00') {
                order.refd = true;
            }
        }

        var chcd = m_request.chcd; // 渠道
        if (chcd === 'WXP') {
            order.img = 'img/wechat@2x.png';
            order.chcd = '微信扫码支付';
        } else {
            order.img = 'img/alipay@2x.png';
            order.chcd = '支付宝扫码支付';
        }

        var transtime = json.system_date;
        order.transtime =  transtime.substring(0, 4) + '-' + transtime.substr(4, 2) + '-' + transtime.substr(6, 2)+ ' ' + transtime.substr(8, 2)	+ ':' + transtime.substr(10, 2)	+ ':' + transtime.substr(12, 2);

        if (response === '00') {
            order.status = '交易成功';
            order.success = true;
        } else if (response === '09') {
            order.status = '未支付';
        } else {
            order.status = '交易失败';
        }

        var money = m_request.txamt;
        if (money && money.length > 0) {
            order.money = Util.getNormalTxamt(money);
        }
        order.orderNum = m_request.orderNum;
        order.username = username;
        order.goodsInfo = m_request.goodsInfo;

        return order;
    };
    var _fillOrderDetailInfoToTemplate = function(order) {
        var source = $('#orderDetailInfoTemplate').html();
        var template = Handlebars.compile(source);
        return template(order);
    };
    return {
        init: init,
        loadedData: loadedData
    };
}(window);
