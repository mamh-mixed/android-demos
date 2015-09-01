// 这里面是扫码支付的按钮事件和逻辑处理
var ScanCode = function() {
    var init = function() {
        events();
    };
    var events = function() {
        $('#switchBtn').find('.button-switch').on('tap', function(e) {
            return;
        });
        // 商户扫码和用户扫码切换事件
        $('#switchBtn').on('tap', function(e) {
            e.preventDefault();
            e.stopPropagation();
            var $this = $(this);
            if ($this.hasClass('button-active')) {
                $this.find('.button-switch').text('商户扫码');
                $this.removeClass('button-active');
                merchantScanCode();
            } else {
                $this.find('.button-switch').text('用户扫码');
                $this.addClass('button-active');
                userScanCode();
            }
        });
        // 开始扫码或者生成二维码的按钮事件
        $('#btnScan')[CLICK](function(e) {
            scan();
        });
        $('#btnAlipay')[CLICK](function(e) {
            apay();
        });
        $('#btnTenpay')[CLICK](function(e) {
            wpay();
        });
    };
    // 用户扫码
    var userScanCode = function() {
        document.getElementById('btnScan').innerHTML = '生成二维码';
    };
    // 商户扫码
    var merchantScanCode = function() {
        document.getElementById('btnScan').innerHTML = '开始扫描';
    };
    // 阿里支付
    var apay = function() {
        document.getElementById('btnAlipay').className = 'active';
        document.getElementById('btnTenpay').className = '';
    };
    // 微信支付
    var wpay = function() {
        document.getElementById('btnAlipay').className = '';
        document.getElementById('btnTenpay').className = 'active';
    };
    // 开始扫码或者生成二维码的按钮事件
    var scan = function() {
        var sum = Calculator.getSum(),
            busicd = document.getElementById('btnScan').innerHTML,
            chcd = 'ALP';

        if (document.getElementById('btnTenpay').className === 'active') {
            chcd = 'WXP';
        }
        if (busicd === '开始扫描') {
            busicd = 'PURC';
        } else {
            busicd = 'PAUT';
        }
        if (sum === 0) {
            window.alert('金额不能为0');
            return;
        }
        var json = {
                'sum': sum,
                'busicd': busicd,
                'chcd': chcd
            },
            jsonStr = JSON.stringify(json);

        Util.sendRequestToDevice('scancode', jsonStr);
    };
    var setKeyAndClientId = function(data) {
        Util.setKey(data.key);
        Util.setClientId(data.cilentId);
    };
    return {
        init: init,
        setKeyAndClientId: setKeyAndClientId
    };
}(window);
