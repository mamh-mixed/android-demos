// 密码修改，账户修改和限额提升的事件和业务逻辑
var PAQ = function() {
    var init = function() {

        events();
    };
    var events = function() {
        // 密码修改按钮事件
        $('#updatePwdBtn')[CLICK](function() {
            var oldpwd = document.getElementById('oldpwd').value,
                newpwd1 = document.getElementById('newpwd1').value,
                newpwd2 = document.getElementById('newpwd2').value;

            if (oldpwd.length === 0) {
                window.alert("原密码不能为空!");
                return;
            }

            if (newpwd1.length === 0) {
                window.alert("新密码不能为空!");
                return;
            }

            if (newpwd1.length < 6) {
                window.alert("密码长度不能小于六位!");
                return;
            }

            if (newpwd1 !== newpwd2) {
                window.alert("两次密码不一致!");
                return;
            }

            var json = {
                "oldpwd": oldpwd,
                "newpwd": newpwd1
            };

            var jsonStr = JSON.stringify(json);

            Util.sendRequestToDevice('updatepassword', jsonStr);
        });
        // 账户修改按钮事件
        $('#updateAccountBtn')[CLICK](function() {
            var bank_open = document.getElementById('bank_open').value,
                payee = document.getElementById('payee').value,
                payee_card = document.getElementById('payee_card').value,
                phone_num = document.getElementById('phone_num').value;

            if (bank_open.length === 0) {
                window.alert("开户行不能为空!");
                return;
            }

            if (payee.length === 0) {
                window.alert("姓名不能为空!");
                return;
            }

            if (payee_card.length === 0) {
                window.alert("银行卡号不能为空!");
                return;
            }

            if (!Util.luhmCheck(payee_card)) {
                window.alert("银行卡号验证未通过!");
                return;
            }

            if (!Util.validatemobile(phone_num)) {
                return;
            }

            var json = {
                "bank_open": bank_open,
                "payee": payee,
                "payee_card": payee_card,
                "phone_num": phone_num
            };
            var jsonStr = JSON.stringify(json);

            Util.sendRequestToDevice('updateaccount', jsonStr);
        });
        // 限额提升按钮事件
        $('#limitIncreaseBtn')[CLICK](function() {
            var limitincrease_phone_num = document.getElementById('limitincrease_phone_num').value,
                limitincrease_payee = document.getElementById('limitincrease_payee').value,
                limitincrease_email = document.getElementById('limitincrease_email').value;

            if (!Util.validatemobile(limitincrease_phone_num)) {
                return;
            }

            if (limitincrease_payee.length === 0) {
                window.alert("姓名不能为空!");
                return;
            }

            if (limitincrease_email.length === 0) {
                window.alert("邮箱不能为空!");
                return;
            }

            if (!Util.isEmail(limitincrease_email)) {
                window.alert("邮箱格式不正确!");
                return;
            }

            var json = {
                "payee": limitincrease_payee,
                "phone_num": limitincrease_phone_num,
                "email": limitincrease_email
            };
            var jsonStr = JSON.stringify(json);

            Util.sendRequestToDevice('limitincrease', jsonStr);
        });
    };
    var _fillUserInfo = function(info) {
        document.getElementById('bank_open').value = info.bank_open;
        document.getElementById('payee').value = info.payee;
        document.getElementById('payee_card').value = info.payee_card;
        document.getElementById('phone_num').value = info.phone_num;
    };
    var getCurrentAccount = function() {
        var url = UrlMap.getInfo,
            password = Util.getPassword(),
            username = Util.getUsername(),
            key = Util.getKey(),
            clientid = Util.getClientId(),
            now = new Date(),
            transtime = now.Format("YYYYMMddHHmmss"),
            dataArray = []; // 用来存储发送到服务器上的数据的key=value 的数组

        if (password) {
            password = hex_md5(password);
        }

        dataArray.push('clientid=' + clientid);
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
                var info = data.info;
                if (!info) {
                    window.alert('获取用户信息失败');
                    return;
                }
                _fillUserInfo(info);
               }
        });
    };
    return {
        init: init,
        getCurrentAccount: getCurrentAccount
    };
}(window);

var CloudCashierBridge = function() {
    var init = function(data) {

    };
    var showSection = function(target) {
        $('#'+target).show()
        .siblings('section').hide();
        $('#menu').find('a[href="#/'+target+'"]').parent().addClass('active')
        .siblings().removeClass('active');
    };
    var saveUserData = function(data) {
        Util.setUsername(data.username);
        Util.setPassword(data.password);
        Util.setDevice(data.device);
        Util.setKey(data.key);
        Util.setClientId(data.clientid);


        if (data.target !== ''){
            showSection(data.target);
        }
        if(data.target === 'transManage'){
            TransManage.freshData();
        }
    };
    return {
        init: init,
        saveUserData: saveUserData
    };
}(window);

$(function() {
    var scanPage = function() {

    };
    var transManage = function() {
        // TransManage.loaded();
        // TransManage.initScroll();
        TransManage.freshData();
    };
    var tranDetail = function(orderId) {
        console.log(orderId);
        TransDetail.loadedData(orderId);
    };
    var pwdChange = function() {
        // console.log("pwdChange");
    };
    var accountChange = function() {
        PAQ.getCurrentAccount();
    };
    var quotaExtend = function() {
        // console.log("quotaExtend");
    };

    var allroutes = function() {
        var hash = window.location.hash,
            idx = hash.lastIndexOf('/'),
            route;
        if (idx > 2) {
            route = hash.slice(2, idx);
            $('#toggleMenu').hide();
        } else {
            route = hash.slice(2);
            $('#toggleMenu').show();
            $('#menu').find('a[href$="' + route + '"]').parent().addClass('active')
                .siblings().removeClass('active');
        }
        var sections = $('section');
        var section;

        section = sections.filter('[data-route=' + route + ']');

        if (section.length) {
            sections.hide(250);
            section.show(250);
        }
    };

    //
    // define the routing table.
    //
    var routes = {
        '/': scanPage,
        '/scanPage': scanPage,
        '/transManage': transManage,
        '/tranDetail/:orderId': tranDetail,
        '/pwdChange': pwdChange,
        '/accountChange': accountChange,
        '/quotaExtend': quotaExtend
    };

    // instantiate the router.
    var router = Router(routes);
    // a global configuration setting.
    router.configure({
        on: allroutes
    });

    router.init();

    // 菜单切换事件
    $('#toggleMenu').on(CLICK, function(e) {
        e.stopPropagation();
        $('#sidebar').toggleClass('menu-show');
    });
    // 菜单切换事件
    $('#menu > li').on(CLICK, function() {
        $(this).addClass('active')
            .siblings().removeClass('active');
        // $(this).find('a').triggerHandler('tap');
        // $(this).find('a').trigger('click');
        $('#sidebar').removeClass('menu-show');
    });
    // 菜单显示后，点击主页任何地方可以关闭菜单
    $('#mainContent')[CLICK](function(e) {
        e.stopPropagation();
        if (!$('#sidebar').hasClass('menu-show')) {
            return;
        }
        $('#sidebar').removeClass('menu-show');
    });

    // 安全退出事件
    $('#safeExitBtn')[CLICK](function() {
        Util.sendRequestToDevice('safeexit', '');
    });

    if ((/android/gi).test(navigator.appVersion)){
        Util.setUsername(Util.getUrlParam("username"));
        Util.setPassword(Util.getUrlParam("password"));
        Util.setDevice(Util.getUrlParam("device"));
        Util.setKey(Util.getUrlParam("key"));
        Util.setClientId(Util.getUrlParam("clientid"));
    }

    // 计算器页面事件初始化
    Calculator.init();
    ScanCode.init();
    // 订单管理页面事件初始化
    TransManage.init();
    // 订单详情页面事件初始化
    TransDetail.init();

    PAQ.init();
});



var setAlert=function(msg){
    alert(msg);
}
